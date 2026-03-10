package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/val/inventory/internal/models"
)

type ItemRepo struct {
	db *sqlx.DB
}

func NewItemRepo(db *sqlx.DB) *ItemRepo {
	return &ItemRepo{db: db}
}

func (r *ItemRepo) List(page, perPage int, search, cat string) ([]models.Item, int, error) {
	offset := (page - 1) * perPage
	where := "WHERE 1=1"
	args := []any{}
	i := 1

	if search != "" {
		where += fmt.Sprintf(` AND (itcd ILIKE $%d OR "desc" ILIKE $%d)`, i, i+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
		i += 2
	}
	if cat != "" {
		where += fmt.Sprintf(" AND cat = $%d", i)
		args = append(args, cat)
		i++
	}

	var total int
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM item "+where, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, perPage, offset)
	query := fmt.Sprintf(`SELECT * FROM item %s ORDER BY itcd LIMIT $%d OFFSET $%d`, where, i, i+1)
	var items []models.Item
	if err := r.db.Select(&items, query, args...); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *ItemRepo) Get(itcd string) (*models.Item, error) {
	var item models.Item
	err := r.db.Get(&item, `SELECT * FROM item WHERE itcd = $1`, itcd)
	return &item, err
}

func (r *ItemRepo) Create(item *models.Item) error {
	_, err := r.db.NamedExec(`
		INSERT INTO item (
			itcd, set_code, "desc", fac, stamp, purity, melt,
			gross_wt, net_wt, ghat_wt, kundan_wt, extra_wt, gold_wt,
			kundn_wt, nakash_wt, stone_wt, prl_wt,
			rby_ct, rby_gm, eme_ct, eme_gm, plk_ct, plk_gm,
			stn1_name, stn1_ct, stn1_gm,
			stn2_name, stn2_ct, stn2_gm,
			stn3_name, stn3_ct, stn3_gm,
			prls_gm, prlb_gm, rbyb_gm, emeb_gm, emebb_gm,
			prl1_name, prl1_gm, srby_gm, seme_gm, scz_gm, snav_gm,
			prl2_name, prl2_gm, prl3_name, prl3_gm,
			narr1, narr2, mk_chrg, w_per,
			recpt_date, lot_no, issue_no, issue_date,
			ac_code, stk_code, cat, kp, tag, cat1, cna, exb
		) VALUES (
			:itcd, :set_code, :desc, :fac, :stamp, :purity, :melt,
			:gross_wt, :net_wt, :ghat_wt, :kundan_wt, :extra_wt, :gold_wt,
			:kundn_wt, :nakash_wt, :stone_wt, :prl_wt,
			:rby_ct, :rby_gm, :eme_ct, :eme_gm, :plk_ct, :plk_gm,
			:stn1_name, :stn1_ct, :stn1_gm,
			:stn2_name, :stn2_ct, :stn2_gm,
			:stn3_name, :stn3_ct, :stn3_gm,
			:prls_gm, :prlb_gm, :rbyb_gm, :emeb_gm, :emebb_gm,
			:prl1_name, :prl1_gm, :srby_gm, :seme_gm, :scz_gm, :snav_gm,
			:prl2_name, :prl2_gm, :prl3_name, :prl3_gm,
			:narr1, :narr2, :mk_chrg, :w_per,
			:recpt_date, :lot_no, :issue_no, :issue_date,
			:ac_code, :stk_code, :cat, :kp, :tag, :cat1, :cna, :exb
		)`, item)
	return err
}

func (r *ItemRepo) Update(item *models.Item) error {
	_, err := r.db.NamedExec(`
		UPDATE item SET
			set_code=:set_code, "desc"=:desc, fac=:fac, stamp=:stamp, purity=:purity, melt=:melt,
			gross_wt=:gross_wt, net_wt=:net_wt, ghat_wt=:ghat_wt, kundan_wt=:kundan_wt,
			extra_wt=:extra_wt, gold_wt=:gold_wt, kundn_wt=:kundn_wt, nakash_wt=:nakash_wt,
			stone_wt=:stone_wt, prl_wt=:prl_wt,
			rby_ct=:rby_ct, rby_gm=:rby_gm, eme_ct=:eme_ct, eme_gm=:eme_gm,
			plk_ct=:plk_ct, plk_gm=:plk_gm,
			stn1_name=:stn1_name, stn1_ct=:stn1_ct, stn1_gm=:stn1_gm,
			stn2_name=:stn2_name, stn2_ct=:stn2_ct, stn2_gm=:stn2_gm,
			stn3_name=:stn3_name, stn3_ct=:stn3_ct, stn3_gm=:stn3_gm,
			prls_gm=:prls_gm, prlb_gm=:prlb_gm, rbyb_gm=:rbyb_gm,
			emeb_gm=:emeb_gm, emebb_gm=:emebb_gm,
			prl1_name=:prl1_name, prl1_gm=:prl1_gm, srby_gm=:srby_gm,
			seme_gm=:seme_gm, scz_gm=:scz_gm, snav_gm=:snav_gm,
			prl2_name=:prl2_name, prl2_gm=:prl2_gm,
			prl3_name=:prl3_name, prl3_gm=:prl3_gm,
			narr1=:narr1, narr2=:narr2, mk_chrg=:mk_chrg, w_per=:w_per,
			recpt_date=:recpt_date, lot_no=:lot_no, issue_no=:issue_no, issue_date=:issue_date,
			ac_code=:ac_code, stk_code=:stk_code, cat=:cat, kp=:kp, tag=:tag,
			cat1=:cat1, cna=:cna, exb=:exb,
			updated_at=NOW()
		WHERE itcd=:itcd`, item)
	return err
}

func (r *ItemRepo) Delete(itcd string) error {
	_, err := r.db.Exec(`DELETE FROM item WHERE itcd = $1`, itcd)
	return err
}

func (r *ItemRepo) Tags() ([]string, error) {
	var tags []string
	err := r.db.Select(&tags, `SELECT itcd FROM itblk ORDER BY itcd`)
	return tags, err
}

func (r *ItemRepo) AddTag(itcd string) error {
	_, err := r.db.Exec(`INSERT INTO itblk (itcd) VALUES ($1) ON CONFLICT DO NOTHING`, itcd)
	return err
}
