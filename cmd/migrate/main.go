// migrate imports legacy FoxPlus .DBF files into PostgreSQL.
//
// Usage:
//
//	go run ./cmd/migrate --dbf /path/to/VAL --dsn "postgres://..."
package main

import (
	"database/sql"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lib/pq"
)

// ─── DBF reader ──────────────────────────────────────────────────────────────

type DBFField struct {
	Name   string
	Type   byte // C, N, D, L
	Length int
	Dec    int
}

type DBFFile struct {
	NumRecords int
	Fields     []DBFField
	RecordSize int
	f          *os.File
	dataStart  int64
}

func OpenDBF(path string) (*DBFFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	hdr := make([]byte, 32)
	if _, err = io.ReadFull(f, hdr); err != nil {
		return nil, err
	}

	numRec := int(binary.LittleEndian.Uint32(hdr[4:8]))
	hdrSize := int(binary.LittleEndian.Uint16(hdr[8:10]))
	recSize := int(binary.LittleEndian.Uint16(hdr[10:12]))

	numFields := (hdrSize - 32 - 1) / 32
	fields := make([]DBFField, 0, numFields)

	for i := 0; i < numFields; i++ {
		fd := make([]byte, 32)
		if _, err = io.ReadFull(f, fd); err != nil {
			break
		}
		if fd[0] == 0x0D {
			break
		}
		name := strings.TrimRight(string(fd[0:11]), "\x00")
		fields = append(fields, DBFField{
			Name:   strings.ToLower(name),
			Type:   fd[11],
			Length: int(fd[16]),
			Dec:    int(fd[17]),
		})
	}

	// Seek to data start
	if _, err = f.Seek(int64(hdrSize), io.SeekStart); err != nil {
		return nil, err
	}

	return &DBFFile{
		NumRecords: numRec,
		Fields:     fields,
		RecordSize: recSize,
		f:          f,
		dataStart:  int64(hdrSize),
	}, nil
}

func (d *DBFFile) Close() { d.f.Close() }

// ReadRecord returns a map of fieldName->string for each record.
// Returns nil when EOF.
func (d *DBFFile) ReadRecord() map[string]string {
	buf := make([]byte, d.RecordSize)
	if _, err := io.ReadFull(d.f, buf); err != nil {
		return nil
	}
	if buf[0] == 0x1A { // EOF marker
		return nil
	}
	if buf[0] == '*' { // deleted record
		return map[string]string{"__deleted": "1"}
	}

	rec := make(map[string]string, len(d.Fields))
	pos := 1
	for _, f := range d.Fields {
		raw := string(buf[pos : pos+f.Length])
		pos += f.Length
		switch f.Type {
		case 'C':
			rec[f.Name] = strings.TrimRight(raw, " ")
		case 'N':
			rec[f.Name] = strings.TrimSpace(raw)
		case 'D':
			raw = strings.TrimSpace(raw)
			if len(raw) == 8 {
				// YYYYMMDD
				rec[f.Name] = raw[:4] + "-" + raw[4:6] + "-" + raw[6:8]
			} else {
				rec[f.Name] = ""
			}
		case 'L':
			rec[f.Name] = strings.TrimSpace(raw)
		default:
			rec[f.Name] = strings.TrimSpace(raw)
		}
	}
	return rec
}

// ReadAll reads all non-deleted records.
func (d *DBFFile) ReadAll() []map[string]string {
	var records []map[string]string
	for {
		rec := d.ReadRecord()
		if rec == nil {
			break
		}
		if rec["__deleted"] == "1" {
			continue
		}
		records = append(records, rec)
	}
	return records
}

// ─── Importer ────────────────────────────────────────────────────────────────

type Importer struct {
	db     *sql.DB
	dbfDir string
}

func (imp *Importer) openDBF(name string) (*DBFFile, error) {
	// Try exact case, then upper
	paths := []string{
		filepath.Join(imp.dbfDir, name+".DBF"),
		filepath.Join(imp.dbfDir, strings.ToUpper(name)+".DBF"),
		filepath.Join(imp.dbfDir, strings.ToLower(name)+".dbf"),
	}
	for _, p := range paths {
		if dbf, err := OpenDBF(p); err == nil {
			return dbf, nil
		}
	}
	return nil, fmt.Errorf("DBF not found: %s", name)
}

func str(rec map[string]string, key string) sql.NullString {
	v := rec[key]
	if v == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: v, Valid: true}
}

