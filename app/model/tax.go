package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
	"strconv"
)

type TaxData struct {
	//Id          int     `json:"id" db:"id"`
	//MonthTax    int     `json:"report_month" db:"month"`
	//YearTax     int     `json:"report_year" db:"year"`
	//CompanyName string  `json:"company_name" db:""`
	//Address     string  `json:"address" db:"Address"`
	//TaxId       string  `json:"tax_id" db:"TaxId"`
	//MonthSend   float64 `json:"month_send" db:"month_send"`
	//CreateBy    string  `json:"create_by" db:"create_by"`
	//ListDoc     []*Data `json:"details"`
	MonthTax         string  `json:"report_month" db:"month"`
	YearTax          string  `json:"report_year" db:"year"`
	CompanyName      string  `json:"company_name" db:""`
	EntrePreneurName string  `json:"entre_preneur_name" `
	Address          string  `json:"address" db:"Address"`
	TaxId            string  `json:"tax_id" db:"TaxId"`
	MonthSend        float64 `json:"month_send" db:"month_send"`
	TaxRate          int     `json:"tax_rate" db:"tax_rate"`
	CreateBy         string  `json:"create_by" db:"create_by"`
	ListDoc          []*Data `json:"details"`
}

type Data struct {
	DaySend         float64 `json:"day_send" db:"day_send"`
	DocDate         string  `json:"doc_date" db:"doc_date"`
	DocNo           string  `json:"doc_no" db:"doc_no"`
	CustomerName    string  `json:"customer_name" db:"customer_name"`
	CustTaxId       string  `json:"cust_tax_id" db:"cust_tax_id"`
	TaxNo           string  `json:"tax_no" db:"tax_no"`
	SumOfItemAmount float64 `json:"sum_of_item_amount" db:"sum_of_item_amount"`
	BeforeTaxAmount float64 `json:"before_tax_amount" db:"before_tax_amount"`
	TaxAmount       float64 `json:"tax_amount" db:"tax_amount"`
	TotalAmount     float64 `json:"sum_total_amount" db:"total_amount"`
}

