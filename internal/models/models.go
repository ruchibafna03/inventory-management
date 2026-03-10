package models

import "time"

// ─── Account / Party ────────────────────────────────────────────────────────

type GroupMaster struct {
	GCode string `db:"g_code" json:"g_code"`
	GDesc string `db:"g_desc" json:"g_desc"`
}

type AccountMaster struct {
	ACCode  string  `db:"ac_code"  json:"ac_code"`
	SubCode *string `db:"sub_code" json:"sub_code,omitempty"`
	Desc    string  `db:"desc"     json:"desc"`
	AName   *string `db:"a_name"   json:"a_name,omitempty"`
	Amt     float64 `db:"amt"      json:"amt"`
	Rmk     *string `db:"rmk"      json:"rmk,omitempty"`
	Rate    float64 `db:"rate"     json:"rate"`
	Per     float64 `db:"per"      json:"per"`
	Opb     float64 `db:"opb"      json:"opb"`
	Date    *time.Time `db:"date"  json:"date,omitempty"`
	Cat     *string `db:"cat"      json:"cat,omitempty"`
	Dia     *string `db:"dia"      json:"dia,omitempty"`
	GCode   *string `db:"g_code"   json:"g_code,omitempty"`
}

type AddressMaster struct {
	ACCode string  `db:"ac_code" json:"ac_code"`
	Add1   *string `db:"add1"    json:"add1,omitempty"`
	Add2   *string `db:"add2"    json:"add2,omitempty"`
	Add3   *string `db:"add3"    json:"add3,omitempty"`
	Pin    *string `db:"pin"     json:"pin,omitempty"`
	TelR1  *string `db:"tel_r1"  json:"tel_r1,omitempty"`
	TelO1  *string `db:"tel_o1"  json:"tel_o1,omitempty"`
	Mobile *string `db:"mobile"  json:"mobile,omitempty"`
	Lst    *string `db:"lst"     json:"lst,omitempty"`
	Cst    *string `db:"cst"     json:"cst,omitempty"`
	PanNo  *string `db:"panno"   json:"panno,omitempty"`
	TinNo  *string `db:"tinno"   json:"tinno,omitempty"`
}

// ─── Rates ──────────────────────────────────────────────────────────────────

type Rate struct {
	ID    int       `db:"id"     json:"id"`
	Date  time.Time `db:"date"   json:"date"`
	Rate  float64   `db:"rate"   json:"rate"`
	Rate1 float64   `db:"rate1"  json:"rate1"`
	SRate float64   `db:"s_rate" json:"s_rate"`
	GTag  string    `db:"gtag"   json:"gtag"`
	QTag  string    `db:"qtag"   json:"qtag"`
	XTag  string    `db:"xtag"   json:"xtag"`
	Sprn  string    `db:"sprn"   json:"sprn"`
}

// ─── Item ───────────────────────────────────────────────────────────────────

