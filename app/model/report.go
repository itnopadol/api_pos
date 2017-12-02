package model

import (
	"github.com/jmoiron/sqlx"
	"fmt"
)

type Report struct {
	ReportMonth string `json:"report_month" db:"report_month"`
	ReportYear string `json:"report_year" db:"report_year"`
	CompanyName string `json:"company_name" db:"company_name"`
	EntrePreneurName string `json:"entre_preneur_name" db:"entre_preneur_name"`
	Address string `json:"address" db:"address"`
	TaxId string `json:"tax_id" db:"tax_id"`
	TaxRate int `json:"tax_rate" db:"tax_rate"`
	Details []*Detail `json:"details"`
}

type Detail struct {
	DocDate string `json:"doc_date" db:"doc_date"`
	DocNo string `json:"doc_no" db:"doc_no"`
	CustomerName string `json:"customer_name" db:"customer_name"`
	CustTaxId string `json:"cust_tax_id" db:"cust_tax_id"`
	SumOfItemAmount float64 `json:"sum_of_item_amount" db:"sum_of_item_amount"`
	BeforeTaxAmount float64 `json:"before_tax_amount" db:"before_tax_amount"`
	TaxAmount float64 `json:"tax_amount" db:"tax_amount"`
	SumTotalAmount float64 `json:"sum_total_amount" db:"sum_total_amount"`
}


func (r *Report)ReportTax(db *sqlx.DB, report_month string, report_year string) error {
	sql := `select ? as report_month,? as report_year, ifnull(company_name,'') as company_name,ifnull(company_name,'') as entre_preneur_name,ifnull(address,'') as address,ifnull(tax_id,'') as tax_id,ifnull(tax_rate,7) as tax_rate from config `
	fmt.Println("report_month = ",report_month)
	db.Get(r, sql, report_month, report_year)

	sqlsub :=`select doc_date,doc_no,'เงินสด' as customer_name,"" as cust_tax_id,item_amount as sum_of_item_amount,before_tax_amount,tax_amount,total_amount as sum_total_amount from sale where month(doc_date) = ? and year(doc_date) = ? order by doc_date, que_id`
	err := db.Select(&r.Details, sqlsub, report_month, report_year)
	if err != nil {
		return err
	}
	return nil
}