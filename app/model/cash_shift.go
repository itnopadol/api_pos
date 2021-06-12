package model

import (
	//"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	//"strconv"
	//"net"
	//"bufio"
	//"github.com/knq/escpos"
	"bufio"
	"net"
	"strconv"

	"github.com/itnopadol/api_pos/hw"
	"github.com/knq/escpos"
	//"github.com/itnopadol/bc_api/bc_api/config"
	//"github.com/itnopadol/bc_api/bc_api/config"
)

type Shift struct {
	Id             int64     `json:"id" db:"id"`
	CompanyID      int64     `json:"company_id" db:"company_id"`
	BranchID       int64     `json:"branch_id" db:"branch_id"`
	HostCode       string    `json:"host_code" db:"host_code"`
	DocDate        string    `json:"doc_date" db:"doc_date"`
	ChangeBegin    float64   `json:"change_begin" db:"change_begin"`
	ChangeAmount   float64   `json:"change_amount" db:"change_amount"`
	CashAmount     float64   `json:"cash_amount" db:"cash_amount"`
	ExpensesAmount float64   `json:"expenses_amount" db:"expenses_amount"`
	MyDescription  string    `json:"my_description" db:"my_description"`
	IsClosed       int       `json:"is_closed" db:"is_closed"`
	CreatedBy      string    `json:"created_by" db:"created_by"`
	Created        time.Time `json:"created" db:"created"`
	EditedBy       string    `json:"edited_by" db:"edited_by"`
	Edited         time.Time `json:"edited" db:"edited"`
	ClosedBy       string    `json:"closed_by" db:"closed_by"`
	Closed         time.Time `json:"closed" db:"closed"`

	SumChangeBegin    float64 `json:"sum_change_begin" db:"sum_change_begin"`
	SumCashAmount     float64 `json:"sum_cash_amount" db:"sum_cash_amount"`
	SumExpensesAmount float64 `json:"sum_expenses_amount" db:"sum_expenses_amount"`
}

func (ch *Shift) SaveShift(db *sqlx.DB) (string, error) {
	now := time.Now()
	fmt.Println("yyyy-mm-dd date format : ", now.AddDate(0, 0, 0).Format("2006-01-02"))
	DocDate := now.AddDate(0, 0, 0).Format("2006-01-02")

	var get_company_id int64
	var get_branch_id int64

	if ch.CompanyID == 0 {
		get_company_id = 1
	} else {
		get_company_id = ch.CompanyID
	}

	if ch.BranchID == 0 {
		get_branch_id = 1
	} else {
		get_branch_id = ch.BranchID
	}

	ch.DocDate = DocDate
	ch.Created = now

	//fmt.Println("HostCode = ", ch.HostCode)

	var checkCount int
	sqlCheckExist := `select count(host_code) as vCount from cash_shift where host_code = ? and is_closed = 0 ` //and doc_date = ?`
	err := db.Get(&checkCount, sqlCheckExist, ch.HostCode)                                                      //, ch.DocDate)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	//fmt.Println("checkCount =", checkCount)

	if checkCount > 1 {
		return "มีจุดขายที่เปิดมากกว่า 1 จุด", nil
	}

	if checkCount == 0 {
		sql := `INSERT INTO cash_shift(company_id,branch_id,host_code,doc_date,change_begin,change_amount,cash_amount,expenses_amount,my_description,created_by,created) VALUES(?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
		res, err := db.Exec(sql, get_company_id, get_branch_id, ch.HostCode, ch.DocDate, ch.ChangeBegin, ch.ChangeAmount, ch.CashAmount, ch.ExpensesAmount, ch.MyDescription, ch.CreatedBy)
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}

		Id, _ := res.LastInsertId()
		ch.Id = Id

		sqlsub := `UPDATE host set status = 1 where host_code = ?`
		_, err = db.Exec(sqlsub, ch.HostCode)
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}
		// } else {

		// 	sql := `UPDATE cash_shift set is_closed = 0 ,edited_by = ?, edited = CURRENT_TIMESTAMP() where  host_code = ? and doc_date  = ?`
		// 	_, err = db.Exec(sql, ch.EditedBy, ch.HostCode, ch.DocDate)
		// 	if err != nil {
		// 		fmt.Println(err.Error())
		// 		return "", err
		// 	}

		// 	sqlsub := `UPDATE host set status = 1 where host_code = ?`
		// 	_, err = db.Exec(sqlsub, ch.HostCode)
		// 	if err != nil {
		// 		fmt.Println(err.Error())
		// 		return "", err
		// 	}
	}

	return "", nil
}

func (ch *Shift) UpdateShift(db *sqlx.DB) (string, error) {
	now := time.Now()
	fmt.Println("yyyy-mm-dd date format : ", now.AddDate(0, 0, 0).Format("2006-01-02"))

	ch.Edited = now

	var checkCount int
	sqlCheckExist := `select count(host_code) as vCount from cash_shift where host_code = ? and is_closed = 0 ` // and doc_date = ?`
	err := db.Get(&checkCount, sqlCheckExist, ch.HostCode)                                                       //, ch.DocDate)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	if checkCount > 1 {
		return "มีจุดขายที่เปิดมากกว่า 1 จุด", nil
	}
	fmt.Println("Begin Have Docno", ch.HostCode)
	if checkCount != 0 {
		fmt.Println("Have Docno")
		sql := `UPDATE cash_shift set change_begin = ?, change_amount = ?,cash_amount = ?,expenses_amount = ?,my_description=?,edited_by = ?, edited = CURRENT_TIMESTAMP() where  host_code = ? and is_closed = 0` // and doc_date  = ?`
		_, err = db.Exec(sql, ch.ChangeBegin, ch.ChangeAmount, ch.CashAmount, ch.ExpensesAmount, ch.MyDescription, ch.EditedBy, ch.HostCode)                                                     //, ch.DocDate)
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}

		//sqlsub := `UPDATE host set status = 1 where host_code = ?`
		//fmt.Println("SQL Update Host =", sqlsub)
		//_, err = db.Exec(sqlsub, ch.HostCode)
		//if err != nil {
		//	fmt.Println(err.Error())
		//	return err
		//}
	}

	return "", nil
}