type Item struct {
	ItCd       string     `db:"itcd"       json:"itcd"`
	SetCode    *string    `db:"set_code"   json:"set_code,omitempty"`
	Desc       string     `db:"desc"       json:"desc"`
	Fac        float64    `db:"fac"        json:"fac"`
	Stamp      *string    `db:"stamp"      json:"stamp,omitempty"`
	Purity     *string    `db:"purity"     json:"purity,omitempty"`
	Melt       float64    `db:"melt"       json:"melt"`
	GrossWt    float64    `db:"gross_wt"   json:"gross_wt"`
	NetWt      float64    `db:"net_wt"     json:"net_wt"`
	GhatWt     float64    `db:"ghat_wt"    json:"ghat_wt"`
	KundanWt   float64    `db:"kundan_wt"  json:"kundan_wt"`
	ExtraWt    float64    `db:"extra_wt"   json:"extra_wt"`
	GoldWt     float64    `db:"gold_wt"    json:"gold_wt"`
	KundnWt    float64    `db:"kundn_wt"   json:"kundn_wt"`
	NakashWt   float64    `db:"nakash_wt"  json:"nakash_wt"`
	StoneWt    float64    `db:"stone_wt"   json:"stone_wt"`
	PrlWt      float64    `db:"prl_wt"     json:"prl_wt"`
	RbyCt      float64    `db:"rby_ct"     json:"rby_ct"`
	RbyGm      float64    `db:"rby_gm"     json:"rby_gm"`
	EmeCt      float64    `db:"eme_ct"     json:"eme_ct"`
	EmeGm      float64    `db:"eme_gm"     json:"eme_gm"`
	PlkCt      float64    `db:"plk_ct"     json:"plk_ct"`
	PlkGm      float64    `db:"plk_gm"     json:"plk_gm"`
	Stn1Name   *string    `db:"stn1_name"  json:"stn1_name,omitempty"`
	Stn1Ct     float64    `db:"stn1_ct"    json:"stn1_ct"`
	Stn1Gm     float64    `db:"stn1_gm"    json:"stn1_gm"`
	Stn2Name   *string    `db:"stn2_name"  json:"stn2_name,omitempty"`
	Stn2Ct     float64    `db:"stn2_ct"    json:"stn2_ct"`
	Stn2Gm     float64    `db:"stn2_gm"    json:"stn2_gm"`
	Stn3Name   *string    `db:"stn3_name"  json:"stn3_name,omitempty"`
	Stn3Ct     float64    `db:"stn3_ct"    json:"stn3_ct"`
	Stn3Gm     float64    `db:"stn3_gm"    json:"stn3_gm"`
	PrlsGm     float64    `db:"prls_gm"    json:"prls_gm"`
	PrlbGm     float64    `db:"prlb_gm"    json:"prlb_gm"`
	RbybGm     float64    `db:"rbyb_gm"    json:"rbyb_gm"`
	EmebGm     float64    `db:"emeb_gm"    json:"emeb_gm"`
	EmebGm2    float64    `db:"emebb_gm"   json:"emebb_gm"`
	Prl1Name   *string    `db:"prl1_name"  json:"prl1_name,omitempty"`
	Prl1Gm     float64    `db:"prl1_gm"    json:"prl1_gm"`
	SrbyGm     float64    `db:"srby_gm"    json:"srby_gm"`
	SemeGm     float64    `db:"seme_gm"    json:"seme_gm"`
	SczGm      float64    `db:"scz_gm"     json:"scz_gm"`
	SnavGm     float64    `db:"snav_gm"    json:"snav_gm"`
	Prl2Name   *string    `db:"prl2_name"  json:"prl2_name,omitempty"`
	Prl2Gm     float64    `db:"prl2_gm"    json:"prl2_gm"`
	Prl3Name   *string    `db:"prl3_name"  json:"prl3_name,omitempty"`
	Prl3Gm     float64    `db:"prl3_gm"    json:"prl3_gm"`
	Narr1      *string    `db:"narr1"      json:"narr1,omitempty"`
	Narr2      *string    `db:"narr2"      json:"narr2,omitempty"`
	MkChrg     float64    `db:"mk_chrg"    json:"mk_chrg"`
	WPer       float64    `db:"w_per"      json:"w_per"`
	RecptDate  *time.Time `db:"recpt_date" json:"recpt_date,omitempty"`
	LotNo      *string    `db:"lot_no"     json:"lot_no,omitempty"`
	IssueNo    *string    `db:"issue_no"   json:"issue_no,omitempty"`
	IssueDate  *time.Time `db:"issue_date" json:"issue_date,omitempty"`
	ACCode     *string    `db:"ac_code"    json:"ac_code,omitempty"`
	StkCode    *string    `db:"stk_code"   json:"stk_code,omitempty"`
	Cat        *string    `db:"cat"        json:"cat,omitempty"`
	Kp         *string    `db:"kp"         json:"kp,omitempty"`
	Tag        *string    `db:"tag"        json:"tag,omitempty"`
	Cat1       *string    `db:"cat1"       json:"cat1,omitempty"`
	Cna        *string    `db:"cna"        json:"cna,omitempty"`
	Exb        *string    `db:"exb"        json:"exb,omitempty"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}

// ─── Lot ────────────────────────────────────────────────────────────────────

type Lot struct {
	TNo       string     `db:"t_no"        json:"t_no"`
	TDate     time.Time  `db:"t_date"      json:"t_date"`
	ACCode    *string    `db:"ac_code"     json:"ac_code,omitempty"`
	Party     *string    `db:"party"       json:"party,omitempty"`
	LotNo     string     `db:"lot_no"      json:"lot_no"`
	Nos       int        `db:"nos"         json:"nos"`
	GrossWt   float64    `db:"gross_wt"    json:"gross_wt"`
	KundanWt  float64    `db:"kundan_wt"   json:"kundan_wt"`
	NakashWt  float64    `db:"nakash_wt"   json:"nakash_wt"`
	StoneWt   float64    `db:"stone_wt"    json:"stone_wt"`
	PearlWt   float64    `db:"pearl_wt"    json:"pearl_wt"`
	Nos1      int        `db:"nos1"        json:"nos1"`
	GrossWt1  float64    `db:"gross_wt1"   json:"gross_wt1"`
	AdjWt     float64    `db:"adj_wt"      json:"adj_wt"`
	KundanWt1 float64    `db:"kundan_wt1"  json:"kundan_wt1"`
	NakashWt1 float64    `db:"nakash_wt1"  json:"nakash_wt1"`
	StoneWt1  float64    `db:"stone_wt1"   json:"stone_wt1"`
	PearlWt1  float64    `db:"pearl_wt1"   json:"pearl_wt1"`
	Amt       float64    `db:"amt"         json:"amt"`
	GRt       float64    `db:"g_rt"        json:"g_rt"`
	Desc      *string    `db:"desc"        json:"desc,omitempty"`
}

// ─── Issue ──────────────────────────────────────────────────────────────────

type Issue struct {
	TNo       string     `db:"t_no"        json:"t_no"`
	TNo1      *string    `db:"t_no1"       json:"t_no1,omitempty"`
	Pre       *string    `db:"pre"         json:"pre,omitempty"`
	TDate     time.Time  `db:"t_date"      json:"t_date"`
	DueDate   *time.Time `db:"due_date"    json:"due_date,omitempty"`
	ACCode    *string    `db:"ac_code"     json:"ac_code,omitempty"`
	ACCode1   *string    `db:"ac_code1"    json:"ac_code1,omitempty"`
	Tag       *string    `db:"tag"         json:"tag,omitempty"`
	Nos       float64    `db:"nos"         json:"nos"`
	GrossWt   float64    `db:"gross_wt"    json:"gross_wt"`
	KundnWt   float64    `db:"kundn_wt"    json:"kundn_wt"`
	NakashWt  float64    `db:"nakash_wt"   json:"nakash_wt"`
	StoneWt   float64    `db:"stone_wt"    json:"stone_wt"`
	NetWt     float64    `db:"net_wt"      json:"net_wt"`
	NetPer    float64    `db:"net_per"     json:"net_per"`
	KdnPer    float64    `db:"kdn_per"     json:"kdn_per"`
	NksPer    float64    `db:"nks_per"     json:"nks_per"`
	NetPure   float64    `db:"net_pure"    json:"net_pure"`
	KdnPure   float64    `db:"kdn_pure"    json:"kdn_pure"`
	NksPure   float64    `db:"nks_pure"    json:"nks_pure"`
	StnAmt    float64    `db:"stn_amt"     json:"stn_amt"`
	StnRate   float64    `db:"stn_rate"    json:"stn_rate"`
	StnPure   float64    `db:"stn_pure"    json:"stn_pure"`
	PureWt    float64    `db:"pure_wt"     json:"pure_wt"`
	Discount  float64    `db:"discount"    json:"discount"`
	PureFinal float64    `db:"pure_final"  json:"pure_final"`
	KundanWt  float64    `db:"kundan_wt"   json:"kundan_wt"`
	LotWt     float64    `db:"lot_wt"      json:"lot_wt"`
	Rate      float64    `db:"rate"        json:"rate"`
	Narr      *string    `db:"narr"        json:"narr,omitempty"`
	Narr1     *string    `db:"narr1"       json:"narr1,omitempty"`
	Desc      *string    `db:"desc"        json:"desc,omitempty"`
	BillAmt   float64    `db:"bill_amt"    json:"bill_amt"`
	Tc        *string    `db:"tc"          json:"tc,omitempty"`
	Exb       *string    `db:"exb"         json:"exb,omitempty"`
	CreatedAt time.Time  `db:"created_at"  json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"  json:"updated_at"`
	Details   []IssueDetail `db:"-"        json:"details,omitempty"`
	PartyName *string    `db:"party_name"  json:"party_name,omitempty"`
}

