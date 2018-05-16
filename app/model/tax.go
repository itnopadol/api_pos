package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
	"strconv"
)

type TaxData struct {
	Id        int     `json:"id" db:"id"`
	MonthTax  int     `json:"month_tax" db:"month"`
	YearTax   int     `json:"year_tax" db:"year"`
	MonthSend float64 `json:"month_send" db:"month_send"`
	ArCode    string  `json:"ar_code" db:"ar_code"`
	ArName    string  `json:"ar_name" db:"ar_name"`
	TaxNo     string  `json:"tax_no" db:"tax_no"`
	Address   string  `json:"address" db:"address"`
	CreateBy  string  `json:"create_by" db:"create_by"`
	ListDoc   []*Data `json:"data"`
}

type Data struct {
	DaySend         float64 `json:"day_send" db:"day_send"`
	DocDate         string  `json:"doc_date" db:"doc_date"`
	DocNo           string  `json:"doc_no" db:"doc_no"`
	BeforeTaxAmount float64 `json:"before_tax_amount" db:"before_tax_amount"`
	TaxAmount       float64 `json:"tax_amount" db:"tax_amount"`
	TotalAmount     float64 `json:"total_amount" db:"total_amount"`
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

	sql := `select count(docDate) as day1 from (select distinct docDate from Queue where  docDate between ? and ?) as q`
	err := db.Get(&vDay, sql, begindate, enddate)
	if err != nil {
		fmt.Println("Count Day =", err.Error())
		return err
	}

	fmt.Println("Count Day =", vDay)

	sqlsum := `select sum(qid) as sumtotal from Queue where  docDate between ? and ?`
	err = db.Get(&vSumAll, sqlsum, begindate, enddate)
	if err != nil {
		fmt.Println("vSumAll =", err.Error())
		return err
	}

	fmt.Println("Sum All = ", vSumAll)

	//dateString := "2018-03-01"
	BeginDate, err := time.Parse("2006-01-02", begindate);
	fmt.Println("begindate,enddate,total", begindate, enddate, vDay, SendTaxAmount)

	fmt.Println("last_number = ", last_number)

	tax.YearTax = BeginDate.Year()
	tax.MonthTax = int(BeginDate.Month())
	tax.MonthSend = SendTaxAmount

	for i := 0; i < vDay; i++ {
		var vTotalDay float64
		var vAmountPerDay float64
		var vPercentDay float64

		DateAdd := BeginDate.AddDate(0, 0, i).Format("2006-01-02")

		sql := `select sum(qId) as totalamount from Queue where  docDate = ?`
		err = db.Get(&vTotalDay, sql, DateAdd)
		if err != nil {
			fmt.Println("vTotal =", err.Error())
			return err
		}

		fmt.Println("total day=", vTotalDay)

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

		sqlsub := `select docDate as doc_date,docNo as doc_no,93 as before_tax_amount,7 as tax_amount,qId as total_amount from Queue where docdate = ? order by docno`
		err = db.Select(&bill, sqlsub, DateAdd)
		fmt.Println("sqlsub = ", sqlsub, DateAdd)
		if err != nil {
			fmt.Println("sqlsub =", err.Error())
			return nil
		}

		last_number1 = 1
		for _, d := range bill {

			var sumtotal float64

			sqlcheck := `select sum(ifnull(totalamount,0)) as sumtotal from Test_Sum_Vat where ifnull(SendDayTax,'') = ?`
			err = db.Get(&sumtotal, sqlcheck, DateAdd)
			if err != nil {
				fmt.Println("sqlcheck", err.Error())
			}

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

				new_tax_no := vyear1 + vmonth1 + vday1 + "-" + snumber //เลขที่เอกสารใหม่ส่งสรรพกร

				sqlins := `Insert into tax_temp(month_tax,year_tax,doc_date,month_send,day_send,doc_no,tax_no,before_tax_amount,tax_amount,total_amount,create_by,created) values(?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
				_, err = db.Exec(sqlins, tax.MonthTax, tax.YearTax, d.DocDate, tax.MonthSend, vAmountPerDay, d.DocNo, new_tax_no, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount, tax.CreateBy)
				fmt.Println("sqlins", sqlins, d.DocDate, d.DocNo, d.BeforeTaxAmount, d.TaxAmount, d.TotalAmount)
				if err != nil {
					fmt.Println("sqlins", err.Error())
					return err
				}

				last_number1 = last_number1 + 1 //เพิ่มเลขที่บิล

			}
		}
	}

	sqldata := `select month_tax,year_tax,doc_date,month_send,day_send,doc_no,tax_no,before_tax_amount,tax_amount,total_amount,create_by,created from tax_temp where doc_date between ? and ?`
	err = db.Select(tax, sqldata, begindate, enddate)
	if err != nil {
		return nil
	}
	return nil

}