func (ch *Shift) ClosedShift(db *sqlx.DB) (string, error) {
	now := time.Now()
	fmt.Println("yyyy-mm-dd date format : ", now.AddDate(0, 0, 0).Format("2006-01-02"))

	ch.Closed = now
	var checkCount int
	sqlCheckExist := `select count(host_code) as vCount from cash_shift where host_code = ? and is_closed = 0` // and doc_date = ?`
	err := db.Get(&checkCount, sqlCheckExist, ch.HostCode)                                                     //, ch.DocDate)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	fmt.Println("docdate = ", ch.DocDate)
	fmt.Println("host_code = ", ch.HostCode)

	if checkCount > 1 {
		return "มีจุดขายที่เปิดมากกว่า 1 จุด", nil
	}

	if checkCount != 0 {
		sql := `UPDATE cash_shift set change_amount = ?,cash_amount = ?,expenses_amount = ?,
				is_closed = ?,closed_by = ?, closed = CURRENT_TIMESTAMP() where  host_code = ? and is_closed = 0`//and doc_date  = ?`
		_, err = db.Exec(sql, ch.ChangeAmount, ch.CashAmount, ch.ExpensesAmount, 1, ch.ClosedBy, ch.HostCode)//, ch.DocDate)
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}

		sqlsub := `UPDATE host set status = 0 where host_code = ?`
		_, err = db.Exec(sqlsub, ch.HostCode)
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}
	}

	return "", nil
}

func (ch *Shift) ShiftDetails(db *sqlx.DB, host_code string, doc_date string) error {

	fmt.Println("doc_date = ", doc_date)
	fmt.Println("host_code = ", host_code)
	sql := `select host_code,doc_date,change_begin,change_amount,cash_amount,expenses_amount,ifnull(my_description,'') as my_description,is_closed,created_by,created from cash_shift where  host_code = ? and doc_date = ?`
	fmt.Println("sql = ", sql)
	err := db.Get(ch, sql, host_code, doc_date)
	if err != nil {
		return err
	}
	return nil
}


func (ch *Shift) ShiftLastID(db *sqlx.DB, host_code string) error {
	fmt.Println("host_code = ", host_code)
	sql := `select id,company_id,branch_id,host_code,doc_date,change_begin,change_amount,cash_amount,expenses_amount,ifnull(my_description,'') as my_description,is_closed,created_by,created from cash_shift where  host_code = ?  and is_closed = 0 order by id desc limit 1`
	fmt.Println("sql = ", sql)
	err := db.Get(ch, sql, host_code)
	if err != nil {
		return err
	}
	return nil
}

func (ch *Shift) SearchShiftByKeyword(db *sqlx.DB, host_code string, doc_date string) (shifts []*Shift, err error) {
	sql := `select host_code,doc_date,change_begin,change_amount,cash_amount,expenses_amount,ifnull(my_description,'') as my_description,is_closed,created_by,created from cash_shift where  host_code = ? order by docdate desc limit 20`
	err = db.Select(&ch, sql, host_code, doc_date)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return shifts, nil
}

func makeline(pt hw.PosPrinter) {
	pt.SetTextSize(0, 0)
	pt.SetFont("A")
	pt.WriteStringLines("-----------------------------------------\n")
}