type IssueDetail struct {
	ID        int        `db:"id"         json:"id"`
	TNo       string     `db:"t_no"       json:"t_no"`
	TDate     *time.Time `db:"t_date"     json:"t_date,omitempty"`
	RecptNo   *string    `db:"recpt_no"   json:"recpt_no,omitempty"`
	RecptDate *time.Time `db:"recpt_date" json:"recpt_date,omitempty"`
	ItCd      *string    `db:"itcd"       json:"itcd,omitempty"`
	ACCode    *string    `db:"ac_code"    json:"ac_code,omitempty"`
	ItDesc    *string    `db:"itdesc"     json:"itdesc,omitempty"`
	Nos       float64    `db:"nos"        json:"nos"`
	GrossWt   float64    `db:"gross_wt"   json:"gross_wt"`
	KundnWt   float64    `db:"kundn_wt"   json:"kundn_wt"`
	NakashWt  float64    `db:"nakash_wt"  json:"nakash_wt"`
	StoneWt   float64    `db:"stone_wt"   json:"stone_wt"`
	NetWt     float64    `db:"net_wt"     json:"net_wt"`
	Cna       *string    `db:"cna"        json:"cna,omitempty"`
	GhatWt    float64    `db:"ghat_wt"    json:"ghat_wt"`
	KundanWt  float64    `db:"kundan_wt"  json:"kundan_wt"`
	ExtraWt   float64    `db:"extra_wt"   json:"extra_wt"`
	Cat       *string    `db:"cat"        json:"cat,omitempty"`
	Narr      *string    `db:"narr"       json:"narr,omitempty"`
	Narr1     *string    `db:"narr1"      json:"narr1,omitempty"`
	Narr2     *string    `db:"narr2"      json:"narr2,omitempty"`
	Exb       *string    `db:"exb"        json:"exb,omitempty"`
}

