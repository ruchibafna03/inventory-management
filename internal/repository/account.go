package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/val/inventory/internal/models"
)

type AccountRepo struct {
	db *sqlx.DB
}

func NewAccountRepo(db *sqlx.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func (r *AccountRepo) ListGroups() ([]models.GroupMaster, error) {
	var groups []models.GroupMaster
	err := r.db.Select(&groups, `SELECT * FROM grpmst ORDER BY g_code`)
	return groups, err
}

func (r *AccountRepo) CreateGroup(g *models.GroupMaster) error {
	_, err := r.db.NamedExec(`INSERT INTO grpmst (g_code, g_desc) VALUES (:g_code, :g_desc)`, g)
	return err
}

func (r *AccountRepo) UpdateGroup(g *models.GroupMaster) error {
	_, err := r.db.NamedExec(`UPDATE grpmst SET g_desc=:g_desc WHERE g_code=:g_code`, g)
	return err
}

func (r *AccountRepo) DeleteGroup(gCode string) error {
	_, err := r.db.Exec(`DELETE FROM grpmst WHERE g_code = $1`, gCode)
	return err
}

func (r *AccountRepo) List(page, perPage int, search, gCode string) ([]models.AccountMaster, int, error) {
	offset := (page - 1) * perPage
	where := "WHERE 1=1"
	args := []any{}
	i := 1

	if search != "" {
		where += fmt.Sprintf(` AND (ac_code ILIKE $%d OR "desc" ILIKE $%d)`, i, i+1)
		args = append(args, "%"+search+"%", "%"+search+"%")
		i += 2
	}
	if gCode != "" {
		where += fmt.Sprintf(" AND g_code = $%d", i)
		args = append(args, gCode)
		i++
	}

	var total int
	if err := r.db.Get(&total, "SELECT COUNT(*) FROM famst "+where, args...); err != nil {
		return nil, 0, err
	}

	args = append(args, perPage, offset)
	query := fmt.Sprintf(`SELECT * FROM famst %s ORDER BY ac_code LIMIT $%d OFFSET $%d`, where, i, i+1)
	var accounts []models.AccountMaster
	if err := r.db.Select(&accounts, query, args...); err != nil {
		return nil, 0, err
	}
	return accounts, total, nil
}

func (r *AccountRepo) Get(acCode string) (*models.AccountMaster, error) {
	var acc models.AccountMaster
	err := r.db.Get(&acc, `SELECT * FROM famst WHERE ac_code = $1`, acCode)
	return &acc, err
}

func (r *AccountRepo) GetAddress(acCode string) (*models.AddressMaster, error) {
	var addr models.AddressMaster
	err := r.db.Get(&addr, `SELECT * FROM add_master WHERE ac_code = $1`, acCode)
	return &addr, err
}

func (r *AccountRepo) Create(acc *models.AccountMaster) error {
	_, err := r.db.NamedExec(`
		INSERT INTO famst (ac_code, sub_code, "desc", a_name, amt, rmk, rate, per, opb, "date", cat, dia, g_code)
		VALUES (:ac_code, :sub_code, :desc, :a_name, :amt, :rmk, :rate, :per, :opb, :date, :cat, :dia, :g_code)`, acc)
	return err
}

func (r *AccountRepo) Update(acc *models.AccountMaster) error {
	_, err := r.db.NamedExec(`
		UPDATE famst SET
			sub_code=:sub_code, "desc"=:desc, a_name=:a_name, amt=:amt, rmk=:rmk,
			rate=:rate, per=:per, opb=:opb, "date"=:date, cat=:cat, dia=:dia, g_code=:g_code
		WHERE ac_code=:ac_code`, acc)
	return err
}

func (r *AccountRepo) UpsertAddress(addr *models.AddressMaster) error {
	_, err := r.db.NamedExec(`
		INSERT INTO add_master (ac_code, add1, add2, add3, pin, tel_r1, tel_o1, mobile, lst, cst, panno, tinno)
		VALUES (:ac_code, :add1, :add2, :add3, :pin, :tel_r1, :tel_o1, :mobile, :lst, :cst, :panno, :tinno)
		ON CONFLICT (ac_code) DO UPDATE SET
			add1=EXCLUDED.add1, add2=EXCLUDED.add2, add3=EXCLUDED.add3,
			pin=EXCLUDED.pin, tel_r1=EXCLUDED.tel_r1, tel_o1=EXCLUDED.tel_o1,
			mobile=EXCLUDED.mobile, lst=EXCLUDED.lst, cst=EXCLUDED.cst,
			panno=EXCLUDED.panno, tinno=EXCLUDED.tinno`, addr)
	return err
}

func (r *AccountRepo) Delete(acCode string) error {
	_, err := r.db.Exec(`DELETE FROM famst WHERE ac_code = $1`, acCode)
	return err
}

// OutstandingBalances returns account-wise outstanding
func (r *AccountRepo) OutstandingBalances(gCode string) ([]models.AccountMaster, error) {
	where := ""
	args := []any{}
	if gCode != "" {
		where = " WHERE g_code = $1"
		args = append(args, gCode)
	}
	var accounts []models.AccountMaster
	err := r.db.Select(&accounts, `SELECT * FROM famst`+where+` ORDER BY ac_code`, args...)
	return accounts, err
}