func (s *Shift) PrintSendDailyTotal(db *sqlx.DB, host_code string, doc_date string) (shifts []*Shift, err error) {
	var sql string
	s.DocDate = doc_date
	s.HostCode = host_code
	fmt.Println("host_code = ", host_code)
	fmt.Println("DOCDATE = ", doc_date, s.DocDate)

	//var today = time.Now()
	//date := fmt.Sprintf("Date %s Time %s", today.Format("02/01/2006"), today.Format("15:04:05"))
	config := new(Config)
	config = GetConfig(db)

	//fmt.Println("printer close shift ",H.PrinterPort)
	f, err := net.Dial("tcp", config.Printer4Port)
	if err != nil {
		fmt.Println("err ", err.Error())
		return nil, err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	p := escpos.New(w)

	pt := hw.PosPrinter{p, w}
	pt.Init()
	pt.SetLeftMargin(40)

	//////////////////////////////////////////////////////////////////////////////////////
	//vDocDate time.Time

	pt.SetCharaterCode(26)
	pt.SetAlign("center")
	pt.SetTextSize(0, 0)
	pt.WriteStringLines("สรุปยอดนำส่งประจำวัน : " + s.DocDate)
	pt.LineFeed()
	pt.SetTextSize(0, 0)
	//pt.PrintRegistrationBitImage(byte(h.LogoImageId), 0)
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	if s.HostCode == "" {
		sql = `select id,host_code,doc_date,change_begin,change_amount,cash_amount,expenses_amount,(select sum(change_begin) from cash_shift where doc_date = a.doc_date) as sum_change_begin,(select sum(cash_amount) from cash_shift where doc_date = a.doc_date) as sum_cash_amount,(select sum(expenses_amount) from cash_shift where doc_date = a.doc_date) as sum_expenses_amount from cash_shift a where doc_date = ? order by host_code`
		err = db.Select(&shifts, sql, s.DocDate)
	} else {
		sql = `select id,host_code,doc_date,change_begin,change_amount,cash_amount,expenses_amount,(select sum(change_begin) from cash_shift where host_code = a.host_code and doc_date = a.doc_date) as sum_change_begin,(select sum(cash_amount) from cash_shift where host_code = a.host_code and doc_date = a.doc_date) as sum_cash_amount,(select sum(expenses_amount) from cash_shift where host_code = a.host_code and doc_date = a.doc_date) as sum_expenses_amount from cash_shift a where host_code = ? and doc_date = ? order by host_code;
 `
		err = db.Select(&shifts, sql, s.HostCode, s.DocDate)
	}
	fmt.Println("sql = ", sql, s.HostCode, s.DocDate)
	if err != nil {
		fmt.Println("error sql ", err.Error())
		return nil, err
	}

	pt.SetFont("B")
	//pt.SetAlign("left")
	pt.WriteStringLines(" จุดขาย")
	pt.WriteStringLines("     ")
	pt.WriteStringLines("มูลค่าเงินสด")
	pt.WriteStringLines("      ")
	pt.WriteStringLines("มูลค่าเงินทอน")
	pt.WriteStringLines("     ")
	pt.WriteStringLines("มูลค่าสุทธิ" + "\n")
	//pt.FormfeedN(3)
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	var vLineNumber int
	var vNetAmount float64
	var vSumNetAmount float64

	for _, s := range shifts {
		vLineNumber = vLineNumber + 1
		vNetAmount = s.CashAmount - s.ChangeBegin

		fmt.Println("Cash =", vLineNumber, s.CashAmount)
		fmt.Println("Expense =", vLineNumber, s.ExpensesAmount)

		pt.SetFont("A")
		pt.WriteStringLines(strconv.Itoa(vLineNumber) + "." + s.HostCode)
		pt.WriteStringLines("     " + CommaFloat(s.CashAmount))
		pt.WriteStringLines("      " + CommaFloat(s.ChangeBegin))
		pt.WriteStringLines("     " + CommaFloat(vNetAmount) + "\n")
		//pt.FormfeedN(3)
	}
	makeline(pt)
	////////////////////////////////////////////////////////////////////////////////////

	vSumNetAmount = shifts[0].SumCashAmount - shifts[0].SumChangeBegin

	fmt.Println("SumCashAmount = ", CommaFloat(shifts[0].SumCashAmount))
	pt.SetFont("B")
	pt.WriteStringLines("รวมเป็นเงิน ")
	pt.SetFont("A")
	pt.WriteStringLines(" ")
	pt.WriteStringLines(CommaFloat(shifts[0].SumCashAmount))
	pt.WriteStringLines("      ")
	pt.WriteStringLines(CommaFloat(shifts[0].SumChangeBegin))
	pt.WriteStringLines("     ")
	pt.WriteStringLines(CommaFloat(vSumNetAmount) + "\n")
	makeline(pt)
	pt.Cut()
	pt.End()

	return shifts, nil
}