// ─── Receipt ────────────────────────────────────────────────────────────────

type Receipt struct {
	TNo       string     `db:"t_no"        json:"t_no"`
	TNo1      *string    `db:"t_no1"       json:"t_no1,omitempty"`
	Pre       *string    `db:"pre"         json:"pre,omitempty"`
	TDate     time.Time  `db:"t_date"      json:"t_date"`
	DueDate   *time.Time `db:"due_date"    json:"due_date,omitempty"`
	ACCode    *string    `db:"ac_code"     json:"ac_code,omitempty"`
	ACCode1   *string    `db:"ac_code1"    json:"ac_code1,omitempty"`
	Tag       *string    `db:"tag"         json:"tag,omitempty"`
	Nos       float64    `db:"nos"         json:"nos"`
	GrossWt   float64    `db:"gross_wt"    json:"gross_wt"`
	KundnWt   float64    `db:"kundn_wt"    json:"kundn_wt"`
	NakashWt  float64    `db:"nakash_wt"   json:"nakash_wt"`
	StoneWt   float64    `db:"stone_wt"    json:"stone_wt"`
	NetWt     float64    `db:"net_wt"      json:"net_wt"`
	NetPer    float64    `db:"net_per"     json:"net_per"`
	KdnPer    float64    `db:"kdn_per"     json:"kdn_per"`
	NksPer    float64    `db:"nks_per"     json:"nks_per"`
	NetPure   float64    `db:"net_pure"    json:"net_pure"`
	KdnPure   float64    `db:"kdn_pure"    json:"kdn_pure"`
	NksPure   float64    `db:"nks_pure"    json:"nks_pure"`
	StnAmt    float64    `db:"stn_amt"     json:"stn_amt"`
	StnRate   float64    `db:"stn_rate"    json:"stn_rate"`
	StnPure   float64    `db:"stn_pure"    json:"stn_pure"`
	PureWt    float64    `db:"pure_wt"     json:"pure_wt"`
	Discount  float64    `db:"discount"    json:"discount"`
	PureFinal float64    `db:"pure_final"  json:"pure_final"`
	KundanWt  float64    `db:"kundan_wt"   json:"kundan_wt"`
	LotWt     float64    `db:"lot_wt"      json:"lot_wt"`
	Rate      float64    `db:"rate"        json:"rate"`
	Narr      *string    `db:"narr"        json:"narr,omitempty"`
	Narr1     *string    `db:"narr1"       json:"narr1,omitempty"`
	Desc      *string    `db:"desc"        json:"desc,omitempty"`
	BillAmt   float64    `db:"bill_amt"    json:"bill_amt"`
	Tc        *string    `db:"tc"          json:"tc,omitempty"`
	CreatedAt time.Time  `db:"created_at"  json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"  json:"updated_at"`
	Details   []ReceiptDetail `db:"-"       json:"details,omitempty"`
	PartyName *string    `db:"party_name"  json:"party_name,omitempty"`
}

