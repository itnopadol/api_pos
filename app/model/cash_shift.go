package model

import "time"

type Shift struct {
	Id int `json:"id" db:"id"`
	HostId int `json:"host_id" db:"host_id"`
	DocDate time.Time `json:"doc_date" db:"doc_date"`
	ChangeAmount float64 `json:"change_amount" db:"change_amount"`
	ExpensesAmount float64 `json:"expenses_amount" db:"expenses_amount"`
	CreatedBy string `json:"created_by" db:"created_by"`
	Created time.Time `json:"created" db:"created"`
	EditedBy string `json:"edited_by" db:"edited_by"`
	Edited time.Time `json:"edited" db:"edited"`

}