func (tax *TaxData) GenTaxData(db *sqlx.DB, begindate string, enddate string, SendTaxAmount float64) error {
	var vDay int;
	var vSumAll float64;
	var last_number1 int
	var last_number string
	var snumber string
	var intyear int
	var vyear string

	var intmonth int
	var intmonth1 int
	var vmonth string
	var vmonth1 string
	var lenmonth int

	var intday int
	var intday1 int
	var vday string
	var vday1 string
	var lenday int

	//sql := `select count(doc_date) as day1 from (select distinct doc_date from sale where  doc_date between ? and ?) as q`
	//err := db.Get(&vDay, sql, begindate, enddate)
	//if err != nil {
	//	fmt.Println("Count Day =", err.Error())
	//	return err
	//}

	BeginDate, err := time.Parse("2006-01-02", begindate);
	fmt.Println("begindate,enddate,total", begindate, enddate, vDay, SendTaxAmount)

	fmt.Println("Day of Month = ", daysIn(BeginDate.Month(), BeginDate.Year()))

	vDay = daysIn(BeginDate.Month(), BeginDate.Year())

	fmt.Println("Count Day =", vDay)

	config := new(Config)
	config = GetConfig(db)

	//sqlsum := `select sum(total_amount) as sumtotal from sale where  doc_date between ? and ?`
	sqlsum := `select 	ifnull(sum(total_amount),0) as totalamount  from  sale  where doc_date between ? and ? and ifnull(doc_no,'') <> '' and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date between ? and ? and item_id in (select id from item where code like '%N%')) `
	err = db.Get(&vSumAll, sqlsum, begindate, enddate, begindate, enddate)
	if err != nil {
		fmt.Println("vSumAll =", err.Error())
		return err
	}

	fmt.Println("Sum All = ", vSumAll)

	//dateString := "2018-03-01"

	fmt.Println("last_number = ", last_number)

	tax.YearTax = strconv.Itoa(BeginDate.Year())
	tax.MonthTax = strconv.Itoa(int(BeginDate.Month()))
	tax.MonthSend = SendTaxAmount
	tax.CompanyName = config.CompanyName
	tax.EntrePreneurName = config.CompanyName
	tax.Address = config.Address
	tax.TaxId = config.TaxId
	tax.TaxRate = config.TaxRate

	sqldel_taxtemp := `delete from tax_temp where doc_date between ? and ?`
	fmt.Println("sqldel_taxtemp = ", sqldel_taxtemp, begindate, enddate, begindate, enddate)
	_, err = db.Exec(sqldel_taxtemp, begindate, enddate)
	if err != nil {
		fmt.Println("sqldel_taxtemp =", err.Error())
		return nil
	}

	for i := 0; i < vDay; i++ {
		var vTotalDay float64
		var vAmountPerDay float64
		var vPercentDay float64

		DateAdd := BeginDate.AddDate(0, 0, i).Format("2006-01-02")

		//sql := `select ifnull(sum(total_amount),0) as totalamount from sale where  doc_date = ?`
		sql := `select 	ifnull(sum(total_amount),0) as totalamount  from  sale  where doc_date = ? and ifnull(doc_no,'') <> '' and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date = ? and item_id in (select id from item where code like '%N%'))  order by doc_no`
		err = db.Get(&vTotalDay, sql, DateAdd, DateAdd)
		fmt.Println("DateAdd = ", DateAdd)
		if err != nil {
			fmt.Println("vTotal =", err.Error())
			return err
		}

		fmt.Println("total day=", vTotalDay)

		if vTotalDay != 0 {
			vPercentDay = (vTotalDay * 100) / vSumAll

			fmt.Println("vPercentDay =", vPercentDay)

			vAmountPerDay = (SendTaxAmount * vPercentDay) / 100

			fmt.Println("vAmountPerDay = ", vAmountPerDay)

			bill := tax.ListDoc
			sqldel := `delete from Test_Sum_Vat where  SendDayTax = ?`
			fmt.Println("sqldel = ", sqldel, DateAdd)
			_, err = db.Exec(sqldel, DateAdd)
			if err != nil {
				fmt.Println("sqldel =", err.Error())
				return nil
			}

			//sqlsub := `select doc_date,doc_no,ifnull(before_tax_amount,0) as before_tax_amount,ifnull(tax_amount,0) as tax_amount,ifnull(total_amount,'') as total_amount,'เงินสด'  as customer_name,ifnull(before_tax_amount,0) as sum_of_item_amount from sale where doc_date = ? and ifnull(doc_no,'') <> '' order by doc_no`
			sqlsub := `select 	doc_date,doc_no,ifnull(before_tax_amount,0) as before_tax_amount,ifnull(tax_amount,0) as tax_amount,ifnull(total_amount,'') as total_amount,'เงินสด'  as customer_name,ifnull(before_tax_amount,0) as sum_of_item_amount from 	sale where doc_date = ? and ifnull(doc_no,'') <> '' and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date = ? and item_id in (select id from item where code like '%N%'))  order by doc_no`
			err = db.Select(&bill, sqlsub, DateAdd, DateAdd)
			fmt.Println("sqlsub = ", sqlsub, DateAdd)
			if err != nil {
				fmt.Println("sqlsub =", err.Error())
				return nil
			}

			last_number1 = 1
			for _, d := range bill {

				var sumtotal float64

				sqlcheck := `select sum(ifnull(total_amount,0)) as sumtotal from tax_temp where doc_date = ?`
				err = db.Get(&sumtotal, sqlcheck, DateAdd)
				if err != nil {
					fmt.Println("sqlcheck", err.Error())
				}

				fmt.Println("sumtotal = ", sumtotal, " vAmountPerDay =", vAmountPerDay)

				if sumtotal < vAmountPerDay {

					last_number = strconv.Itoa(last_number1)

					DateGenDoc, err := time.Parse("2006-01-02", DateAdd);
					if (DateGenDoc.Year() >= 2560) {
						intyear = DateGenDoc.Year()
					} else {
						intyear = DateGenDoc.Year() + 543
					}

					vyear = strconv.Itoa(intyear)
					vyear1 := vyear[2:len(vyear)]

					intmonth = int(DateGenDoc.Month())
					intmonth1 = int(intmonth)
					vmonth = strconv.Itoa(intmonth1)

					lenmonth = len(vmonth)

					if (lenmonth == 1) {
						vmonth1 = "0" + vmonth
					} else {
						vmonth1 = vmonth
					}

					intday = int(DateGenDoc.Day())
					intday1 = int(intday)
					vday = strconv.Itoa(intday1)

					lenday = len(vday)

					if (lenday == 1) {
						vday1 = "0" + vday
					} else {
						vday1 = vday
					}

					if (len(string(last_number)) == 1) {
						snumber = "000" + last_number
					}
					if (len(string(last_number)) == 2) {
						snumber = "00" + last_number
					}
					if (len(string(last_number)) == 3) {
						snumber = "0" + last_number
					}
					if (len(string(last_number)) == 4) {
						snumber = last_number
					}

					new_tax_no := "01" + vyear1 + vmonth1 + vday1 + "-" + snumber //เลขที่เอกสารใหม่ส่งสรรพกร

					sqlins := `Insert into tax_temp(month_tax,year_tax,doc_date,month_send,day_send,doc_no,tax_no,before_tax_amount,tax_amount,total_amount,create_by,created) values(?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
					fmt.Println("Insert tax_temp = ", tax.MonthTax, tax.YearTax, d.DocDate, tax.MonthSend, vAmountPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount, tax.CreateBy)
					_, err = db.Exec(sqlins, tax.MonthTax, tax.YearTax, DateAdd, tax.MonthSend, vAmountPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount, tax.CreateBy)
					fmt.Println("sqlins", sqlins, d.DocDate, d.DocNo, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount)
					if err != nil {
						fmt.Println("sqlins", err.Error())
						return err
					}

					last_number1 = last_number1 + 1 //เพิ่มเลขที่บิล

				}
			}
		}

	}

	tax.YearTax = strconv.Itoa(BeginDate.Year())
	tax.MonthTax = strconv.Itoa(int(BeginDate.Month()))
	tax.MonthSend = SendTaxAmount

	sqldata := `select doc_date,day_send ,doc_no,tax_no,before_tax_amount,tax_amount,total_amount from tax_temp where doc_date between ? and ?`
	err = db.Select(&tax.ListDoc, sqldata, begindate, enddate)
	if err != nil {
		return err
	}
	return nil

}