type ReceiptDetail struct {
	ID       int        `db:"id"        json:"id"`
	TNo      string     `db:"t_no"      json:"t_no"`
	TDate    *time.Time `db:"t_date"    json:"t_date,omitempty"`
	ACCode   *string    `db:"ac_code"   json:"ac_code,omitempty"`
	ItCd     *string    `db:"itcd"      json:"itcd,omitempty"`
	Fac      int        `db:"fac"       json:"fac"`
	ItDesc   *string    `db:"itdesc"    json:"itdesc,omitempty"`
	LotNo    *string    `db:"lot_no"    json:"lot_no,omitempty"`
	GrossWt  float64    `db:"gross_wt"  json:"gross_wt"`
	NetWt    float64    `db:"net_wt"    json:"net_wt"`
	GhatWt   float64    `db:"ghat_wt"   json:"ghat_wt"`
	KundanWt float64    `db:"kundan_wt" json:"kundan_wt"`
	ExtraWt  float64    `db:"extra_wt"  json:"extra_wt"`
	GoldWt   float64    `db:"gold_wt"   json:"gold_wt"`
	KundnWt  float64    `db:"kundn_wt"  json:"kundn_wt"`
	NakashWt float64    `db:"nakash_wt" json:"nakash_wt"`
	StoneWt  float64    `db:"stone_wt"  json:"stone_wt"`
	PrlWt    float64    `db:"prl_wt"    json:"prl_wt"`
	BillAmt  float64    `db:"bill_amt"  json:"bill_amt"`
	Cat      *string    `db:"cat"       json:"cat,omitempty"`
	Narr1    *string    `db:"narr1"     json:"narr1,omitempty"`
	Narr2    *string    `db:"narr2"     json:"narr2,omitempty"`
}

// ─── Sale ───────────────────────────────────────────────────────────────────

