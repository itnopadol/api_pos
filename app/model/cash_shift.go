package model

import (
	"time"
	"github.com/jmoiron/sqlx"
	"fmt"
)

type Shift struct {
	Id int64 `json:"id" db:"id"`
	HostId string `json:"host_id" db:"host_id"`
	DocDate string `json:"doc_date" db:"doc_date"`
	ChangeBegin float64 `json:"change_begin" db:"change_begin"`
	ChangeAmount float64 `json:"change_amount" db:"change_amount"`
	CashAmount float64 `json:"cash_amount" db:"cash_amount"`
	ExpensesAmount float64 `json:"expenses_amount" db:"expenses_amount"`
	IsClosed int `json:"is_closed" db:"is_closed"`
	CreatedBy string `json:"created_by" db:"created_by"`
	Created time.Time `json:"created" db:"created"`
	EditedBy string `json:"edited_by" db:"edited_by"`
	Edited time.Time `json:"edited" db:"edited"`
	ClosedBy string `json:"closed_by" db:"closed_by"`
	Closed time.Time `json:"closed" db:"closed"`
}


func (ch *Shift)SaveShift(db *sqlx.DB) error{
	now := time.Now()
	fmt.Println("yyyy-mm-dd date format : ", now.AddDate(0, 0, 0).Format("2006-01-02"))
	DocDate := now.AddDate(0, 0, 0).Format("2006-01-02")

	ch.DocDate = DocDate
	ch.Created = now

	var checkCount int
	sqlCheckExist := `select count(host_id) as vCount from cash_shift where doc_date = ?`
	err := db.Get(&checkCount, sqlCheckExist, ch.DocDate)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if (checkCount == 0){
		sql := `INSERT INTO cash_shift(host_id,doc_date,change_begin,change_amount,cash_amount,expenses_amount,created_by,created) VALUES(?,?,?,?,?,?,?,?)`
		res, err := db.Exec(sql, ch.HostId, ch.DocDate, ch.ChangeBegin, ch.ChangeAmount, ch.CashAmount, ch.ExpensesAmount, ch.CreatedBy, ch.Created)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		Id, _ := res.LastInsertId()
		ch.Id = Id
	}else{

	}

	return nil
}

func (ch *Shift)UpdateShift(db *sqlx.DB) error{
	now := time.Now()
	fmt.Println("yyyy-mm-dd date format : ", now.AddDate(0, 0, 0).Format("2006-01-02"))

	ch.Edited = now

	var checkCount int
	sqlCheckExist := `select count(host_id) as vCount from cash_shift where doc_date = ?`
	err := db.Get(&checkCount, sqlCheckExist, ch.DocDate)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if (checkCount != 0) {
		sql := `UPDATE cash_shift set change_begin = ?, change_amount = ?,cash_amount = ?,expenses_amount = ?, edited_by = ?, edited = ? where  host_id = ? and doc_date  = ?`
		_, err = db.Exec(sql,ch.ChangeBegin, ch.ChangeAmount, ch.CashAmount, ch.ExpensesAmount, ch.EditedBy, ch.Edited, ch.HostId, ch.DocDate)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func (ch *Shift)ClosedShift(db *sqlx.DB) error{
	now := time.Now()
	fmt.Println("yyyy-mm-dd date format : ", now.AddDate(0, 0, 0).Format("2006-01-02"))

	ch.Closed = now

	var checkCount int
	sqlCheckExist := `select count(host_id) as vCount from cash_shift where doc_date = ?`
	err := db.Get(&checkCount, sqlCheckExist, ch.DocDate)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if (checkCount != 0) {
		sql := `UPDATE cash_shift set change_amount = ?,cash_amount = ?,expenses_amount = ?, is_closed = ?,closed_by = ?, closed = ? where  host_id = ? and doc_date  = ?`
		_, err = db.Exec(sql, ch.ChangeAmount, ch.CashAmount, ch.ExpensesAmount, 1, ch.ClosedBy, ch.Closed, ch.HostId, ch.DocDate)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func (ch *Shift)ShiftDetails(db *sqlx.DB, host_id string) error{
	sql := `select host_id,doc_date,change_begin,change_amount,cash_amount,expenses_amount,is_closed,created_by,created from cash_shift where  host_id = ? `
	err := db.Get(ch, sql, host_id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}


func (ch *Shift)SearchShiftByKeyword(db *sqlx.DB, host_id string, doc_date string)(shifts []*Shift, err error){
	sql := `select change_begin,change_amount,cash_amount,expenses_amount,is_closed,created_by,created from cash_shift where  host_id = ? order by docdate desc limit 20`
	err = db.Select(&ch, sql, host_id, doc_date)
	if err != nil {
		fmt.Println(err.Error())
		return  nil, err
	}
	return shifts, nil
}