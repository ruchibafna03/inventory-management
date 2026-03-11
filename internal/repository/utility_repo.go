package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UtilityRepo struct {
	db *sqlx.DB
}

func NewUtilityRepo(db *sqlx.DB) *UtilityRepo {
	return &UtilityRepo{db: db}
}

// ─── Data Summary ────────────────────────────────────────────────────────────

type TableCount struct {
	Table string `db:"tbl"   json:"table"`
	Count int    `db:"count" json:"count"`
}

func (r *UtilityRepo) DataSummary() ([]TableCount, error) {
	var rows []TableCount
	err := r.db.Select(&rows, `
		SELECT tbl, cnt AS count FROM (VALUES
			('Items',        (SELECT COUNT(*) FROM item)),
			('Accounts',     (SELECT COUNT(*) FROM famst)),
			('Issues',       (SELECT COUNT(*) FROM issue)),
			('Issue Lines',  (SELECT COUNT(*) FROM isdetl)),
			('Receipts',     (SELECT COUNT(*) FROM recpt)),
			('Receipt Lines',(SELECT COUNT(*) FROM rcdetl)),
			('Sales',        (SELECT COUNT(*) FROM sale)),
			('Purchases',    (SELECT COUNT(*) FROM purch)),
			('Lots',         (SELECT COUNT(*) FROM lot)),
			('Transactions', (SELECT COUNT(*) FROM tran))
		) AS t(tbl, cnt)
	`)
	return rows, err
}

// ─── Orphan Check ────────────────────────────────────────────────────────────

type OrphanDetail struct {
	ID     int    `db:"id"     json:"id"`
	TNo    string `db:"t_no"   json:"t_no"`
	ItCd   string `db:"itcd"   json:"itcd"`
	ItDesc string `db:"itdesc" json:"itdesc"`
}

func (r *UtilityRepo) OrphanIssueDetails() ([]OrphanDetail, error) {
	var rows []OrphanDetail
	err := r.db.Select(&rows, `
		SELECT d.id, d.t_no, d.itcd, COALESCE(d.itdesc, '') AS itdesc
		FROM isdetl d
		WHERE d.itcd IS NOT NULL
		  AND NOT EXISTS (SELECT 1 FROM item i WHERE i.itcd = d.itcd)
		ORDER BY d.t_no, d.itcd
	`)
	return rows, err
}

func (r *UtilityRepo) OrphanReceiptDetails() ([]OrphanDetail, error) {
	var rows []OrphanDetail
	err := r.db.Select(&rows, `
		SELECT d.id, d.t_no, d.itcd, COALESCE(d.itdesc, '') AS itdesc
		FROM rcdetl d
		WHERE d.itcd IS NOT NULL
		  AND NOT EXISTS (SELECT 1 FROM item i WHERE i.itcd = d.itcd)
		ORDER BY d.t_no, d.itcd
	`)
	return rows, err
}

// ─── Stock Position ──────────────────────────────────────────────────────────

type StockItem struct {
	ItCd        string   `db:"itcd"         json:"itcd"`
	Desc        string   `db:"desc"         json:"desc"`
	Purity      *string  `db:"purity"       json:"purity"`
	GrossWt     float64  `db:"gross_wt"     json:"gross_wt"`
	NetWt       float64  `db:"net_wt"       json:"net_wt"`
	GoldWt      float64  `db:"gold_wt"      json:"gold_wt"`
	LotNo       *string  `db:"lot_no"       json:"lot_no"`
	Cat         *string  `db:"cat"          json:"cat"`
	Tag         *string  `db:"tag"          json:"tag"`
	IssueNo     *string  `db:"issue_no"     json:"issue_no"`
	AcCode      *string  `db:"ac_code"      json:"ac_code"`
	KarigarName string   `db:"karigar_name" json:"karigar_name"`
}

func (r *UtilityRepo) StockPosition() ([]StockItem, error) {
	var rows []StockItem
	err := r.db.Select(&rows, `
		SELECT
			i.itcd, i.desc, i.purity, i.gross_wt, i.net_wt, i.gold_wt,
			i.lot_no, i.cat, i.tag, i.issue_no, i.ac_code,
			COALESCE(f.desc, '') AS karigar_name
		FROM item i
		LEFT JOIN famst f ON f.ac_code = i.ac_code
		ORDER BY i.itcd
	`)
	return rows, err
}

// ─── Item History ────────────────────────────────────────────────────────────

type HistoryEvent struct {
	EventType string    `db:"event_type" json:"event_type"`
	EventDate time.Time `db:"event_date" json:"event_date"`
	VouchNo   string    `db:"vouch_no"   json:"vouch_no"`
	AcCode    string    `db:"ac_code"    json:"ac_code"`
	Party     string    `db:"party"      json:"party"`
	GrossWt   float64   `db:"gross_wt"   json:"gross_wt"`
	NetWt     float64   `db:"net_wt"     json:"net_wt"`
	Narr      string    `db:"narr"       json:"narr"`
}

