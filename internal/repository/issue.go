package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/val/inventory/internal/models"
)

type IssueRepo struct {
	db *sqlx.DB
}

func NewIssueRepo(db *sqlx.DB) *IssueRepo {
	return &IssueRepo{db: db}
}

func (r *IssueRepo) List(page, perPage int, acCode, tag string) ([]models.Issue, int, error) {
	offset := (page - 1) * perPage
	where := "WHERE 1=1"
	args := []any{}
	i := 1

	if acCode != "" {
		where += fmt.Sprintf(" AND i.ac_code = $%d", i)
		args = append(args, acCode)
		i++
	}
	if tag != "" {
		where += fmt.Sprintf(" AND i.tag = $%d", i)
		args = append(args, tag)
		i++
	}

	var total int
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM issue i "+where, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, perPage, offset)
	query := fmt.Sprintf(`
		SELECT i.*, f."desc" AS party_name
		FROM issue i
		LEFT JOIN famst f ON f.ac_code = i.ac_code
		%s ORDER BY i.t_date DESC LIMIT $%d OFFSET $%d`, where, i, i+1)

	issues := []models.Issue{}
	if err := r.db.Select(&issues, query, args...); err != nil {
		return nil, 0, err
	}
	return issues, total, nil
}

func (r *IssueRepo) Get(tNo string) (*models.Issue, error) {
	var issue models.Issue
	err := r.db.Get(&issue, `
		SELECT i.*, f."desc" AS party_name
		FROM issue i
		LEFT JOIN famst f ON f.ac_code = i.ac_code
		WHERE i.t_no = $1`, tNo)
	if err != nil {
		return nil, err
	}

	details := []models.IssueDetail{}
	if err := r.db.Select(&details, `SELECT * FROM isdetl WHERE t_no = $1 ORDER BY id`, tNo); err != nil {
		return nil, err
	}
	issue.Details = details
	return &issue, nil
}

func (r *IssueRepo) Create(issue *models.Issue, details []models.IssueDetail) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(`
		INSERT INTO issue (
			t_no, t_no1, pre, t_date, due_date, ac_code, ac_code1, tag,
			nos, gross_wt, kundn_wt, nakash_wt, stone_wt, net_wt,
			net_per, kdn_per, nks_per, net_pure, kdn_pure, nks_pure,
			stn_amt, stn_rate, stn_pure, pure_wt, discount, pure_final,
			kundan_wt, lot_wt, rate, narr, narr1, "desc", bill_amt, tc, exb
		) VALUES (
			:t_no, :t_no1, :pre, :t_date, :due_date, :ac_code, :ac_code1, :tag,
			:nos, :gross_wt, :kundn_wt, :nakash_wt, :stone_wt, :net_wt,
			:net_per, :kdn_per, :nks_per, :net_pure, :kdn_pure, :nks_pure,
			:stn_amt, :stn_rate, :stn_pure, :pure_wt, :discount, :pure_final,
			:kundan_wt, :lot_wt, :rate, :narr, :narr1, :desc, :bill_amt, :tc, :exb
		)`, issue)
	if err != nil {
		return err
	}

	for _, d := range details {
		d.TNo = issue.TNo
		_, err = tx.NamedExec(`
			INSERT INTO isdetl (
				t_no, t_date, recpt_no, recpt_date, itcd, ac_code, itdesc,
				nos, gross_wt, kundn_wt, nakash_wt, stone_wt, net_wt,
				cna, ghat_wt, kundan_wt, extra_wt,
				cat, narr, narr1, narr2, exb
			) VALUES (
				:t_no, :t_date, :recpt_no, :recpt_date, :itcd, :ac_code, :itdesc,
				:nos, :gross_wt, :kundn_wt, :nakash_wt, :stone_wt, :net_wt,
				:cna, :ghat_wt, :kundan_wt, :extra_wt,
				:cat, :narr, :narr1, :narr2, :exb
			)`, &d)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *IssueRepo) Update(issue *models.Issue, details []models.IssueDetail) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NamedExec(`
		UPDATE issue SET
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
			tc=:tc, exb=:exb, updated_at=NOW()
		WHERE t_no=:t_no`, issue)
	if err != nil {
		return err
	}

	if _, err = tx.Exec(`DELETE FROM isdetl WHERE t_no = $1`, issue.TNo); err != nil {
		return err
	}

	for _, d := range details {
		d.TNo = issue.TNo
		_, err = tx.NamedExec(`
			INSERT INTO isdetl (
				t_no, t_date, recpt_no, recpt_date, itcd, ac_code, itdesc,
				nos, gross_wt, kundn_wt, nakash_wt, stone_wt, net_wt,
				cna, ghat_wt, kundan_wt, extra_wt, cat, narr, narr1, narr2, exb
			) VALUES (
				:t_no, :t_date, :recpt_no, :recpt_date, :itcd, :ac_code, :itdesc,
				:nos, :gross_wt, :kundn_wt, :nakash_wt, :stone_wt, :net_wt,
				:cna, :ghat_wt, :kundan_wt, :extra_wt, :cat, :narr, :narr1, :narr2, :exb
			)`, &d)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *IssueRepo) Delete(tNo string) error {
	_, err := r.db.Exec(`DELETE FROM issue WHERE t_no = $1`, tNo)
	return err
}
