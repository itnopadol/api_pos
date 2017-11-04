package model

import (
	"time"
	"github.com/jmoiron/sqlx"
	"fmt"
	"strconv"
	"net"
	"bufio"
	"github.com/knq/escpos"
	"github.com/itnopadol/api_pos/hw"
)

type Shift struct {
	Id int64 `json:"id" db:"id"`
	HostCode string `json:"host_code" db:"host_code"`
	DocDate string `json:"doc_date" db:"doc_date"`
	ChangeBegin float64 `json:"change_begin" db:"change_begin"`
	ChangeAmount float64 `json:"change_amount" db:"change_amount"`
	CashAmount float64 `json:"cash_amount" db:"cash_amount"`
	ExpensesAmount float64 `json:"expenses_amount" db:"expenses_amount"`
	MyDescription string `json:"my_description" db:"my_description"`
	IsClosed int `json:"is_closed" db:"is_closed"`
	CreatedBy string `json:"created_by" db:"created_by"`
	Created time.Time `json:"created" db:"created"`
	EditedBy string `json:"edited_by" db:"edited_by"`
	Edited time.Time `json:"edited" db:"edited"`
	ClosedBy string `json:"closed_by" db:"closed_by"`
	Closed time.Time `json:"closed" db:"closed"`
	Sale []*Sale `json:"sale"`
}


func (ch *Shift)SaveShift(db *sqlx.DB) error{
	now := time.Now()
	fmt.Println("yyyy-mm-dd date format : ", now.AddDate(0, 0, 0).Format("2006-01-02"))
	DocDate := now.AddDate(0, 0, 0).Format("2006-01-02")

	ch.DocDate = DocDate
	ch.Created = now

	var checkCount int
	sqlCheckExist := `select count(host_code) as vCount from cash_shift where doc_date = ?`
	err := db.Get(&checkCount, sqlCheckExist, ch.DocDate)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if (checkCount == 0){
		sql := `INSERT INTO cash_shift(host_code,doc_date,change_begin,change_amount,cash_amount,expenses_amount,my_description,created_by,created) VALUES(?,?,?,?,?,?,?,?,?)`
		res, err := db.Exec(sql, ch.HostCode, ch.DocDate, ch.ChangeBegin, ch.ChangeAmount, ch.CashAmount, ch.ExpensesAmount, ch.MyDescription, ch.CreatedBy, ch.Created)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		Id, _ := res.LastInsertId()
		ch.Id = Id

		sqlsub := `UPDATE host set status = 1 where host_code = ?`
		_, err = db.Exec(sqlsub, ch.HostCode)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}else{

	}

	return nil
}