func (r *UtilityRepo) ItemHistory(itcd string) ([]HistoryEvent, error) {
	var rows []HistoryEvent
	err := r.db.Select(&rows, `
		SELECT event_type, event_date, vouch_no, ac_code, party, gross_wt, net_wt, narr
		FROM (
			-- Lot (origin of item)
			SELECT
				'Lot Received' AS event_type,
				l.t_date       AS event_date,
				l.t_no         AS vouch_no,
				COALESCE(l.ac_code, '')  AS ac_code,
				COALESCE(l.party, '')    AS party,
				i.gross_wt,
				i.net_wt,
				COALESCE(l.desc, '')     AS narr
			FROM item i
			JOIN lot l ON l.lot_no = i.lot_no
			WHERE i.itcd = $1

			UNION ALL

			-- Issue lines
			SELECT
				'Issued'              AS event_type,
				d.t_date              AS event_date,
				d.t_no                AS vouch_no,
				COALESCE(iss.ac_code, '') AS ac_code,
				COALESCE(iss.ac_code1,'') AS party,
				d.gross_wt,
				d.net_wt,
				COALESCE(d.narr, '')  AS narr
			FROM isdetl d
			JOIN issue iss ON iss.t_no = d.t_no
			WHERE d.itcd = $1

			UNION ALL

			-- Receipt lines
			SELECT
				'Received Back'       AS event_type,
				d.t_date              AS event_date,
				d.t_no                AS vouch_no,
				COALESCE(rc.ac_code, '') AS ac_code,
				''                    AS party,
				d.gross_wt,
				d.net_wt,
				COALESCE(d.narr1, '') AS narr
			FROM rcdetl d
			JOIN recpt rc ON rc.t_no = d.t_no
			WHERE d.itcd = $1
		) h
		ORDER BY event_date, vouch_no
	`, itcd)
	return rows, err
}

// ─── Code Change ─────────────────────────────────────────────────────────────

type CodeChange struct {
	ID    int       `db:"id"    json:"id"`
	ItCdF string    `db:"itcdf" json:"itcdf"`
	ItCdT string    `db:"itcdt" json:"itcdt"`
	Date  time.Time `db:"date"  json:"date"`
}

func (r *UtilityRepo) ChangeItemCode(from, to string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`UPDATE item SET itcd = $1 WHERE itcd = $2`, to, from); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE isdetl SET itcd = $1 WHERE itcd = $2`, to, from); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE rcdetl SET itcd = $1 WHERE itcd = $2`, to, from); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE itblk SET itcd = $1 WHERE itcd = $2`, to, from); err != nil {
		return err
	}
	if _, err = tx.Exec(`INSERT INTO cdchange (itcdf, itcdt, "date") VALUES ($1, $2, CURRENT_DATE)`, from, to); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *UtilityRepo) CodeChangeHistory() ([]CodeChange, error) {
	var rows []CodeChange
	err := r.db.Select(&rows, `
		SELECT id, itcdf, itcdt, "date" FROM cdchange ORDER BY id DESC LIMIT 200
	`)
	return rows, err
}

// ─── Blocked Items ───────────────────────────────────────────────────────────

type BlockedItem struct {
	ItCd    string  `db:"itcd"     json:"itcd"`
	Desc    string  `db:"desc"     json:"desc"`
	GrossWt float64 `db:"gross_wt" json:"gross_wt"`
}

func (r *UtilityRepo) BlockedItems() ([]BlockedItem, error) {
	var rows []BlockedItem
	err := r.db.Select(&rows, `
		SELECT i.itcd, i.desc, i.gross_wt
		FROM itblk b
		JOIN item i ON i.itcd = b.itcd
		ORDER BY i.itcd
	`)
	return rows, err
}

func (r *UtilityRepo) BlockItem(itcd string) error {
	_, err := r.db.Exec(`INSERT INTO itblk(itcd) VALUES($1) ON CONFLICT DO NOTHING`, itcd)
	return err
}

func (r *UtilityRepo) UnblockItem(itcd string) error {
	_, err := r.db.Exec(`DELETE FROM itblk WHERE itcd = $1`, itcd)
	return err
}

// ─── Password Change ─────────────────────────────────────────────────────────

type UserRecord struct {
	ID       int    `db:"id"        json:"id"`
	Username string `db:"username"  json:"username"`
	FullName string `db:"full_name" json:"full_name"`
	Role     string `db:"role"      json:"role"`
	Active   bool   `db:"active"    json:"active"`
}

func (r *UtilityRepo) ListUsers() ([]UserRecord, error) {
	var rows []UserRecord
	err := r.db.Select(&rows, `
		SELECT id, username, COALESCE(full_name,'') AS full_name, role, active
		FROM users ORDER BY username
	`)
	return rows, err
}

func (r *UtilityRepo) ChangePassword(username, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	res, err := r.db.Exec(
		`UPDATE users SET password_hash = $1 WHERE username = $2`,
		string(hash), username,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("user %q not found", username)
	}
	return nil
}
