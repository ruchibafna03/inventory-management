package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/val/inventory/internal/models"
)

type RateRepo struct {
	db *sqlx.DB
}

func NewRateRepo(db *sqlx.DB) *RateRepo {
	return &RateRepo{db: db}
}

func (r *RateRepo) List() ([]models.Rate, error) {
	var rates []models.Rate
	err := r.db.Select(&rates, `SELECT * FROM rate ORDER BY date DESC`)
	return rates, err
}

func (r *RateRepo) Latest() (*models.Rate, error) {
	var rate models.Rate
	err := r.db.Get(&rate, `SELECT * FROM rate ORDER BY date DESC LIMIT 1`)
	return &rate, err
}

func (r *RateRepo) Create(rate *models.Rate) error {
	_, err := r.db.NamedExec(`
		INSERT INTO rate ("date", rate, rate1, s_rate, gtag, qtag, xtag, sprn)
		VALUES (:date, :rate, :rate1, :s_rate, :gtag, :qtag, :xtag, :sprn)
		ON CONFLICT ("date") DO UPDATE SET
			rate=EXCLUDED.rate, rate1=EXCLUDED.rate1, s_rate=EXCLUDED.s_rate,
			gtag=EXCLUDED.gtag, qtag=EXCLUDED.qtag, xtag=EXCLUDED.xtag, sprn=EXCLUDED.sprn`, rate)
	return err
}

func (r *RateRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM rate WHERE id = $1`, id)
	return err
}

// ─── Lots ────────────────────────────────────────────────────────────────────

type LotRepo struct {
	db *sqlx.DB
}

func NewLotRepo(db *sqlx.DB) *LotRepo {
	return &LotRepo{db: db}
}

func (r *LotRepo) List(acCode string) ([]models.Lot, error) {
	where := ""
	args := []any{}
	if acCode != "" {
		where = " WHERE l.ac_code = $1"
		args = append(args, acCode)
	}
	var lots []models.Lot
	err := r.db.Select(&lots, `SELECT * FROM lot`+where+` ORDER BY t_date DESC`, args...)
	return lots, err
}

func (r *LotRepo) Get(tNo string) (*models.Lot, error) {
	var lot models.Lot
	err := r.db.Get(&lot, `SELECT * FROM lot WHERE t_no = $1`, tNo)
	return &lot, err
}

func (r *LotRepo) Create(lot *models.Lot) error {
	_, err := r.db.NamedExec(`
		INSERT INTO lot (t_no, t_date, ac_code, party, lot_no, nos, gross_wt, kundan_wt, nakash_wt, stone_wt, pearl_wt, nos1, gross_wt1, adj_wt, kundan_wt1, nakash_wt1, stone_wt1, pearl_wt1, amt, g_rt, "desc")
		VALUES (:t_no, :t_date, :ac_code, :party, :lot_no, :nos, :gross_wt, :kundan_wt, :nakash_wt, :stone_wt, :pearl_wt, :nos1, :gross_wt1, :adj_wt, :kundan_wt1, :nakash_wt1, :stone_wt1, :pearl_wt1, :amt, :g_rt, :desc)`, lot)
	return err
}

func (r *LotRepo) Update(lot *models.Lot) error {
	_, err := r.db.NamedExec(`
		UPDATE lot SET
			t_date=:t_date, ac_code=:ac_code, party=:party, lot_no=:lot_no,
			nos=:nos, gross_wt=:gross_wt, kundan_wt=:kundan_wt,
			nakash_wt=:nakash_wt, stone_wt=:stone_wt, pearl_wt=:pearl_wt,
			nos1=:nos1, gross_wt1=:gross_wt1, adj_wt=:adj_wt,
			kundan_wt1=:kundan_wt1, nakash_wt1=:nakash_wt1,
			stone_wt1=:stone_wt1, pearl_wt1=:pearl_wt1,
			amt=:amt, g_rt=:g_rt, "desc"=:desc
		WHERE t_no=:t_no`, lot)
	return err
}

func (r *LotRepo) Delete(tNo string) error {
	_, err := r.db.Exec(`DELETE FROM lot WHERE t_no = $1`, tNo)
	return err
}
