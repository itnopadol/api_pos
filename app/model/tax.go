package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
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
	NoVat           float64 `json:"no_vat" db:"no_vat"`
	TotalAmount     float64 `json:"sum_total_amount" db:"total_amount"`
}

func (tax *TaxData) GenTaxWithNoVatData(db *sqlx.DB, begindate string, enddate string, SendNoVatAmount float64, SendTotalAmount float64) error {
	var vDay int
	var vSumAll float64
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

	BeginDate, err := time.Parse("2006-01-02", begindate)
	fmt.Println("begindate,enddate,total", begindate, enddate, vDay, SendNoVatAmount)

	fmt.Println("Day of Month = ", daysIn(BeginDate.Month(), BeginDate.Year()))

	vDay = daysIn(BeginDate.Month(), BeginDate.Year())

	fmt.Println("Count Day =", vDay)

	config := new(Config)
	config = GetConfig(db)

	sqlsum := `select ifnull(sum(no_vat),0)  as novat from (select doc_date,doc_no,ifnull(sum(before_tax_amount),0) as before_tax_amount,ifnull(sum(tax_amount),0) as  tax_amount,ifnull(sum(no_vat_amount),0) as no_vat,ifnull(sum(amount),0) as total_amount,'เงินสด'  as customer_name,ifnull(sum(before_tax_amount),0) as sum_of_item_amount
			from
			(
				select a.doc_no,a.doc_date,ifnull(b.amount,0) as amount,c.code,b.item_name,
				case when c.code like '%n' then 0
				else round((b.amount*100)/107,2) end as  before_tax_amount ,
				case when c.code like '%n' then 0
				else round(b.amount- ((b.amount*100)/107),2) end as  tax_amount,
				case when c.code like '%n' then b.amount
				else 0 end as  no_vat_amount
				from 	sale a
				inner join sale_sub b on a.id = b.sale_id
				inner join item c on b.item_id = c.id
				where a.doc_date between ? and ? and ifnull(a.doc_no,'') <> ''
			) as rs
			group by doc_no,doc_date) as result
			where no_vat > 0
			order by doc_date`
	fmt.Println("sql sum =", sqlsum, begindate, enddate)
	err = db.Get(&vSumAll, sqlsum, begindate, enddate) //, begindate, enddate)
	if err != nil {
		fmt.Println("vSumAll =", err.Error())
		return err
	}

	fmt.Println("Sum All = ", vSumAll)

	//dateString := "2018-03-01"

	fmt.Println("last_number = ", last_number)

	tax.YearTax = strconv.Itoa(BeginDate.Year())
	tax.MonthTax = strconv.Itoa(int(BeginDate.Month()))
	tax.MonthSend = SendNoVatAmount
	tax.CompanyName = config.CompanyName
	tax.EntrePreneurName = config.CompanyName
	tax.Address = config.Address
	tax.TaxId = config.TaxId
	tax.TaxRate = config.TaxRate

	sqldel_taxtemp := `delete from tax_temp_all where doc_date between ? and ?`
	fmt.Println("sqldel_taxtemp_all = ", sqldel_taxtemp, begindate, enddate, begindate, enddate)
	_, err = db.Exec(sqldel_taxtemp, begindate, enddate)
	if err != nil {
		fmt.Println("sqldel_taxtemp_all =", err.Error())
		return nil
	}

	for i := 0; i < vDay; i++ {
		var vTotalDay float64
		var vAmountPerDay float64
		var vPercentDay float64

		DateAdd := BeginDate.AddDate(0, 0, i).Format("2006-01-02")

		//sql := `select ifnull(sum(total_amount),0) as totalamount from sale where  doc_date = ?`
		//sql := `select 	ifnull(sum(total_amount),0) as totalamount  from  sale  where doc_date = ? and ifnull(doc_no,'') <> '' order by doc_no` //and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date = ? and item_id in (select id from item where code like '%N%'))  order by doc_no`

		sql := `select ifnull(sum(no_vat),0) as novat from (select doc_date,doc_no,ifnull(sum(before_tax_amount),0) as before_tax_amount,ifnull(sum(tax_amount),0) as  tax_amount,ifnull(sum(no_vat_amount),0) as no_vat,ifnull(sum(amount),0) as total_amount,'เงินสด'  as customer_name,ifnull(sum(before_tax_amount),0) as sum_of_item_amount
			from
			(
				select a.doc_no,a.doc_date,ifnull(b.amount,0) as amount,c.code,b.item_name,
				case when c.code like '%n' then 0
				else round((b.amount*100)/107,2) end as  before_tax_amount ,
				case when c.code like '%n' then 0
				else round(b.amount- ((b.amount*100)/107),2) end as  tax_amount,
				case when c.code like '%n' then b.amount
				else 0 end as  no_vat_amount
				from 	sale a
				inner join sale_sub b on a.id = b.sale_id
				inner join item c on b.item_id = c.id
				where a.doc_date = ? and ifnull(a.doc_no,'') <> ''
			) as rs
			group by doc_no,doc_date) as result 
			where no_vat > 0
			order by doc_date`
		err = db.Get(&vTotalDay, sql, DateAdd) //, DateAdd)
		fmt.Println("DateAdd = ", DateAdd)
		if err != nil {
			fmt.Println("vTotal =", err.Error())
			return err
		}

		fmt.Println("total day=", vTotalDay)

		if vTotalDay != 0 {
			vPercentDay = (vTotalDay * 100) / vSumAll

			fmt.Println("vPercentDay =", vPercentDay, vAmountPerDay, vTotalDay)

			vAmountPerDay = (SendNoVatAmount * vPercentDay) / 100

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
			//sqlsub := `select 	doc_date,doc_no,ifnull(before_tax_amount,0) as before_tax_amount,ifnull(tax_amount,0) as tax_amount,ifnull(total_amount,'') as total_amount,'เงินสด'  as customer_name,ifnull(before_tax_amount,0) as sum_of_item_amount from 	sale where doc_date = ? and ifnull(doc_no,'') <> '' order by doc_no ` // and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date = ? and item_id in (select id from item where code like '%N%'))  order by doc_no`
			sqlsub := `select * from (select doc_date,doc_no,ifnull(sum(before_tax_amount),0) as before_tax_amount,ifnull(sum(tax_amount),0) as  tax_amount,ifnull(sum(no_vat_amount),0) as no_vat,ifnull(sum(amount),0) as total_amount,'เงินสด'  as customer_name,ifnull(sum(before_tax_amount),0) as sum_of_item_amount
			from
			(
				select a.doc_no,a.doc_date,b.amount,c.code,b.item_name,
				case when c.code like '%n' then 0
				else round((b.amount*100)/107,2) end as  before_tax_amount ,
				case when c.code like '%n' then 0
				else round(b.amount- ((b.amount*100)/107),2) end as  tax_amount,
				case when c.code like '%n' then b.amount
				else 0 end as  no_vat_amount
				from 	sale a
				inner join sale_sub b on a.id = b.sale_id
				inner join item c on b.item_id = c.id
				where a.doc_date = ?
			) as rs
			group by doc_no,doc_date) as result
			where no_vat > 0
			order by doc_date`
			err = db.Select(&bill, sqlsub, DateAdd) //, DateAdd)
			fmt.Println("sqlsub = ", sqlsub, DateAdd)
			if err != nil {
				fmt.Println("sqlsub =", err.Error())
				return nil
			}

			last_number1 = 1
			for _, d := range bill {

				var sumtotal float64

				//sqlcheck := `select sum(ifnull(total_amount,0)) as sumtotal from tax_temp where doc_date = ?`
				sqlcheck := `select sum(ifnull(no_vat,0)) as sumtotal from tax_temp_all where doc_date = ? and no_vat > 0`
				err = db.Get(&sumtotal, sqlcheck, DateAdd)
				if err != nil {
					fmt.Println("sqlcheck", err.Error())
				}

				fmt.Println("sumtotal = ", sumtotal, " vAmountPerDay =", vAmountPerDay)

				if sumtotal < vAmountPerDay {

					last_number = strconv.Itoa(last_number1)

					DateGenDoc, err := time.Parse("2006-01-02", DateAdd)
					if DateGenDoc.Year() >= 2560 {
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

					if lenmonth == 1 {
						vmonth1 = "0" + vmonth
					} else {
						vmonth1 = vmonth
					}

					intday = int(DateGenDoc.Day())
					intday1 = int(intday)
					vday = strconv.Itoa(intday1)

					lenday = len(vday)

					if lenday == 1 {
						vday1 = "0" + vday
					} else {
						vday1 = vday
					}

					if len(string(last_number)) == 1 {
						snumber = "000" + last_number
					}
					if len(string(last_number)) == 2 {
						snumber = "00" + last_number
					}
					if len(string(last_number)) == 3 {
						snumber = "0" + last_number
					}
					if len(string(last_number)) == 4 {
						snumber = last_number
					}

					new_tax_no := "01" + vyear1 + vmonth1 + vday1 + "-" + snumber //เลขที่เอกสารใหม่ส่งสรรพกร

					fmt.Println("day send = ", vAmountPerDay)

					sqlins := `Insert into tax_temp_all(month_tax,year_tax,doc_date,month_send,month_no_vat_send,day_send,doc_no,tax_no,before_tax_amount,tax_amount,no_vat,total_amount,create_by,created) values(?,?,?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
					//fmt.Println("Insert tax_temp Head= ", tax.MonthTax, tax.YearTax, d.DocDate, tax.MonthSend, vAmountPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.NoVat, d.TotalAmount, tax.CreateBy)
					_, err = db.Exec(sqlins, tax.MonthTax, tax.YearTax, DateAdd, SendTotalAmount, tax.MonthSend, vAmountPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.NoVat, d.TotalAmount, tax.CreateBy)
					fmt.Println("sqlins head = ", sqlins, d.DocDate, d.DocNo, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount)
					if err != nil {
						fmt.Println("sqlins", err.Error())
						return err
					}

					last_number1 = last_number1 + 1 //เพิ่มเลขที่บิล

				}
			}
		}

	}

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	var sum_total_tax float64

	var vSumAllTotal float64
	//var last_numberVat1 int
	//var last_number string
	//var snumberVat string
	//var intyearVat int
	//var vyearVat string
	//
	//var intmonthVat int
	//var intmonthVat1 int
	//var vmonthVat string
	//var vmonthVat1 string
	//var lenmonthVat int
	//
	//var intdayVat int
	//var intdayVat1 int
	//var vdayVat string
	//var vdayVat1 string
	//var lendayVat int

	sqlsumvat := `select ifnull(sum(total_amount),0)  as total_amount from (select doc_date,doc_no,ifnull(sum(before_tax_amount),0) as before_tax_amount,ifnull(sum(tax_amount),0) as  tax_amount,ifnull(sum(no_vat_amount),0) as no_vat,ifnull(sum(amount),0) as total_amount,'เงินสด'  as customer_name,ifnull(sum(before_tax_amount),0) as sum_of_item_amount
			from
			(
				select a.doc_no,a.doc_date,b.amount,c.code,b.item_name,
				case when c.code like '%n' then 0
				else round((b.amount*100)/107,2) end as  before_tax_amount ,
				case when c.code like '%n' then 0
				else round(b.amount- ((b.amount*100)/107),2) end as  tax_amount,
				case when c.code like '%n' then b.amount
				else 0 end as  no_vat_amount
				from 	sale a
				inner join sale_sub b on a.id = b.sale_id
				inner join item c on b.item_id = c.id
				where a.doc_date between ? and ? and ifnull(a.doc_no,'') <> ''
			) as rs
			group by doc_no,doc_date) as result  
			where no_vat = 0  
			order by doc_date`
	err = db.Get(&vSumAllTotal, sqlsumvat, begindate, enddate) //, begindate, enddate)
	if err != nil {
		fmt.Println("vSumAll =", err.Error())
		return err
	}

	sql_total := `select ifnull(sum(total_amount),0) as sum_total_tax from tax_temp_all where year_tax = ? and month_tax = ?`
	err = db.Get(&sum_total_tax, sql_total, tax.YearTax, tax.MonthTax)
	if err != nil {
		fmt.Println("sqlcheck", err.Error())
	}

	fmt.Println("sum_total_tax", sum_total_tax)

	TaxTotalRemain := SendTotalAmount - sum_total_tax

	fmt.Println("TaxTotalRemain,SendTotalAmount", TaxTotalRemain, SendTotalAmount)

	//if sum_total_tax < SendTotalAmount {
	for i := 0; i < vDay; i++ {
		var vTotalVatDay float64
		var vAmountVatPerDay float64
		var vPercentVatDay float64

		DateAdd := BeginDate.AddDate(0, 0, i).Format("2006-01-02")

		sql := `select ifnull(sum(total_amount),0) as total_amount from (select doc_date,doc_no,ifnull(sum(before_tax_amount),0) as before_tax_amount,ifnull(sum(tax_amount),0) as  tax_amount,ifnull(sum(no_vat_amount),0) as no_vat,ifnull(sum(amount),0) as total_amount,'เงินสด'  as customer_name,ifnull(sum(before_tax_amount),0) as sum_of_item_amount
			from
			(
				select a.doc_no,a.doc_date,ifnull(b.amount,0) as amount,c.code,b.item_name,
				case when c.code like '%n' then 0
				else round((ifnull(b.amount,0)*100)/107,2) end as  before_tax_amount ,
				case when c.code like '%n' then 0
				else round(ifnull(b.amount,0)- ((ifnull(b.amount,0)*100)/107),2) end as  tax_amount,
				case when c.code like '%n' then ifnull(b.amount,0)
				else 0 end as  no_vat_amount
				from 	sale a
				inner join sale_sub b on a.id = b.sale_id
				inner join item c on b.item_id = c.id
				where a.doc_date = ?
			) as rs
			group by doc_no,doc_date) as result
			where no_vat = 0  
			order by doc_date`
		err = db.Get(&vTotalVatDay, sql, DateAdd)
		fmt.Println("DateAdd = ", DateAdd)
		if err != nil {
			fmt.Println("vTotal =", err.Error())
			return err
		}

		fmt.Println("total day=", vTotalVatDay, vSumAllTotal)

		if vTotalVatDay != 0 {
			vPercentVatDay = (vTotalVatDay * 100) / vSumAllTotal

			fmt.Println("vPercentDay =", vPercentVatDay, vTotalVatDay)

			vAmountVatPerDay = (TaxTotalRemain * vPercentVatDay) / 100
			//vAmountVatPerDay = (SendTotalAmount * vPercentVatDay) / 100

			fmt.Println("vAmountPerDay = ", vAmountVatPerDay)

			bill := tax.ListDoc
			//sqldel := `delete from Test_Sum_Vat where  SendDayTax = ?`
			//fmt.Println("sqldel = ", sqldel, DateAdd)
			//_, err = db.Exec(sqldel, DateAdd)
			//if err != nil {
			//	fmt.Println("sqldel =", err.Error())
			//	return nil
			//}

			//sqlsub := `select doc_date,doc_no,ifnull(before_tax_amount,0) as before_tax_amount,ifnull(tax_amount,0) as tax_amount,ifnull(total_amount,'') as total_amount,'เงินสด'  as customer_name,ifnull(before_tax_amount,0) as sum_of_item_amount from sale where doc_date = ? and ifnull(doc_no,'') <> '' order by doc_no`
			//sqlsub := `select 	doc_date,doc_no,ifnull(before_tax_amount,0) as before_tax_amount,ifnull(tax_amount,0) as tax_amount,ifnull(total_amount,'') as total_amount,'เงินสด'  as customer_name,ifnull(before_tax_amount,0) as sum_of_item_amount from 	sale where doc_date = ? and ifnull(doc_no,'') <> '' order by doc_no ` // and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date = ? and item_id in (select id from item where code like '%N%'))  order by doc_no`
			sqlsub := `select * from (select doc_date,doc_no,ifnull(sum(before_tax_amount),0) as before_tax_amount,ifnull(sum(tax_amount),0) as  tax_amount,ifnull(sum(no_vat_amount),0) as no_vat,ifnull(sum(amount),0) as total_amount,'เงินสด'  as customer_name,ifnull(sum(before_tax_amount),0) as sum_of_item_amount
			from
			(
				select a.doc_no,a.doc_date,ifnull(b.amount,0) as amount,c.code,b.item_name,
				case when c.code like '%n' then 0
				else round((b.amount*100)/107,2) end as  before_tax_amount ,
				case when c.code like '%n' then 0
				else round(b.amount- ((b.amount*100)/107),2) end as  tax_amount,
				case when c.code like '%n' then b.amount
				else 0 end as  no_vat_amount
				from 	sale a
				inner join sale_sub b on a.id = b.sale_id
				inner join item c on b.item_id = c.id
				where a.doc_date = ?  
			) as rs
			group by doc_no,doc_date) as result where no_vat = 0  order by doc_date`
			err = db.Select(&bill, sqlsub, DateAdd) //, DateAdd)
			fmt.Println("sqlsub = ", sqlsub, DateAdd)
			if err != nil {
				fmt.Println("sqlsub =", err.Error())
				return nil
			}

			//last_number1 = 1

			sql_last_no := `select convert(right(max(tax_no),4),UNSIGNED INTEGER)+1 as last_number1  from tax_temp_all where doc_date = ? `
			err = db.Get(&last_number1, sql_last_no, DateAdd)
			if err != nil {
				fmt.Println("sql_last_no", err.Error())
			}
			fmt.Println("sql_last_no = ", sql_last_no, DateAdd, last_number1)

			for _, d := range bill {

				var sumtotal float64

				sqlcheck := `select ifnull(sum(ifnull(total_amount,0)),0) as sumtotal from tax_temp_all where doc_date = ? and no_vat = 0`
				err = db.Get(&sumtotal, sqlcheck, DateAdd)
				if err != nil {
					fmt.Println("sqlcheck", err.Error())
				}

				fmt.Println("sumtotal = ", sumtotal, " vAmountPerDay =", vAmountVatPerDay, "last number = ", last_number1)

				if sumtotal < vAmountVatPerDay {

					last_number = strconv.Itoa(last_number1)

					fmt.Println("last_number = ", last_number, "     ", strconv.Itoa(last_number1))

					DateGenDoc, err := time.Parse("2006-01-02", DateAdd)
					if DateGenDoc.Year() >= 2560 {
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

					if lenmonth == 1 {
						vmonth1 = "0" + vmonth
					} else {
						vmonth1 = vmonth
					}

					intday = int(DateGenDoc.Day())
					intday1 = int(intday)
					vday = strconv.Itoa(intday1)

					lenday = len(vday)

					if lenday == 1 {
						vday1 = "0" + vday
					} else {
						vday1 = vday
					}

					fmt.Println("len(string(last_number)) = ", len(string(last_number)))

					if len(string(last_number)) == 1 {
						snumber = "000" + last_number
					}
					if len(string(last_number)) == 2 {
						snumber = "00" + last_number
					}
					if len(string(last_number)) == 3 {
						snumber = "0" + last_number
					}
					if len(string(last_number)) == 4 {
						snumber = last_number
					}

					new_tax_no_sub := "01" + vyear1 + vmonth1 + vday1 + "-" + snumber //เลขที่เอกสารใหม่ส่งสรรพกร

					fmt.Println("day send new_tax_no_sub = ", vAmountVatPerDay, new_tax_no_sub)

					sqlins := `Insert into tax_temp_all(month_tax,year_tax,doc_date,month_send,month_no_vat_send,day_send,doc_no,tax_no,before_tax_amount,tax_amount,no_vat,total_amount,create_by,created) values(?,?,?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
					//fmt.Println("Insert tax_temp Sub = ", tax.MonthTax, tax.YearTax, d.DocDate, SendTotalAmount, tax.MonthSend, vAmountVatPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.NoVat, d.TotalAmount, tax.CreateBy)
					_, err = db.Exec(sqlins, tax.MonthTax, tax.YearTax, DateAdd, SendTotalAmount, tax.MonthSend, vAmountVatPerDay, d.DocNo, new_tax_no_sub, d.BeforeTaxAmount, d.TaxAmount, d.NoVat, d.TotalAmount, tax.CreateBy)
					fmt.Println("sqlins sub =", sqlins, d.DocDate, d.DocNo, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount)
					if err != nil {
						fmt.Println("sqlins", err.Error())
						return err
					}

					last_number1 = last_number1 + 1 //เพิ่มเลขที่บิล

				}
			}
		}

	}
	//}

	//	sql := `select count(*) as count_day
	//from
	//(
	//select distinct doc_date from (select doc_date,doc_no,ifnull(sum(before_tax_amount),0) as before_tax_amount,ifnull(sum(tax_amount),0) as  tax_amount,ifnull(sum(no_vat_amount),0) as no_vat,ifnull(sum(amount),0) as total_amount,'เงินสด'  as customer_name,ifnull(sum(before_tax_amount),0) as sum_of_item_amount
	//			from
	//			(
	//				select a.doc_no,a.doc_date,b.amount,c.code,b.item_name,
	//				case when c.code like '%n' then 0
	//				else round((b.amount*100)/107,2) end as  before_tax_amount ,
	//				case when c.code like '%n' then 0
	//				else round(b.amount- ((b.amount*100)/107),2) end as  tax_amount,
	//				case when c.code like '%n' then b.amount
	//				else 0 end as  no_vat_amount
	//				from 	sale a
	//				inner join sale_sub b on a.id = b.sale_id
	//				inner join item c on b.item_id = c.id
	//				where a.doc_date between '2019/09/01' and '2019/09/30'
	//			) as rs
	//			group by doc_no,doc_date) as result
	//			where no_vat = 0
	//			group  by doc_date
	//) as count_day`
	//	err = db.Get(&count_day, sql)
	//	if err != nil {
	//		fmt.Println("sqlcheck", err.Error())
	//	}

	tax.YearTax = strconv.Itoa(BeginDate.Year())
	tax.MonthTax = strconv.Itoa(int(BeginDate.Month()))
	tax.MonthSend = 0 //SendNoVat

	sqldata := `select doc_date,day_send ,tax_no as doc_no,doc_no as tax_no,ifnull(before_tax_amount,0) as sum_of_item_amount,ifnull(tax_amount,0) as tax_amount,ifnull(no_vat,0) as no_vat,ifnull(total_amount,0) as total_amount from tax_temp_all where doc_date between ? and ? order by doc_date, doc_no`
	err = db.Select(&tax.ListDoc, sqldata, begindate, enddate)
	if err != nil {
		return err
	}
	return nil
}

func (tax *TaxData) GenTaxData(db *sqlx.DB, company_id int64, branch_id int64, begindate string, enddate string, SendTotalAmount float64) error {
	var vDay int
	var vSumAll float64
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

	BeginDate, err := time.Parse("2006-01-02", begindate)
	fmt.Println("begindate,enddate,total", begindate, enddate, vDay, SendTotalAmount)

	fmt.Println("Day of Month = ", daysIn(BeginDate.Month(), BeginDate.Year()))

	vDay = daysIn(BeginDate.Month(), BeginDate.Year())

	fmt.Println("Count Day =", vDay)

	config := new(Config)
	config = GetConfig(db)

	sqlsum := `select 	ifnull(sum(total_amount),0) as totalamount  from  sale  where company_id = ? and branch_id = ? and doc_date between ? and ? and ifnull(doc_no,'') <> ''` //and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date between ? and ? and item_id in (select id from item where code like '%N%')) `
	fmt.Println("sql sum =", sqlsum, begindate, enddate)
	err = db.Get(&vSumAll, sqlsum, company_id, branch_id, begindate, enddate) //, begindate, enddate)
	if err != nil {
		fmt.Println("vSumAll =", err.Error())
		return err
	}

	fmt.Println("Sum All = ", vSumAll)

	//dateString := "2018-03-01"

	fmt.Println("last_number = ", last_number)

	tax.YearTax = strconv.Itoa(BeginDate.Year())
	tax.MonthTax = strconv.Itoa(int(BeginDate.Month()))
	tax.MonthSend = SendTotalAmount
	tax.CompanyName = config.CompanyName
	tax.EntrePreneurName = config.CompanyName
	tax.Address = config.Address
	tax.TaxId = config.TaxId
	tax.TaxRate = config.TaxRate

	sqldel_taxtemp := `delete from tax_temp where company_id = ? and branch_id = ? and doc_date between ? and ?`
	fmt.Println("sqldel_taxtemp = ", sqldel_taxtemp, company_id, branch_id, begindate, enddate, begindate, enddate)
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
		sql := `select 	ifnull(sum(total_amount),0) as totalamount  from  sale  where company_id = ? and branch_id = ? and doc_date = ? and ifnull(doc_no,'') <> '' order by doc_no` //and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date = ? and item_id in (select id from item where code like '%N%'))  order by doc_no`
		err = db.Get(&vTotalDay, sql, company_id, branch_id, DateAdd)                                                                                                                //, DateAdd)
		fmt.Println("DateAdd = ", DateAdd)
		if err != nil {
			fmt.Println("vTotal =", err.Error())
			return err
		}

		fmt.Println("total day=", vTotalDay)

		if vTotalDay != 0 {
			vPercentDay = (vTotalDay * 100) / vSumAll

			fmt.Println("vPercentDay =", vPercentDay, vAmountPerDay, vTotalDay)

			vAmountPerDay = (SendTotalAmount * vPercentDay) / 100

			fmt.Println("vAmountPerDay = ", vAmountPerDay)

			bill := tax.ListDoc
			sqldel := `delete from Test_Sum_Vat where compny_id = ? and branch_id = ? and  SendDayTax = ?`
			fmt.Println("sqldel = ", sqldel, DateAdd)
			_, err = db.Exec(sqldel, company_id, branch_id, DateAdd)
			if err != nil {
				fmt.Println("sqldel =", err.Error())
				return nil
			}

			//sqlsub := `select doc_date,doc_no,ifnull(before_tax_amount,0) as before_tax_amount,ifnull(tax_amount,0) as tax_amount,ifnull(total_amount,'') as total_amount,'เงินสด'  as customer_name,ifnull(before_tax_amount,0) as sum_of_item_amount from sale where doc_date = ? and ifnull(doc_no,'') <> '' order by doc_no`
			//sqlsub := `select 	doc_date,doc_no,ifnull(before_tax_amount,0) as before_tax_amount,ifnull(tax_amount,0) as tax_amount,ifnull(total_amount,'') as total_amount,'เงินสด'  as customer_name,ifnull(before_tax_amount,0) as sum_of_item_amount from 	sale where doc_date = ? and ifnull(doc_no,'') <> '' order by doc_no ` // and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date = ? and item_id in (select id from item where code like '%N%'))  order by doc_no`
			sqlsub := `select * from (select doc_date,doc_no,ifnull(sum(before_tax_amount),0) as before_tax_amount,ifnull(sum(tax_amount),0) as  tax_amount,ifnull(sum(no_vat_amount),0) as no_vat,ifnull(sum(amount),0) as total_amount,'เงินสด'  as customer_name,ifnull(sum(before_tax_amount),0) as sum_of_item_amount
			from
			(
				select a.doc_no,a.doc_date,b.amount,c.code,b.item_name,
				case when c.code like '%n' then 0 
				else round((b.amount*100)/107,2) end as  before_tax_amount ,
				case when c.code like '%n' then 0 
				else round(b.amount- ((b.amount*100)/107),2) end as  tax_amount,
				case when c.code like '%n' then b.amount 
				else 0 end as  no_vat_amount  
				from 	sale a 
				inner join sale_sub b on a.id = b.sale_id 
				inner join item c on b.item_id = c.id
				where a.company_id=? a.branch_id = ? and a.doc_date = ?
			) as rs 
			group by doc_no,doc_date) as result order by doc_date`
			err = db.Select(&bill, sqlsub, company_id, branch_id, DateAdd) //, DateAdd)
			fmt.Println("sqlsub = ", sqlsub, DateAdd)
			if err != nil {
				fmt.Println("sqlsub =", err.Error())
				return nil
			}

			last_number1 = 1
			for _, d := range bill {

				var sumtotal float64

				sqlcheck := `select sum(ifnull(total_amount,0)) as sumtotal from tax_temp where company_id = ? and branch_id = ? and doc_date = ?`
				err = db.Get(&sumtotal, sqlcheck, company_id, branch_id, DateAdd)
				if err != nil {
					fmt.Println("sqlcheck", err.Error())
				}

				fmt.Println("sumtotal = ", sumtotal, " vAmountPerDay =", vAmountPerDay)

				if sumtotal < vAmountPerDay {

					last_number = strconv.Itoa(last_number1)

					DateGenDoc, err := time.Parse("2006-01-02", DateAdd)
					if DateGenDoc.Year() >= 2560 {
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

					if lenmonth == 1 {
						vmonth1 = "0" + vmonth
					} else {
						vmonth1 = vmonth
					}

					intday = int(DateGenDoc.Day())
					intday1 = int(intday)
					vday = strconv.Itoa(intday1)

					lenday = len(vday)

					if lenday == 1 {
						vday1 = "0" + vday
					} else {
						vday1 = vday
					}

					if len(string(last_number)) == 1 {
						snumber = "000" + last_number
					}
					if len(string(last_number)) == 2 {
						snumber = "00" + last_number
					}
					if len(string(last_number)) == 3 {
						snumber = "0" + last_number
					}
					if len(string(last_number)) == 4 {
						snumber = last_number
					}

					new_tax_no := "01" + vyear1 + vmonth1 + vday1 + "-" + snumber //เลขที่เอกสารใหม่ส่งสรรพกร

					fmt.Println("day send = ", vAmountPerDay)

					sqlins := `Insert into tax_temp(company_id,branch_id,month_tax,year_tax,doc_date,month_send,day_send,doc_no,tax_no,before_tax_amount,tax_amount,no_vat,total_amount,create_by,created) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
					fmt.Println("Insert tax_temp = ", tax.MonthTax, tax.YearTax, d.DocDate, tax.MonthSend, vAmountPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.NoVat, d.TotalAmount, tax.CreateBy)
					_, err = db.Exec(sqlins, company_id, branch_id, tax.MonthTax, tax.YearTax, DateAdd, tax.MonthSend, vAmountPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.NoVat, d.TotalAmount, tax.CreateBy)
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
	tax.MonthSend = SendTotalAmount

	sqldata := `select doc_date,day_send ,tax_no as doc_no,doc_no as tax_no,before_tax_amount as sum_of_item_amount,tax_amount,no_vat,total_amount from tax_temp where company_id = ? and branch_id = ? and doc_date between ? and ?`
	err = db.Select(&tax.ListDoc, sqldata, company_id,branch_id, begindate, enddate)
	if err != nil {
		return err
	}
	return nil

	//var vDay int;
	//var vSumAll float64;
	//var last_number1 int
	//var last_number string
	//var snumber string
	//var intyear int
	//var vyear string
	//
	//var intmonth int
	//var intmonth1 int
	//var vmonth string
	//var vmonth1 string
	//var lenmonth int
	//
	//var intday int
	//var intday1 int
	//var vday string
	//var vday1 string
	//var lenday int
	//
	////sql := `select count(doc_date) as day1 from (select distinct doc_date from sale where  doc_date between ? and ?) as q`
	////err := db.Get(&vDay, sql, begindate, enddate)
	////if err != nil {
	////	fmt.Println("Count Day =", err.Error())
	////	return err
	////}
	//
	//BeginDate, err := time.Parse("2006-01-02", begindate);
	//fmt.Println("begindate,enddate,total", begindate, enddate, vDay, SendTaxAmount)
	//
	//fmt.Println("Day of Month = ", daysIn(BeginDate.Month(), BeginDate.Year()))
	//
	//vDay = daysIn(BeginDate.Month(), BeginDate.Year())
	//
	//fmt.Println("Count Day =", vDay)
	//
	//config := new(Config)
	//config = GetConfig(db)
	//
	////sqlsum := `select sum(total_amount) as sumtotal from sale where  doc_date between ? and ?`
	//sqlsum := `select 	ifnull(sum(total_amount),0) as totalamount  from  sale  where doc_date between ? and ? and ifnull(doc_no,'') <> '' and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date between ? and ? and item_id in (select id from item where code like '%N%')) `
	//err = db.Get(&vSumAll, sqlsum, begindate, enddate, begindate, enddate)
	//if err != nil {
	//	fmt.Println("vSumAll =", err.Error())
	//	return err
	//}
	//
	//fmt.Println("Sum All = ", vSumAll)
	//
	////dateString := "2018-03-01"
	//
	//fmt.Println("last_number = ", last_number)
	//
	//tax.YearTax = strconv.Itoa(BeginDate.Year())
	//tax.MonthTax = strconv.Itoa(int(BeginDate.Month()))
	//tax.MonthSend = SendTaxAmount
	//tax.CompanyName = config.CompanyName
	//tax.EntrePreneurName = config.CompanyName
	//tax.Address = config.Address
	//tax.TaxId = config.TaxId
	//tax.TaxRate = config.TaxRate
	//
	//sqldel_taxtemp := `delete from tax_temp where doc_date between ? and ?`
	//fmt.Println("sqldel_taxtemp = ", sqldel_taxtemp, begindate, enddate, begindate, enddate)
	//_, err = db.Exec(sqldel_taxtemp, begindate, enddate)
	//if err != nil {
	//	fmt.Println("sqldel_taxtemp =", err.Error())
	//	return nil
	//}
	//
	//for i := 0; i < vDay; i++ {
	//	var vTotalDay float64
	//	var vAmountPerDay float64
	//	var vPercentDay float64
	//
	//	DateAdd := BeginDate.AddDate(0, 0, i).Format("2006-01-02")
	//
	//	//sql := `select ifnull(sum(total_amount),0) as totalamount from sale where  doc_date = ?`
	//	sql := `select 	ifnull(sum(total_amount),0) as totalamount  from  sale  where doc_date = ? and ifnull(doc_no,'') <> '' and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date = ? and item_id in (select id from item where code like '%N%'))  order by doc_no`
	//	err = db.Get(&vTotalDay, sql, DateAdd, DateAdd)
	//	fmt.Println("DateAdd = ", DateAdd)
	//	if err != nil {
	//		fmt.Println("vTotal =", err.Error())
	//		return err
	//	}
	//
	//	fmt.Println("total day=", vTotalDay)
	//
	//	if vTotalDay != 0 {
	//		vPercentDay = (vTotalDay * 100) / vSumAll
	//
	//		fmt.Println("vPercentDay =", vPercentDay)
	//
	//		vAmountPerDay = (SendTaxAmount * vPercentDay) / 100
	//
	//		fmt.Println("vAmountPerDay = ", vAmountPerDay)
	//
	//		bill := tax.ListDoc
	//		sqldel := `delete from Test_Sum_Vat where  SendDayTax = ?`
	//		fmt.Println("sqldel = ", sqldel, DateAdd)
	//		_, err = db.Exec(sqldel, DateAdd)
	//		if err != nil {
	//			fmt.Println("sqldel =", err.Error())
	//			return nil
	//		}
	//
	//		//sqlsub := `select doc_date,doc_no,ifnull(before_tax_amount,0) as before_tax_amount,ifnull(tax_amount,0) as tax_amount,ifnull(total_amount,'') as total_amount,'เงินสด'  as customer_name,ifnull(before_tax_amount,0) as sum_of_item_amount from sale where doc_date = ? and ifnull(doc_no,'') <> '' order by doc_no`
	//		sqlsub := `select 	doc_date,doc_no,ifnull(before_tax_amount,0) as before_tax_amount,ifnull(tax_amount,0) as tax_amount,ifnull(total_amount,'') as total_amount,'เงินสด'  as customer_name,ifnull(before_tax_amount,0) as sum_of_item_amount from 	sale where doc_date = ? and ifnull(doc_no,'') <> '' and id not in (select a.id from sale a inner join sale_sub b on a.id = b.sale_id where doc_date = ? and item_id in (select id from item where code like '%N%'))  order by doc_no`
	//		err = db.Select(&bill, sqlsub, DateAdd, DateAdd)
	//		fmt.Println("sqlsub = ", sqlsub, DateAdd)
	//		if err != nil {
	//			fmt.Println("sqlsub =", err.Error())
	//			return nil
	//		}
	//
	//		last_number1 = 1
	//		for _, d := range bill {
	//
	//			var sumtotal float64
	//
	//			sqlcheck := `select sum(ifnull(total_amount,0)) as sumtotal from tax_temp where doc_date = ?`
	//			err = db.Get(&sumtotal, sqlcheck, DateAdd)
	//			if err != nil {
	//				fmt.Println("sqlcheck", err.Error())
	//			}
	//
	//			fmt.Println("sumtotal = ", sumtotal, " vAmountPerDay =", vAmountPerDay)
	//
	//			if sumtotal < vAmountPerDay {
	//
	//				last_number = strconv.Itoa(last_number1)
	//
	//				DateGenDoc, err := time.Parse("2006-01-02", DateAdd);
	//				if (DateGenDoc.Year() >= 2560) {
	//					intyear = DateGenDoc.Year()
	//				} else {
	//					intyear = DateGenDoc.Year() + 543
	//				}
	//
	//				vyear = strconv.Itoa(intyear)
	//				vyear1 := vyear[2:len(vyear)]
	//
	//				intmonth = int(DateGenDoc.Month())
	//				intmonth1 = int(intmonth)
	//				vmonth = strconv.Itoa(intmonth1)
	//
	//				lenmonth = len(vmonth)
	//
	//				if (lenmonth == 1) {
	//					vmonth1 = "0" + vmonth
	//				} else {
	//					vmonth1 = vmonth
	//				}
	//
	//				intday = int(DateGenDoc.Day())
	//				intday1 = int(intday)
	//				vday = strconv.Itoa(intday1)
	//
	//				lenday = len(vday)
	//
	//				if (lenday == 1) {
	//					vday1 = "0" + vday
	//				} else {
	//					vday1 = vday
	//				}
	//
	//				if (len(string(last_number)) == 1) {
	//					snumber = "000" + last_number
	//				}
	//				if (len(string(last_number)) == 2) {
	//					snumber = "00" + last_number
	//				}
	//				if (len(string(last_number)) == 3) {
	//					snumber = "0" + last_number
	//				}
	//				if (len(string(last_number)) == 4) {
	//					snumber = last_number
	//				}
	//
	//				new_tax_no := "01" + vyear1 + vmonth1 + vday1 + "-" + snumber //เลขที่เอกสารใหม่ส่งสรรพกร
	//
	//				sqlins := `Insert into tax_temp(month_tax,year_tax,doc_date,month_send,day_send,doc_no,tax_no,before_tax_amount,tax_amount,total_amount,create_by,created) values(?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
	//				fmt.Println("Insert tax_temp = ", tax.MonthTax, tax.YearTax, d.DocDate, tax.MonthSend, vAmountPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount, tax.CreateBy)
	//				_, err = db.Exec(sqlins, tax.MonthTax, tax.YearTax, DateAdd, tax.MonthSend, vAmountPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount, tax.CreateBy)
	//				fmt.Println("sqlins", sqlins, d.DocDate, d.DocNo, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount)
	//				if err != nil {
	//					fmt.Println("sqlins", err.Error())
	//					return err
	//				}
	//
	//				last_number1 = last_number1 + 1 //เพิ่มเลขที่บิล
	//
	//			}
	//		}
	//	}
	//
	//}
	//
	//tax.YearTax = strconv.Itoa(BeginDate.Year())
	//tax.MonthTax = strconv.Itoa(int(BeginDate.Month()))
	//tax.MonthSend = SendTaxAmount
	//
	//sqldata := `select doc_date,day_send ,tax_no as doc_no,doc_no as tax_no,before_tax_amount as sum_of_item_amount,tax_amount,total_amount from tax_temp where doc_date between ? and ?`
	//err = db.Select(&tax.ListDoc, sqldata, begindate, enddate)
	//if err != nil {
	//	return err
	//}
	//return nil

}
