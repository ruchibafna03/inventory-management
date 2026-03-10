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

	_ "github.com/lib/pq"
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

	stmt, err := imp.db.Prepare(`
		INSERT INTO grpmst (g_code, g_desc)
		VALUES ($1, $2)
		ON CONFLICT (g_code) DO UPDATE SET g_desc=EXCLUDED.g_desc`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range records {
		if _, err := stmt.Exec(r["g_code"], r["g_desc"]); err != nil {
			log.Printf("  GRPMST skip %v: %v", r["g_code"], err)
		}
	}
	return nil
}

func (imp *Importer) importFamSt() error {
	dbf, err := imp.openDBF("FAMST")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("FAMST: %d records", len(records))

	stmt, err := imp.db.Prepare(`
		INSERT INTO famst (ac_code, sub_code, "desc", a_name, amt, rmk, rate, per, opb, "date", cat, dia)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (ac_code) DO UPDATE SET
			"desc"=EXCLUDED."desc", amt=EXCLUDED.amt, opb=EXCLUDED.opb`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, r := range records {
		if _, err := stmt.Exec(
			r["ac_code"], str(r, "sub_code"), r["desc"], str(r, "a_name"),
			num(r, "amt"), str(r, "rmk"), num(r, "rate"), num(r, "per"),
			num(r, "opb"), date(r, "date"), str(r, "cat"), str(r, "dia"),
		); err != nil {
			log.Printf("  FAMST skip %v: %v", r["ac_code"], err)
		}
	}
	return nil
}

func (imp *Importer) importItem() error {
	dbf, err := imp.openDBF("ITEM")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("ITEM: %d records", len(records))

	for _, r := range records {
		_, err := imp.db.Exec(`
			INSERT INTO item (
				itcd, set_code, "desc", fac, stamp, purity, melt,
				gross_wt, net_wt, ghat_wt, kundan_wt, extra_wt, gold_wt,
				kundn_wt, nakash_wt, stone_wt, prl_wt,
				rby_ct, rby_gm, eme_ct, eme_gm, plk_ct, plk_gm,
				stn1_name, stn1_ct, stn1_gm, stn2_name, stn2_ct, stn2_gm,
				stn3_name, stn3_ct, stn3_gm,
				prls_gm, prlb_gm, rbyb_gm, emeb_gm, emebb_gm,
				prl1_name, prl1_gm, srby_gm, seme_gm, scz_gm, snav_gm,
				prl2_name, prl2_gm, prl3_name, prl3_gm,
				narr1, narr2, mk_chrg, w_per,
				recpt_date, lot_no, issue_no, issue_date,
				ac_code, stk_code, cat, kp, tag, cat1, cna, exb
			) VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,
				$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,
				$33,$34,$35,$36,$37,$38,$39,$40,$41,$42,$43,$44,$45,$46,$47,
				$48,$49,$50,$51,$52,$53,$54,$55,$56,$57,$58,$59,$60,$61,$62
			) ON CONFLICT (itcd) DO NOTHING`,
			r["itcd"], str(r, "set"), r["desc"], num(r, "fac"),
			str(r, "stamp"), str(r, "purity"), num(r, "melt"),
			num(r, "gross_wt"), num(r, "net_wt"), num(r, "ghat_wt"),
			num(r, "kundan_wt"), num(r, "extra_wt"), num(r, "gold_wt"),
			num(r, "kundn_wt"), num(r, "nakash_wt"), num(r, "stone_wt"), num(r, "prl_wt"),
			num(r, "rby_ct"), num(r, "rby_gm"), num(r, "eme_ct"), num(r, "eme_gm"),
			num(r, "plk_ct"), num(r, "plk_gm"),
			str(r, "stn1_name"), num(r, "stn1_ct"), num(r, "stn1_gm"),
			str(r, "stn2_name"), num(r, "stn2_ct"), num(r, "stn2_gm"),
			str(r, "stn3_name"), num(r, "stn3_ct"), num(r, "stn3_gm"),
			num(r, "prls_gm"), num(r, "prlb_gm"), num(r, "rbyb_gm"),
			num(r, "emeb_gm"), num(r, "emebb_gm"),
			str(r, "prl1_name"), num(r, "prl1_gm"),
			num(r, "srby_gm"), num(r, "seme_gm"), num(r, "scz_gm"), num(r, "snav_gm"),
			str(r, "prl2_name"), num(r, "prl2_gm"),
			str(r, "prl3_name"), num(r, "prl3_gm"),
			str(r, "narr1"), str(r, "narr2"), num(r, "mk_chrg"), num(r, "w_per"),
			date(r, "recpt_date"), str(r, "lot_no"), str(r, "issue_no"), date(r, "issue_date"),
			str(r, "ac_code"), str(r, "stk_code"), str(r, "cat"),
			str(r, "kp"), str(r, "tag"), str(r, "cat1"), str(r, "cna"), str(r, "exb"),
		)
		if err != nil {
			log.Printf("  ITEM skip %v: %v", r["itcd"], err)
		}
	}
	return nil
}

func (imp *Importer) importLot() error {
	dbf, err := imp.openDBF("LOT")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("LOT: %d records", len(records))

	for _, r := range records {
		_, err := imp.db.Exec(`
			INSERT INTO lot (t_no, t_date, ac_code, party, lot_no, nos, gross_wt, kundan_wt, nakash_wt, stone_wt, pearl_wt, nos1, gross_wt1, adj_wt, kundan_wt1, nakash_wt1, stone_wt1, pearl_wt1, amt, g_rt, "desc")
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)
			ON CONFLICT (t_no) DO NOTHING`,
			r["t_no"], date(r, "t_date"), str(r, "ac_code"), str(r, "party"), r["lot_no"],
			int(num(r, "nos")), num(r, "gross_wt"), num(r, "kundan_wt"), num(r, "nakash_wt"),
			num(r, "stone_wt"), num(r, "pearl_wt"), int(num(r, "nos1")), num(r, "gross_wt1"),
			num(r, "adj_wt"), num(r, "kundan_wt1"), num(r, "nakash_wt1"),
			num(r, "stone_wt1"), num(r, "pearl_wt1"), num(r, "amt"), num(r, "g_rt"), str(r, "desc"),
		)
		if err != nil {
			log.Printf("  LOT skip %v: %v", r["t_no"], err)
		}
	}
	return nil
}

func (imp *Importer) importRate() error {
	dbf, err := imp.openDBF("RATE")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("RATE: %d records", len(records))

	for _, r := range records {
		d := date(r, "date")
		if !d.Valid {
			continue
		}
		_, err := imp.db.Exec(`
			INSERT INTO rate ("date", rate, rate1, s_rate, gtag, qtag, xtag, sprn)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
			ON CONFLICT ("date") DO UPDATE SET rate=EXCLUDED.rate, rate1=EXCLUDED.rate1, s_rate=EXCLUDED.s_rate`,
			d.Time, num(r, "rate"), num(r, "rate1"), num(r, "s_rate"),
			r["gtag"], r["qtag"], r["xtag"], r["sprn"],
		)
		if err != nil {
			log.Printf("  RATE skip: %v", err)
		}
	}
	return nil
}

func (imp *Importer) importSale() error {
	dbf, err := imp.openDBF("SALE")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("SALE: %d records", len(records))

	for _, r := range records {
		d := date(r, "vouch_date")
		if !d.Valid {
			continue
		}
		_, err := imp.db.Exec(`
			INSERT INTO sale (vouch_no, ovno, pre, vouch_no1, vouch_date, ac_code, name, add1, add2, phone_no, phone_no1, sman, nos, gross_wt, kundn_wt, nakash_wt, stone_wt, net_wt, net_per, kdn_per, nks_per, net_pure, kdn_pure, nks_pure, stn_amt, stn_rate, stn_pure, pure_wt, discount, pure_final, narr, narr1, btype)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33)
			ON CONFLICT (vouch_no) DO NOTHING`,
			r["vouch_no"], str(r, "ovno"), str(r, "pre"), str(r, "vouch_no1"),
			d.Time, str(r, "ac_code"), str(r, "name"), str(r, "add1"), str(r, "add2"),
			str(r, "phone_no"), str(r, "phone_no1"), str(r, "sman"),
			int(num(r, "nos")), num(r, "gross_wt"), num(r, "kundn_wt"),
			num(r, "nakash_wt"), num(r, "stone_wt"), num(r, "net_wt"),
			num(r, "net_per"), num(r, "kdn_per"), num(r, "nks_per"),
			num(r, "net_pure"), num(r, "kdn_pure"), num(r, "nks_pure"),
			num(r, "stn_amt"), num(r, "stn_rate"), num(r, "stn_pure"),
			num(r, "pure_wt"), num(r, "discount"), num(r, "pure_final"),
			str(r, "narr"), str(r, "narr1"), str(r, "btype"),
		)
		if err != nil {
			log.Printf("  SALE skip %v: %v", r["vouch_no"], err)
		}
	}
	return nil
}

func (imp *Importer) importPurch() error {
	dbf, err := imp.openDBF("PURCH")
	if err != nil {
		return err
	}
	defer dbf.Close()

	records := dbf.ReadAll()
	log.Printf("PURCH: %d records", len(records))

	for _, r := range records {
		d := date(r, "vouch_date")
		if !d.Valid {
			continue
		}
		_, err := imp.db.Exec(`
			INSERT INTO purch (vouch_no, pre, vouch_no1, vouch_date, ac_code, name, add1, add2, rate, gross_wt, net_wt, bill_amt, narr, rmk, tag, post)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
			ON CONFLICT (vouch_no) DO NOTHING`,
			r["vouch_no"], str(r, "pre"), str(r, "vouch_no1"), d.Time,
			str(r, "ac_code"), str(r, "name"), str(r, "add1"), str(r, "add2"),
			num(r, "rate"), num(r, "gross_wt"), num(r, "net_wt"), num(r, "bill_amt"),
			str(r, "narr"), str(r, "rmk"), str(r, "tag"), str(r, "post"),
		)
		if err != nil {
			log.Printf("  PURCH skip %v: %v", r["vouch_no"], err)
		}
	}
	return nil
}

func (imp *Importer) Run() {
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
	imp.Run()
}
