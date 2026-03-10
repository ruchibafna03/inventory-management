package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/val/inventory/internal/models"
)

type SaleRepo struct {
	db *sqlx.DB
}

func NewSaleRepo(db *sqlx.DB) *SaleRepo {
	return &SaleRepo{db: db}
}

func (r *SaleRepo) List(page, perPage int, acCode string) ([]models.Sale, int, error) {
	offset := (page - 1) * perPage
	where := "WHERE 1=1"
	args := []any{}
	i := 1

	if acCode != "" {
		where += fmt.Sprintf(" AND ac_code = $%d", i)
		args = append(args, acCode)
		i++
	}

	var total int
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM sale "+where, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, perPage, offset)
	query := fmt.Sprintf(`SELECT * FROM sale %s ORDER BY vouch_date DESC LIMIT $%d OFFSET $%d`, where, i, i+1)
	var sales []models.Sale
	if err := r.db.Select(&sales, query, args...); err != nil {
		return nil, 0, err
	}
	return sales, total, nil
}

func (r *SaleRepo) Get(vouchNo string) (*models.Sale, error) {
	var sale models.Sale
	err := r.db.Get(&sale, `SELECT * FROM sale WHERE vouch_no = $1`, vouchNo)
	return &sale, err
}

func (r *SaleRepo) Create(sale *models.Sale) error {
	_, err := r.db.NamedExec(`
		INSERT INTO sale (
			vouch_no, ovno, pre, vouch_no1, vouch_date, ac_code, name, add1, add2,
			phone_no, phone_no1, sman, nos, gross_wt, kundn_wt, nakash_wt, stone_wt, net_wt,
			net_per, kdn_per, nks_per, net_pure, kdn_pure, nks_pure,
			stn_amt, stn_rate, stn_pure, pure_wt, discount, pure_final,
			narr, narr1, btype
		) VALUES (
			:vouch_no, :ovno, :pre, :vouch_no1, :vouch_date, :ac_code, :name, :add1, :add2,
			:phone_no, :phone_no1, :sman, :nos, :gross_wt, :kundn_wt, :nakash_wt, :stone_wt, :net_wt,
			:net_per, :kdn_per, :nks_per, :net_pure, :kdn_pure, :nks_pure,
			:stn_amt, :stn_rate, :stn_pure, :pure_wt, :discount, :pure_final,
			:narr, :narr1, :btype
		)`, sale)
	return err
}

func (r *SaleRepo) Update(sale *models.Sale) error {
	_, err := r.db.NamedExec(`
		UPDATE sale SET
			ovno=:ovno, pre=:pre, vouch_no1=:vouch_no1, vouch_date=:vouch_date,
			ac_code=:ac_code, name=:name, add1=:add1, add2=:add2,
			phone_no=:phone_no, phone_no1=:phone_no1, sman=:sman,
			nos=:nos, gross_wt=:gross_wt, kundn_wt=:kundn_wt,
			nakash_wt=:nakash_wt, stone_wt=:stone_wt, net_wt=:net_wt,
			net_per=:net_per, kdn_per=:kdn_per, nks_per=:nks_per,
			net_pure=:net_pure, kdn_pure=:kdn_pure, nks_pure=:nks_pure,
			stn_amt=:stn_amt, stn_rate=:stn_rate, stn_pure=:stn_pure,
			pure_wt=:pure_wt, discount=:discount, pure_final=:pure_final,
			narr=:narr, narr1=:narr1, btype=:btype, updated_at=NOW()
		WHERE vouch_no=:vouch_no`, sale)
	return err
}

func (r *SaleRepo) Delete(vouchNo string) error {
	_, err := r.db.Exec(`DELETE FROM sale WHERE vouch_no = $1`, vouchNo)
	return err
}

// ─── Purchase ───────────────────────────────────────────────────────────────

type PurchaseRepo struct {
	db *sqlx.DB
}

func NewPurchaseRepo(db *sqlx.DB) *PurchaseRepo {
	return &PurchaseRepo{db: db}
}

func (r *PurchaseRepo) List(page, perPage int, acCode string) ([]models.Purchase, int, error) {
	offset := (page - 1) * perPage
	where := "WHERE 1=1"
	args := []any{}
	i := 1

	if acCode != "" {
		where += fmt.Sprintf(" AND ac_code = $%d", i)
		args = append(args, acCode)
		i++
	}

	var total int
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM purch "+where, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, perPage, offset)
	query := fmt.Sprintf(`SELECT * FROM purch %s ORDER BY vouch_date DESC LIMIT $%d OFFSET $%d`, where, i, i+1)
	var purchases []models.Purchase
	if err := r.db.Select(&purchases, query, args...); err != nil {
		return nil, 0, err
	}
	return purchases, total, nil
}

func (r *PurchaseRepo) Get(vouchNo string) (*models.Purchase, error) {
	var p models.Purchase
	err := r.db.Get(&p, `SELECT * FROM purch WHERE vouch_no = $1`, vouchNo)
	return &p, err
}

func (r *PurchaseRepo) Create(p *models.Purchase) error {
	_, err := r.db.NamedExec(`
		INSERT INTO purch (vouch_no, pre, vouch_no1, vouch_date, ac_code, name, add1, add2, rate, gross_wt, net_wt, bill_amt, narr, rmk, tag, post)
		VALUES (:vouch_no, :pre, :vouch_no1, :vouch_date, :ac_code, :name, :add1, :add2, :rate, :gross_wt, :net_wt, :bill_amt, :narr, :rmk, :tag, :post)`, p)
	return err
}

func (r *PurchaseRepo) Update(p *models.Purchase) error {
	_, err := r.db.NamedExec(`
		UPDATE purch SET
			pre=:pre, vouch_no1=:vouch_no1, vouch_date=:vouch_date,
			ac_code=:ac_code, name=:name, add1=:add1, add2=:add2,
			rate=:rate, gross_wt=:gross_wt, net_wt=:net_wt, bill_amt=:bill_amt,
			narr=:narr, rmk=:rmk, tag=:tag, post=:post
		WHERE vouch_no=:vouch_no`, p)
	return err
}

func (r *PurchaseRepo) Delete(vouchNo string) error {
	_, err := r.db.Exec(`DELETE FROM purch WHERE vouch_no = $1`, vouchNo)
	return err
}