type Sale struct {
	VouchNo   string     `db:"vouch_no"    json:"vouch_no"`
	OvNo      *string    `db:"ovno"        json:"ovno,omitempty"`
	Pre       *string    `db:"pre"         json:"pre,omitempty"`
	VouchNo1  *string    `db:"vouch_no1"   json:"vouch_no1,omitempty"`
	VouchDate time.Time  `db:"vouch_date"  json:"vouch_date"`
	ACCode    *string    `db:"ac_code"     json:"ac_code,omitempty"`
	Name      *string    `db:"name"        json:"name,omitempty"`
	Add1      *string    `db:"add1"        json:"add1,omitempty"`
	Add2      *string    `db:"add2"        json:"add2,omitempty"`
	PhoneNo   *string    `db:"phone_no"    json:"phone_no,omitempty"`
	PhoneNo1  *string    `db:"phone_no1"   json:"phone_no1,omitempty"`
	Sman      *string    `db:"sman"        json:"sman,omitempty"`
	Nos       int        `db:"nos"         json:"nos"`
	GrossWt   float64    `db:"gross_wt"    json:"gross_wt"`
	KundnWt   float64    `db:"kundn_wt"    json:"kundn_wt"`
	NakashWt  float64    `db:"nakash_wt"   json:"nakash_wt"`
	StoneWt   float64    `db:"stone_wt"    json:"stone_wt"`
	NetWt     float64    `db:"net_wt"      json:"net_wt"`
	NetPer    float64    `db:"net_per"     json:"net_per"`
	KdnPer    float64    `db:"kdn_per"     json:"kdn_per"`
	NksPer    float64    `db:"nks_per"     json:"nks_per"`
	NetPure   float64    `db:"net_pure"    json:"net_pure"`
	StnAmt    float64    `db:"stn_amt"     json:"stn_amt"`
	StnRate   float64    `db:"stn_rate"    json:"stn_rate"`
	PureWt    float64    `db:"pure_wt"     json:"pure_wt"`
	Discount  float64    `db:"discount"    json:"discount"`
	PureFinal float64    `db:"pure_final"  json:"pure_final"`
	Narr      *string    `db:"narr"        json:"narr,omitempty"`
	Narr1     *string    `db:"narr1"       json:"narr1,omitempty"`
	BType     *string    `db:"btype"       json:"btype,omitempty"`
	CreatedAt time.Time  `db:"created_at"  json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"  json:"updated_at"`
}

// ─── Purchase ───────────────────────────────────────────────────────────────

type Purchase struct {
	VouchNo   string    `db:"vouch_no"   json:"vouch_no"`
	Pre       *string   `db:"pre"        json:"pre,omitempty"`
	VouchNo1  *string   `db:"vouch_no1"  json:"vouch_no1,omitempty"`
	VouchDate time.Time `db:"vouch_date" json:"vouch_date"`
	ACCode    *string   `db:"ac_code"    json:"ac_code,omitempty"`
	Name      *string   `db:"name"       json:"name,omitempty"`
	Add1      *string   `db:"add1"       json:"add1,omitempty"`
	Add2      *string   `db:"add2"       json:"add2,omitempty"`
	Rate      float64   `db:"rate"       json:"rate"`
	GrossWt   float64   `db:"gross_wt"   json:"gross_wt"`
	NetWt     float64   `db:"net_wt"     json:"net_wt"`
	BillAmt   float64   `db:"bill_amt"   json:"bill_amt"`
	Narr      *string   `db:"narr"       json:"narr,omitempty"`
	Rmk       *string   `db:"rmk"        json:"rmk,omitempty"`
	Tag       *string   `db:"tag"        json:"tag,omitempty"`
	Post      *string   `db:"post"       json:"post,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// ─── User ───────────────────────────────────────────────────────────────────

type User struct {
	ID           int       `db:"id"            json:"id"`
	Username     string    `db:"username"      json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	FullName     string    `db:"full_name"     json:"full_name"`
	Role         string    `db:"role"          json:"role"`
	Active       bool      `db:"active"        json:"active"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
}

// ─── Pagination ─────────────────────────────────────────────────────────────

type PagedResult[T any] struct {
	Data    []T `json:"data"`
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