func date(rec map[string]string, key string) sql.NullTime {
	v := rec[key]
	if v == "" || v == "0000-00-00" {
		return sql.NullTime{}
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// strVal returns the string value or nil for use with pq.CopyIn.
func strVal(rec map[string]string, key string) interface{} {
	v := rec[key]
	if v == "" {
		return nil
	}
	return v
}

// dateVal returns a *time.Time or nil for use with pq.CopyIn.
func dateVal(rec map[string]string, key string) interface{} {
	v := rec[key]
	if v == "" || v == "0000-00-00" {
		return nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil
	}
	return t
}

func num(rec map[string]string, key string) float64 {
	v := strings.TrimSpace(rec[key])
	if v == "" {
		return 0
	}
	var f float64
	fmt.Sscanf(v, "%f", &f)
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}
	return f
}

func (imp *Importer) importGrpMst() error {
	dbf, err := imp.openDBF("GRPMST")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("GRPMST: %d records", len(records))

	tx, err := imp.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("grpmst", "g_code", "g_desc"))
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, r := range records {
		stmt.Exec(r["g_code"], r["g_desc"])
	}
	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return err
	}
	stmt.Close()
	return tx.Commit()
}

func (imp *Importer) importFamSt() error {
	dbf, err := imp.openDBF("FAMST")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("FAMST: %d records", len(records))

	tx, err := imp.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("famst", "ac_code", "sub_code", "desc", "a_name", "amt", "rmk", "rate", "per", "opb", "date", "cat", "dia"))
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, r := range records {
		stmt.Exec(
			r["ac_code"], strVal(r, "sub_code"), r["desc"], strVal(r, "a_name"),
			num(r, "amt"), strVal(r, "rmk"), num(r, "rate"), num(r, "per"),
			num(r, "opb"), dateVal(r, "date"), strVal(r, "cat"), strVal(r, "dia"),
		)
	}
	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return err
	}
	stmt.Close()
	return tx.Commit()
}