func (ch *Shift)UpdateShift(db *sqlx.DB) error{
	now := time.Now()
	fmt.Println("yyyy-mm-dd date format : ", now.AddDate(0, 0, 0).Format("2006-01-02"))

	ch.Edited = now

	var checkCount int
	sqlCheckExist := `select count(host_code) as vCount from cash_shift where doc_date = ?`
	err := db.Get(&checkCount, sqlCheckExist, ch.DocDate)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Begin Have Docno",ch.HostCode)
	if (checkCount != 0) {
		fmt.Println("Have Docno")
		sql := `UPDATE cash_shift set change_begin = ?, change_amount = ?,cash_amount = ?,expenses_amount = ?,my_description=?,edited_by = ?, edited = ? where  host_code = ? and doc_date  = ?`
		_, err = db.Exec(sql,ch.ChangeBegin, ch.ChangeAmount, ch.CashAmount, ch.ExpensesAmount, ch.MyDescription, ch.EditedBy, ch.Edited, ch.HostCode, ch.DocDate)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		sqlsub := `UPDATE host set status = 1 where host_code = ?`
		fmt.Println("SQL Update Host =",sqlsub)
		_, err = db.Exec(sqlsub, ch.HostCode)
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

	fmt.Println("docdate = ",ch.DocDate)
	fmt.Println("host_code = ",ch.HostCode)

	var checkCount int
	sqlCheckExist := `select count(host_code) as vCount from cash_shift where doc_date = ?`
	err := db.Get(&checkCount, sqlCheckExist, ch.DocDate)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if (checkCount != 0) {
		sql := `UPDATE cash_shift set change_amount = ?,cash_amount = ?,expenses_amount = ?, is_closed = ?,closed_by = ?, closed = ? where  host_code = ? and doc_date  = ?`
		_, err = db.Exec(sql, ch.ChangeAmount, ch.CashAmount, ch.ExpensesAmount, 1, ch.ClosedBy, ch.Closed, ch.HostCode, ch.DocDate)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		sqlsub := `UPDATE host set status = 0 where host_code = ?`
		_, err = db.Exec(sqlsub, ch.HostCode)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func (ch *Shift)ShiftDetails(db *sqlx.DB, host_code string, doc_date string) error{

	fmt.Println("doc_date = ",doc_date)
	fmt.Println("host_code = ",host_code)
	sql := `select host_code,doc_date,change_begin,change_amount,cash_amount,expenses_amount,ifnull(my_description,'') as my_description,is_closed,created_by,created from cash_shift where  host_code = ? and doc_date = ?`
	fmt.Println("sql = ",sql)
	err := db.Get(ch, sql, host_code, doc_date)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}


func (ch *Shift)SearchShiftByKeyword(db *sqlx.DB, host_code string, doc_date string)(shifts []*Shift, err error){
	sql := `select host_code,doc_date,change_begin,change_amount,cash_amount,expenses_amount,ifnull(my_description,'') as my_description,is_closed,created_by,created from cash_shift where  host_code = ? order by docdate desc limit 20`
	err = db.Select(&ch, sql, host_code, doc_date)
	if err != nil {
		fmt.Println(err.Error())
		return  nil, err
	}
	return shifts, nil
}


func makeline(pt hw.PosPrinter) {
	pt.SetTextSize(0,0)
	pt.SetFont("A")
	pt.WriteStringLines("-----------------------------------------\n")
}

func (s *Sale)PrintSaleDaily(db *sqlx.DB, doc_date string, h *Host)(sales []*Sale, err error){
	s.DocDate = doc_date

	f, err := net.Dial("tcp", "192.168.0.206:9100")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	p := escpos.New(w)

	pt := hw.PosPrinter{p,w}
	pt.Init()
	pt.SetLeftMargin(20)
	//pt.PrintRegistrationBitImage(0, 0)
	pt.WriteRaw([]byte{29,	40,	76,	6,	0,	48,	85,	32,	32,10,10 })

	//////////////////////////////////////////////////////////////////////////////////////
	pt.SetCharaterCode(26)
	pt.SetAlign("center")
	pt.SetTextSize(1, 1)
	pt.WriteStringLines("สรุปยอดขายประจำวัน : "+s.DocDate)
	pt.LineFeed()
	pt.SetTextSize(0, 0)
	pt.PrintRegistrationBitImage(byte(h.LogoImageId), 0)
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	sql := `select id,host_code,que_id,doc_no,doc_date,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted from sale where doc_date = ?  order by id`
	err = db.Select(&sales, sql, doc_date)
	if err != nil {
		return nil, err
	}

	pt.SetFont("B")
	pt.WriteStringLines("   รายการสินค้า" )
	pt.WriteStringLines("     ")
	pt.WriteStringLines("   จำนวน/หน่วย")
	pt.WriteStringLines("  ")
	pt.WriteStringLines("   ราคา" )
	pt.WriteStringLines("    ")
	pt.WriteStringLines("   รวม\n" )
	pt.FormfeedN(3)
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	var vLineNumber int
	for _, sale := range sales{

		vLineNumber = vLineNumber+1
		pt.SetFont("B")
		pt.WriteStringLines("     "+strconv.Itoa(vLineNumber)+"."+sale.DocNo )
		pt.WriteStringLines("     "+strconv.Itoa(vLineNumber)+"."+sale.DocNo )
		pt.WriteStringLines("     "+strconv.FormatFloat(sale.TotalAmount, 'f', -1, 64)+"\n")
		pt.FormfeedN(3)
	}
	makeline(pt)
	////////////////////////////////////////////////////////////////////////////////////
	pt.SetFont("B")
	pt.WriteStringLines("รวมเป็นเงิน ")
	pt.WriteStringLines("                                   ")
	pt.WriteStringLines(strconv.FormatFloat(s.TotalAmount, 'f', -1, 64)+" บาท\n")
	makeline(pt)
	// Footer Area
	pt.SetFont("A")
	pt.SetAlign("center")
	pt.WriteStringLines("รหัสผ่าน Wifi : 999999999")
	pt.Formfeed()
	pt.Write("*** Completed ***")
	pt.Formfeed()
	pt.Cut()
	pt.End()

	return nil, nil
}