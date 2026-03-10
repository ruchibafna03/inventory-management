package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/val/inventory/internal/models"
)

type ReceiptRepo struct {
	db *sqlx.DB
}

func NewReceiptRepo(db *sqlx.DB) *ReceiptRepo {
	return &ReceiptRepo{db: db}
}

func (r *ReceiptRepo) List(page, perPage int, acCode string) ([]models.Receipt, int, error) {
	offset := (page - 1) * perPage
	where := "WHERE 1=1"
	args := []any{}
	i := 1

	if acCode != "" {
		where += fmt.Sprintf(" AND r.ac_code = $%d", i)
		args = append(args, acCode)
		i++
	}

	var total int
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM recpt r "+where, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, perPage, offset)
	query := fmt.Sprintf(`
		SELECT r.*, f."desc" AS party_name
		FROM recpt r
		LEFT JOIN famst f ON f.ac_code = r.ac_code
		%s ORDER BY r.t_date DESC LIMIT $%d OFFSET $%d`, where, i, i+1)

	receipts := []models.Receipt{}
	if err := r.db.Select(&receipts, query, args...); err != nil {
		return nil, 0, err
	}
	return receipts, total, nil
}

func (r *ReceiptRepo) Get(tNo string) (*models.Receipt, error) {
	var recpt models.Receipt
	err := r.db.Get(&recpt, `
		SELECT r.*, f."desc" AS party_name
		FROM recpt r
		LEFT JOIN famst f ON f.ac_code = r.ac_code
		WHERE r.t_no = $1`, tNo)
	if err != nil {
		return nil, err
	}

	details := []models.ReceiptDetail{}
	if err := r.db.Select(&details, `SELECT * FROM rcdetl WHERE t_no = $1 ORDER BY id`, tNo); err != nil {
		return nil, err
	}
	recpt.Details = details
	return &recpt, nil
}

func (r *ReceiptRepo) Create(recpt *models.Receipt, details []models.ReceiptDetail) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(`
		INSERT INTO recpt (
			t_no, t_no1, pre, t_date, due_date, ac_code, ac_code1, tag,
			nos, gross_wt, kundn_wt, nakash_wt, stone_wt, net_wt,
			net_per, kdn_per, nks_per, net_pure, kdn_pure, nks_pure,
			stn_amt, stn_rate, stn_pure, pure_wt, discount, pure_final,
			kundan_wt, lot_wt, rate, narr, narr1, desc, bill_amt, tc
		) VALUES (
			:t_no, :t_no1, :pre, :t_date, :due_date, :ac_code, :ac_code1, :tag,
			:nos, :gross_wt, :kundn_wt, :nakash_wt, :stone_wt, :net_wt,
			:net_per, :kdn_per, :nks_per, :net_pure, :kdn_pure, :nks_pure,
			:stn_amt, :stn_rate, :stn_pure, :pure_wt, :discount, :pure_final,
			:kundan_wt, :lot_wt, :rate, :narr, :narr1, :desc, :bill_amt, :tc
		)`, recpt)
	if err != nil {
		return err
	}

	for _, d := range details {
		d.TNo = recpt.TNo
		_, err = tx.NamedExec(`
			INSERT INTO rcdetl (
				t_no, t_date, ac_code, itcd, fac, itdesc, lot_no,
				gross_wt, net_wt, ghat_wt, kundan_wt, extra_wt, gold_wt,
				kundn_wt, nakash_wt, stone_wt, prl_wt,
				bill_amt, cat, narr1, narr2
			) VALUES (
				:t_no, :t_date, :ac_code, :itcd, :fac, :itdesc, :lot_no,
				:gross_wt, :net_wt, :ghat_wt, :kundan_wt, :extra_wt, :gold_wt,
				:kundn_wt, :nakash_wt, :stone_wt, :prl_wt,
				:bill_amt, :cat, :narr1, :narr2
			)`, &d)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ReceiptRepo) Update(recpt *models.Receipt, details []models.ReceiptDetail) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(`
		UPDATE recpt SET
			t_no1=:t_no1, pre=:pre, t_date=:t_date, due_date=:due_date,
			ac_code=:ac_code, ac_code1=:ac_code1, tag=:tag,
			nos=:nos, gross_wt=:gross_wt, kundn_wt=:kundn_wt,
			nakash_wt=:nakash_wt, stone_wt=:stone_wt, net_wt=:net_wt,
			net_per=:net_per, kdn_per=:kdn_per, nks_per=:nks_per,
			net_pure=:net_pure, kdn_pure=:kdn_pure, nks_pure=:nks_pure,
			stn_amt=:stn_amt, stn_rate=:stn_rate, stn_pure=:stn_pure,
			pure_wt=:pure_wt, discount=:discount, pure_final=:pure_final,
			kundan_wt=:kundan_wt, lot_wt=:lot_wt, rate=:rate,
			narr=:narr, narr1=:narr1, "desc"=:desc, bill_amt=:bill_amt,
			tc=:tc, updated_at=NOW()
		WHERE t_no=:t_no`, recpt)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(`DELETE FROM rcdetl WHERE t_no = $1`, recpt.TNo); err != nil {
		return err
	}

	for _, d := range details {
		d.TNo = recpt.TNo
		_, err = tx.NamedExec(`
			INSERT INTO rcdetl (
				t_no, t_date, ac_code, itcd, fac, itdesc, lot_no,
				gross_wt, net_wt, ghat_wt, kundan_wt, extra_wt, gold_wt,
				kundn_wt, nakash_wt, stone_wt, prl_wt,
				bill_amt, cat, narr1, narr2
			) VALUES (
				:t_no, :t_date, :ac_code, :itcd, :fac, :itdesc, :lot_no,
				:gross_wt, :net_wt, :ghat_wt, :kundan_wt, :extra_wt, :gold_wt,
				:kundn_wt, :nakash_wt, :stone_wt, :prl_wt,
				:bill_amt, :cat, :narr1, :narr2
			)`, &d)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ReceiptRepo) Delete(tNo string) error {
	_, err := r.db.Exec(`DELETE FROM recpt WHERE t_no = $1`, tNo)
	return err
}