func (imp *Importer) importItem() error {
	dbf, err := imp.openDBF("ITEM")
	if err != nil {
		return err
	}
	defer dbf.Close()

	allRecords := dbf.ReadAll()
	seen := map[string]bool{}
	records := allRecords[:0]
	for _, r := range allRecords {
		if !seen[r["itcd"]] {
			seen[r["itcd"]] = true
			records = append(records, r)
		}
	}
	log.Printf("ITEM: %d records (%d dupes dropped)", len(records), len(allRecords)-len(records))

	// Use a temp table (no constraints) for COPY, then INSERT ... ON CONFLICT.
	imp.db.Exec(`DROP TABLE IF EXISTS item_tmp`)
	if _, err := imp.db.Exec(`CREATE TEMP TABLE item_tmp (LIKE item)`); err != nil {
		return fmt.Errorf("create item_tmp: %w", err)
	}

	tx, err := imp.db.Begin()
	if err != nil {
		return err
	}
	cols := []string{
		"itcd", "set_code", "desc", "fac", "stamp", "purity", "melt",
		"gross_wt", "net_wt", "ghat_wt", "kundan_wt", "extra_wt", "gold_wt",
		"kundn_wt", "nakash_wt", "stone_wt", "prl_wt",
		"rby_ct", "rby_gm", "eme_ct", "eme_gm", "plk_ct", "plk_gm",
		"stn1_name", "stn1_ct", "stn1_gm", "stn2_name", "stn2_ct", "stn2_gm",
		"stn3_name", "stn3_ct", "stn3_gm",
		"prls_gm", "prlb_gm", "rbyb_gm", "emeb_gm", "emebb_gm",
		"prl1_name", "prl1_gm", "srby_gm", "seme_gm", "scz_gm", "snav_gm",
		"prl2_name", "prl2_gm", "prl3_name", "prl3_gm",
		"narr1", "narr2", "mk_chrg", "w_per",
		"recpt_date", "lot_no", "issue_no", "issue_date",
		"ac_code", "stk_code", "cat", "kp", "tag", "cat1", "cna", "exb",
	}
	stmt, err := tx.Prepare(pq.CopyIn("item_tmp", cols...))
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, r := range records {
		stmt.Exec(
			r["itcd"], strVal(r, "set"), r["desc"], num(r, "fac"),
			strVal(r, "stamp"), strVal(r, "purity"), num(r, "melt"),
			num(r, "gross_wt"), num(r, "net_wt"), num(r, "ghat_wt"),
			num(r, "kundan_wt"), num(r, "extra_wt"), num(r, "gold_wt"),
			num(r, "kundn_wt"), num(r, "nakash_wt"), num(r, "stone_wt"), num(r, "prl_wt"),
			num(r, "rby_ct"), num(r, "rby_gm"), num(r, "eme_ct"), num(r, "eme_gm"),
			num(r, "plk_ct"), num(r, "plk_gm"),
			strVal(r, "stn1_name"), num(r, "stn1_ct"), num(r, "stn1_gm"),
			strVal(r, "stn2_name"), num(r, "stn2_ct"), num(r, "stn2_gm"),
			strVal(r, "stn3_name"), num(r, "stn3_ct"), num(r, "stn3_gm"),
			num(r, "prls_gm"), num(r, "prlb_gm"), num(r, "rbyb_gm"),
			num(r, "emeb_gm"), num(r, "emebb_gm"),
			strVal(r, "prl1_name"), num(r, "prl1_gm"),
			num(r, "srby_gm"), num(r, "seme_gm"), num(r, "scz_gm"), num(r, "snav_gm"),
			strVal(r, "prl2_name"), num(r, "prl2_gm"),
			strVal(r, "prl3_name"), num(r, "prl3_gm"),
			strVal(r, "narr1"), strVal(r, "narr2"), num(r, "mk_chrg"), num(r, "w_per"),
			dateVal(r, "recpt_date"), strVal(r, "lot_no"), strVal(r, "issue_no"), dateVal(r, "issue_date"),
			strVal(r, "ac_code"), strVal(r, "stk_code"), strVal(r, "cat"),
			strVal(r, "kp"), strVal(r, "tag"), strVal(r, "cat1"), strVal(r, "cna"), strVal(r, "exb"),
		)
	}
	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return err
	}
	stmt.Close()
	// Move from temp to real table; set timestamps since COPY doesn't supply them.
	if _, err := tx.Exec(`INSERT INTO item SELECT *, NOW(), NOW() FROM item_tmp ON CONFLICT (itcd) DO NOTHING`); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (imp *Importer) importLot() error {
	dbf, err := imp.openDBF("LOT")
	if err != nil {
		return err
	}
	defer dbf.Close()

	allRecords := dbf.ReadAll()
	seen := map[string]bool{}
	records := allRecords[:0]
	for _, r := range allRecords {
		if !seen[r["t_no"]] {
			seen[r["t_no"]] = true
			records = append(records, r)
		}
	}
	log.Printf("LOT: %d records (%d dupes dropped)", len(records), len(allRecords)-len(records))

	tx, err := imp.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("lot",
		"t_no", "t_date", "ac_code", "party", "lot_no",
		"nos", "gross_wt", "kundan_wt", "nakash_wt", "stone_wt", "pearl_wt",
		"nos1", "gross_wt1", "adj_wt", "kundan_wt1", "nakash_wt1",
		"stone_wt1", "pearl_wt1", "amt", "g_rt", "desc",
	))
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, r := range records {
		stmt.Exec(
			r["t_no"], dateVal(r, "t_date"), strVal(r, "ac_code"), strVal(r, "party"), r["lot_no"],
			int(num(r, "nos")), num(r, "gross_wt"), num(r, "kundan_wt"), num(r, "nakash_wt"),
			num(r, "stone_wt"), num(r, "pearl_wt"), int(num(r, "nos1")), num(r, "gross_wt1"),
			num(r, "adj_wt"), num(r, "kundan_wt1"), num(r, "nakash_wt1"),
			num(r, "stone_wt1"), num(r, "pearl_wt1"), num(r, "amt"), num(r, "g_rt"), strVal(r, "desc"),
		)
	}
	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return err
	}
	stmt.Close()
	return tx.Commit()
}

func (imp *Importer) importRate() error {
	dbf, err := imp.openDBF("RATE")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("RATE: %d records", len(records))

	tx, err := imp.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("rate", "date", "rate", "rate1", "s_rate", "gtag", "qtag", "xtag", "sprn"))
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, r := range records {
		d := dateVal(r, "date")
		if d == nil {
			continue
		}
		stmt.Exec(d, num(r, "rate"), num(r, "rate1"), num(r, "s_rate"),
			r["gtag"], r["qtag"], r["xtag"], r["sprn"])
	}
	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return err
	}
	stmt.Close()
	return tx.Commit()
}

func (imp *Importer) importSale() error {
	dbf, err := imp.openDBF("SALE")
	if err != nil {
		return err
	}
	defer dbf.Close()

	allRecords := dbf.ReadAll()
	seen := map[string]bool{}
	records := allRecords[:0]
	for _, r := range allRecords {
		if !seen[r["vouch_no"]] {
			seen[r["vouch_no"]] = true
			records = append(records, r)
		}
	}
	log.Printf("SALE: %d records (%d dupes dropped)", len(records), len(allRecords)-len(records))

	tx, err := imp.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("sale",
		"vouch_no", "ovno", "pre", "vouch_no1", "vouch_date",
		"ac_code", "name", "add1", "add2", "phone_no", "phone_no1", "sman",
		"nos", "gross_wt", "kundn_wt", "nakash_wt", "stone_wt", "net_wt",
		"net_per", "kdn_per", "nks_per", "net_pure", "kdn_pure", "nks_pure",
		"stn_amt", "stn_rate", "stn_pure", "pure_wt", "discount", "pure_final",
		"narr", "narr1", "btype",
	))
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, r := range records {
		d := dateVal(r, "vouch_date")
		if d == nil {
			continue
		}
		stmt.Exec(
			r["vouch_no"], strVal(r, "ovno"), strVal(r, "pre"), strVal(r, "vouch_no1"),
			d, strVal(r, "ac_code"), strVal(r, "name"), strVal(r, "add1"), strVal(r, "add2"),
			strVal(r, "phone_no"), strVal(r, "phone_no1"), strVal(r, "sman"),
			int(num(r, "nos")), num(r, "gross_wt"), num(r, "kundn_wt"),
			num(r, "nakash_wt"), num(r, "stone_wt"), num(r, "net_wt"),
			num(r, "net_per"), num(r, "kdn_per"), num(r, "nks_per"),
			num(r, "net_pure"), num(r, "kdn_pure"), num(r, "nks_pure"),
			num(r, "stn_amt"), num(r, "stn_rate"), num(r, "stn_pure"),
			num(r, "pure_wt"), num(r, "discount"), num(r, "pure_final"),
			strVal(r, "narr"), strVal(r, "narr1"), strVal(r, "btype"),
		)
	}
	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return err
	}
	stmt.Close()
	return tx.Commit()
}

func (imp *Importer) importPurch() error {
	dbf, err := imp.openDBF("PURCH")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("PURCH: %d records", len(records))

	tx, err := imp.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(pq.CopyIn("purch",
		"vouch_no", "pre", "vouch_no1", "vouch_date",
		"ac_code", "name", "add1", "add2",
		"rate", "gross_wt", "net_wt", "bill_amt",
		"narr", "rmk", "tag", "post",
	))
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, r := range records {
		d := dateVal(r, "vouch_date")
		if d == nil {
			continue
		}
		stmt.Exec(
			r["vouch_no"], strVal(r, "pre"), strVal(r, "vouch_no1"), d,
			strVal(r, "ac_code"), strVal(r, "name"), strVal(r, "add1"), strVal(r, "add2"),
			num(r, "rate"), num(r, "gross_wt"), num(r, "net_wt"), num(r, "bill_amt"),
			strVal(r, "narr"), strVal(r, "rmk"), strVal(r, "tag"), strVal(r, "post"),
		)
	}
	if _, err := stmt.Exec(); err != nil {
		tx.Rollback()
		return err
	}
	stmt.Close()
	return tx.Commit()
}

func (imp *Importer) truncateInTx(tx *sql.Tx, tables ...string) error {
	for _, t := range tables {
		if _, err := tx.Exec("TRUNCATE TABLE " + t + " CASCADE"); err != nil {
			return fmt.Errorf("truncate %s: %w", t, err)
		}
		log.Printf("  truncated %s", t)
	}
	return nil
}

func (imp *Importer) Run(doTruncate bool) {
	if doTruncate {
		log.Println("==> Truncating all tables...")
		tx, err := imp.db.Begin()
		if err != nil {
			log.Fatalf("truncate: begin: %v", err)
		}
		if err := imp.truncateInTx(tx, "purch", "sale", "rate", "lot", "item", "famst", "grpmst"); err != nil {
			tx.Rollback()
			log.Fatalf("truncate failed: %v", err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("truncate commit: %v", err)
		}
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Account Groups", imp.importGrpMst},
		{"Account Master", imp.importFamSt},
		{"Items", imp.importItem},
		{"Lots", imp.importLot},
		{"Gold Rates", imp.importRate},
		{"Sales", imp.importSale},
		{"Purchases", imp.importPurch},
	}

	for _, step := range steps {
		log.Printf("==> Importing %s...", step.name)
		if err := step.fn(); err != nil {
			log.Printf("    WARNING: %v", err)
		} else {
			log.Printf("    OK")
		}
	}
	log.Println("Migration complete.")
}

func main() {
	dbfDir := flag.String("dbf", ".", "directory containing .DBF files")
	dsn := flag.String("dsn", "", "PostgreSQL DSN (or set DATABASE_URL env)")
	truncate := flag.Bool("truncate", false, "truncate all tables before importing")
	flag.Parse()

	if *dsn == "" {
		*dsn = os.Getenv("DATABASE_URL")
	}
	if *dsn == "" {
		*dsn = "postgres://postgres:postgres@localhost:5432/val_inventory?sslmode=disable"
	}

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	imp := &Importer{db: db, dbfDir: *dbfDir}
	imp.Run(*truncate)
}
