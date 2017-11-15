package model

import (
	"time"
	"fmt"
	"github.com/jmoiron/sqlx"
	"bufio"
	"github.com/knq/escpos"
	"net"
	"github.com/itnopadol/api_pos/hw"
	"strconv"
	"strings"
)

// Sale เป็นหัวเอกสารขายแต่ละครั้ง
type Sale struct {
	Id       uint64     `json:"id" db:"id"`
	HostCode   string   	`json:"host_code" db:"host_code"`
	QueId	int `json:"que_id" db:"que_id"`
	DocNo string `json:"doc_no" db:"doc_no"`
	DocDate string `json:"doc_date" db:"doc_date"`
	TotalAmount    float64    `json:"total_amount" db:"total_amount"`
	PayAmount      float64    `json:"pay_amount" db:"pay_amount"`
	ChangeAmount   float64    `json:"change_amount" db:"change_amount"`
	Type     string     `json:"type" db:"type"`
	TaxRate  int `json:"tax_rate" db:"tax_rate"`
	ItemAmount float64 `json:"item_amount" db:"item_amount"`
	BeforeTaxAmount float64 `json:"before_tax_amount" db:"before_tax_amount"`
	TaxAmount float64 `json:"tax_amount" db:"tax_amount"`
	CreateBy  string `json:"create_by" db:"create_by"`
	Created  *time.Time `json:"-" db:"created"`

	IsPosted bool       `json:"-" db:"is_posted"`
	IsCancel		int        `json:"is_cancel" db:"is_cancel"`
	PostedBy		string `json:"posted_by" db:"posted_by"`
	PostedDatetime 	*time.Time `json:"posted_datetime" db:"posted_datetime"`
	CancelBy 		string `json:"cancel_by" db:"cancel_by"`
	Canceled 	*time.Time `json:"canceled" db:"canceled"`

	SumCashAmount float64 `json:"sum_cash_amount" db:"sum_cash_amount"`
	SumChangeAmount float64 `json:"sum_change_amount" db:"sum_change_amount"`
	SumCashAmountAll float64 `json:"sum_cash_amount_all" db:"sum_cash_amount_all"`
	SumChangeAmountAll float64 `json:"sum_change_amount_all" db:"sum_change_amount_all"`
	NetAmount float64 `json:"net_amount" db:"net_amount"`
	NetAmountAll float64 `json:"net_amount_all" db:"net_amount_all"`
	BillCount int `json:"bill_count" db:"bill_count"`
	BillCountAll int `json:"bill_count_all" db:"bill_count_all"`

	SaleSubs []*SaleSub `json:"sale_subs"`
}

// SaleSub เป็นรายการสินค้าที่ขายใน Sale
type SaleSub struct {
	SaleId    int  `json:"-" db:"sale_id"`
	ItemId    int `json:"item_id" db:"item_id"`
	ItemName  string  `json:"item_name" db:"item_name"`
	ShortName	string `json:"short_name" db:"short_name"`
	Description string `json:"description" db:"description"`
	Price     float64 `json:"price" db:"price"`
	Qty       int     `json:"qty" db:"qty"`
	Unit      string  `json:"unit" db:"unit"`
	Amount	float64 `json:"amount" db:"amount"`
	IsKitchen  int `json:"is_kitchen" db:"is_kitchen"`
	IsAtHome	int `json:"is_athome" db:"is_athome"`
	Line int `json:"line" db:"line"`

}

var vQueID int
var printerIP string

func (s *Sale) ShowChangeAmount()(Amount float64,msg string, err error){
	var Remain float64
	var Change float64

	Change = s.PayAmount - s.TotalAmount //20-20 = 0
	fmt.Println(s.TotalAmount,"    ",s.PayAmount)
	if Change > 0 {
		s.ChangeAmount = Change //80
	} else {
		s.ChangeAmount = 0 //0
	}

	Remain = (s.TotalAmount - s.PayAmount) + s.ChangeAmount //(20-100)+ 80 //(50-20)+0 //(20-20)+0

	switch  {
	case Remain < 0 :
		msg = "Change Amount ="
		Amount = Change
	case Remain == 0 && Change > 0:
		msg = "Change Amount ="
		Amount = Change
	case Remain == 0 && Change == 0:
		msg = "IsOK"
		Amount = Change
	case Remain > 0 :
		msg = "Remain Amount ="
		Amount = Remain
	}

	return Amount, msg, nil
}

func (s *Sale) CheckAmount()(status int, err error){
	var TotalAmount float64
	var ItemAmount float64
	var Amount float64

	TotalAmount = s.TotalAmount

	for _, sub := range s.SaleSubs {
		Amount = sub.Amount
		ItemAmount = ItemAmount+Amount
	}

	if (TotalAmount != ItemAmount){
		status = 1
	}else {
		status = 0
	}

	fmt.Println("TotalAmount =",TotalAmount)
	fmt.Println("ItemAmount", ItemAmount)
	fmt.Println("status", status)

	return status, nil
}

func (s *Sale) SaleSave(db *sqlx.DB) (docno string, err error) {
	var CheckRemain float64
	var CheckChange float64

	now := time.Now()
	s.Created = &now

	DocDate := now.AddDate(0, 0, 0).Format("2006-01-02")

	s.DocDate = DocDate

	CheckChange = s.PayAmount - s.TotalAmount
	if CheckChange > 0 {
		s.ChangeAmount = CheckChange
	} else {
		s.ChangeAmount = 0
	}

	fmt.Println("Total = ",s.TotalAmount,"Pay = ",s.PayAmount)
	CheckRemain = (s.TotalAmount - s.PayAmount) + s.ChangeAmount //(100-50)+0

	fmt.Println("Total = ",s.TotalAmount,"Pay = ",s.PayAmount, "Remain = ",CheckRemain, "Change =",s.ChangeAmount)

	var vTaxAmount float64
	var vBeforeTaxAmount float64

	if CheckRemain == 0 {

		fmt.Println("Host_Code = ",s.HostCode)
		if(s.HostCode != ""){
			checkAmount, _ := s.CheckAmount()
			if err != nil {
				return "Error Check ItemAmount", err
			}

			if(checkAmount==0){
				s.TotalAmount = toFixed(s.TotalAmount,2)
				vTaxAmount = toFixed(s.TotalAmount-((s.TotalAmount*100)/(100+float64(s.TaxRate))),2)
				vBeforeTaxAmount = toFixed(s.TotalAmount-vTaxAmount,2)

				s.BeforeTaxAmount = vBeforeTaxAmount
				s.TaxAmount = vTaxAmount

				s.DocNo = GenDocno(db,s.HostCode)
				s.QueId = LastQueId(db)

				fmt.Println("*Sale.Save() start")
				sql1 := `INSERT INTO sale(host_code,que_id,doc_no,doc_date,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted,create_by,created) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
				fmt.Println("*Sale.Save()",sql1)
				rs, err := db.Exec(sql1,
					s.HostCode,
					s.QueId,
					s.DocNo,
					s.DocDate,
					s.TotalAmount,
					s.PayAmount,
					s.ChangeAmount,
					s.Type,
					s.TaxRate,
					s.ItemAmount,
					s.BeforeTaxAmount,
					s.TaxAmount,
					s.IsCancel,
					s.IsPosted,
					s.CreateBy)
				if err != nil {
					fmt.Printf("Error when db.Exec(sql1) %v", err.Error())
					return "", err
				}
				id, _ := rs.LastInsertId()
				s.Id = uint64(id)
				fmt.Println("s.MachineId =", s.Id)

				var checkPrintSlipKitchen int
				var checkPrintSlipBar int
				checkPrintSlipKitchen = 0
				checkPrintSlipBar = 0

				sql2 := `INSERT INTO sale_sub(sale_id,line,item_id,item_name,short_name,description,price,qty,unit,amount,is_kitchen,is_athome) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`
				for _, ss := range s.SaleSubs {
					fmt.Println("start for range s.SaleSubs")
					rs, err = db.Exec(sql2,
						s.Id,
						ss.Line,
						ss.ItemId,
						ss.ItemName,
						ss.ShortName,
						ss.Description,
						ss.Price,
						ss.Qty,
						ss.Unit,
						ss.Amount,
						ss.IsKitchen,
						ss.IsAtHome)
					if err != nil {
						fmt.Printf("Error when db.Exec(sql2) %v\n", err.Error())
						return "", err
					}
					fmt.Println("Insert sale_sub line ", ss)


					if (ss.IsKitchen == 1){
						checkPrintSlipKitchen = checkPrintSlipKitchen+1
					}

					if (ss.IsKitchen == 0){
						checkPrintSlipBar = checkPrintSlipBar+1
					}
				}

				fmt.Println("checkPrintSlipKitchen = ",checkPrintSlipKitchen)
				//พิมพ์ บิล และ ใบจัดสินค้า
				config := new(Config)
				config = GetConfig(db)

				host := new(Host)
				host = GetHostPrinter(db, s.HostCode)

				fmt.Println("Port1 ",host.HostCode,host.PrinterPort)
				fmt.Println("Port2 " ,config.Printer2Port)

				if (host.PrinterPort != "") {
					err = PrintBill(s, host, config, db)
				}
				if (config.Printer1Port != "") {
					if (checkPrintSlipKitchen > 0) {
						err = printPickup1(s, config, db)
					}
				}
				if (config.Printer2Port != "") {
					if (checkPrintSlipBar > 0) {
						err = printPickup2(s, config, db)
					}
				}

			}else{
				return "มูลค่ารวม ไม่เท่ากับ มูลค่าสินค้า กรุณาตรวจสอบ", err
			}

		}else{
			return "Host Code ไม่แสดง กรุณาตรวจสอบ", err
		}

	}else{
		return "ลูกค้าชำระเงิน ยังไม่ครบกรุณาตรวจสอบ", nil
	}
	fmt.Println("Save data sucess: sale =", s)

	return s.DocNo, nil
}

func (s *Sale) SaleVoid(db *sqlx.DB) error {
	var checkCount int
	sqlCheckExist := `select count(doc_no) as vCount from sale where id = ?`
	err := db.Get(&checkCount, sqlCheckExist, s.Id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if(checkCount!=0) {
		fmt.Println("*Sale.Save() start")
		sql := `Update sale set is_cancel = 1, cancel_by = ?, canceled = CURRENT_TIMESTAMP() where id = ?)`
		fmt.Println("*Sale.Save()", sql)
		_, err = db.Exec(sql, s.Id)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		sqlsub := `Update sale_sub set is_cancel = 1 where sale_id = ?)`
		fmt.Println("*Sale.Save()", sql)
		_, err = db.Exec(sqlsub, s.Id)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func (s *Sale)SearchSales(db *sqlx.DB,host_code string,doc_date string,keyword string) (sales []*Sale, err error){

	sql := `select id,host_code,que_id,doc_no,doc_date,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted from sale where host_code = ? and doc_date = ? and (doc_no like CONCAT("%",?,"%")) order by created desc`
	err = db.Select(&sales, sql, host_code, doc_date, keyword)
	if err != nil {
		return nil, err
	}

	for _, sub := range sales{
		fmt.Println("SaleID = ",sub.Id)
		sqlsub := `select sale_id,line,item_id,item_name,ifnull(short_name,'') as short_name,ifnull(description,'') as description,price,qty,unit,amount,is_kitchen,is_athome from sale_sub where sale_id = ?`
		err = db.Select(&sub.SaleSubs,sqlsub,sub.Id)
		if err != nil {
			return nil, err
		}
	}
	return sales, nil
}

func (s *Sale)SearchSaleById(db *sqlx.DB, id int64) error{
	fmt.Println("ID = ",id)
	sql := `select id,host_code,que_id,doc_no,doc_date,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted from sale where id = ? order by id desc limit 20`
	err := db.Get(s, sql, id)
	if err != nil {
		return err
	}

	fmt.Println("SaleID = ",s.Id)
	sqlsub := `select sale_id,line,item_id,item_name,ifnull(short_name,'') as short_name,ifnull(description,'') as description,price,qty,unit,amount,is_kitchen,is_athome from sale_sub where sale_id = ? `
	err = db.Select(&s.SaleSubs,sqlsub,id)
	if err != nil {
		return  err
	}

	return nil
}

func PrintBill(s *Sale, h *Host, c *Config, db *sqlx.DB)error{
	myPassword := genMikrotikPassword(c)
	fmt.Println("password =",myPassword)

	f, err := net.Dial("tcp", h.PrinterPort)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	p := escpos.New(w)

	pt := hw.PosPrinter{p,w}
	pt.Init()
	pt.SetLeftMargin(10)
	//pt.PrintRegistrationBitImage(0, 0)
	//pt.WriteRaw([]byte{29,	40,	76,	6,	0,	48,	85,	32,	32,10,10 })
	//pt.WriteRaw([]byte{28, 112, 1, 0})
	//pt.WriteRaw([]byte{28, 112, 1, 1})

	//////////////////////////////////////////////////////////////////////////////////////
	pt.SetCharaterCode(26)
	pt.SetAlign("center")
	pt.SetTextSize(1, 1)
	pt.WriteStringLines("คิวเลขที่ : "+strconv.Itoa(s.QueId))
	pt.LineFeed()
	pt.SetTextSize(0, 0)
	pt.SetAlign("center")
	pt.SetFont("B")
	pt.WriteStringLines(c.CompanyName+"\n")
	pt.SetAlign("center")
	pt.WriteStringLines(c.Address+"\n")
	fmt.Println(c.Address)
	pt.SetAlign("center")
	pt.WriteStringLines("ใบเสร็จรับเงิน/ใบกำกับภาษีอย่างย่อ\n")
	pt.WriteStringLines("เลขประตัวผู้เสียภาษี "+c.TaxId+"\n")
	//pt.WriteStringLines("ใบกำกับภาษีอย่างย่อ\n")
	pt.SetAlign("center")
	pt.WriteStringLines(" เลขเครื่อง : "+s.HostCode+"      "+" พนักงาน : "+s.CreateBy +"\n")
	//pt.WriteStringLines(" พนักงาน : "+s.CreateBy+"\n")
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	var CountItem int
	var CountQty int
	for _, subcount := range s.SaleSubs {
		CountItem = CountItem+1
		CountQty = CountQty+subcount.Qty
	}

	fmt.Println("CountItem =",CountItem, CountQty )
	///////////////////////////////////////////////////////////////////////////////////
	pt.SetAlign("left")
	for _, sub := range s.SaleSubs {
		var vAtHome string
		var vDiffEmpty int
		var vDiffOld int
		var vItemPriceAmount string
		var vPrice float64

		if (sub.IsAtHome==1){
			vAtHome = "H"
		}else{
			vAtHome = ""
		}
		pt.SetFont("A")
		if (sub.Description==""){
			pt.WriteStringLines(" "+sub.ItemName+"\n")
		}else{
			pt.WriteStringLines(" "+sub.ItemName+" ("+sub.Description+" )"+"\n")
		}

		vPrice = sub.Amount / float64(sub.Qty)
		vItemPriceAmount = " "+strconv.FormatFloat(vPrice, 'f', -1, 64)+" X "+strconv.Itoa(sub.Qty)+" "+sub.Unit

		vLen := len(vItemPriceAmount)
		vDiff := 25- (vLen/3)

		if (vDiff < 0){
			vDiffEmpty = 0
		}else {
			vDiffEmpty = vDiff
		}

		vDiffOld = vDiffEmpty

		fmt.Println("ItemName=",sub.ItemName)
		fmt.Println("Len",vLen/3)
		fmt.Println("Diff ",vDiff)
		if (sub.Line == 0 ) {
			pt.WriteStringLines(vItemPriceAmount + strings.Repeat(" ", vDiffEmpty))
		}else{
			pt.WriteStringLines(vItemPriceAmount + strings.Repeat(" ", vDiffOld))
		}
		pt.WriteStringLines("      ")
		pt.WriteStringLines(CommaFloat(sub.Amount)+"  "+vAtHome+"\n")
		//pt.FormfeedN(3)
	}
	makeline(pt)
	////////////////////////////////////////////////////////////////////////////////////
	pt.SetFont("A")
	pt.WriteStringLines(" "+strconv.Itoa(CountItem)+" รายการ "+strconv.Itoa(CountQty)+" ชิ้น\n")
	pt.WriteStringLines(" รวม ")
	pt.WriteStringLines("                              ")
	//pt.WriteStringLines(strconv.FormatFloat(s.TotalAmount, 'f', 2, 64)+"\n")
	pt.WriteStringLines(CommaFloat(s.TotalAmount)+"\n")
	////////////////////////////////////////////////////////////////////////////////////
	pt.SetFont("A")
	//pt.WriteStringLines(" มูลค่าสินค้ามีภาษีมูลค่าเพิ่ม"+"                       "+Commaf(vBeforeTaxAmount)+"\n")
	//pt.WriteStringLines(" ภาษีมูลค่าเพิ่ม"+strconv.Itoa(c.TaxRate)+"%"+"                                "+Commaf(vTaxAmount)+"\n")
	pt.WriteStringLines(" ชำระด้วยเงินสด" +"                      "+CommaFloat(s.PayAmount)+"\n")
	pt.WriteStringLines(" เงินทอน" +"                            "+CommaFloat(s.ChangeAmount)+"\n")
	pt.WriteStringLines(" เลขที่ :"+ s.DocNo+"  วันที่ :"+ s.Created.Format("02-01-2006 15:04:05")+"\n")
	//pt.WriteStringLines(" วันที่ :"+ s.Created.Format("02-01-2006 15:04:05")+"\n")
	pt.SetAlign("center")
	pt.WriteStringLines(" H = สั่งกลับบ้าน"+"\n")
	pt.WriteStringLines(" ราคานี้รวมภาษีมูลค่าเพิ่มแล้ว"+"\n")
	makeline(pt)
	 //Footer Area
	pt.SetFont("A")
	pt.SetAlign("center")
	fmt.Println("myPassword After = ",myPassword)
	pt.WriteStringLines(myPassword)
	pt.Formfeed()
	pt.Cut()
	pt.OpenCashBox()
	pt.End()

	return nil
}

func PrintInvoice(s *Sale, c *Config, db *sqlx.DB)error{

	f, err := net.Dial("tcp", c.Printer1Port)

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
	pt.WriteRaw([]byte{28, 112, 1, 0})

	//////////////////////////////////////////////////////////////////////////////////////
	pt.SetCharaterCode(26)
	pt.SetAlign("center")
	pt.SetTextSize(1, 1)
	pt.WriteStringLines("คิวเลขที่ : "+strconv.Itoa(s.QueId))
	pt.LineFeed()
	pt.SetTextSize(0, 0)
	/////////////////////////////////////////////////////////////////////////////////////
	pt.SetFont("B")
	pt.WriteStringLines(c.CompanyName+"\n")
	pt.SetAlign("left")
	pt.WriteStringLines("เลขประจำตัวผู้เสียภาษี : "+c.TaxId)
	pt.SetAlign("right")

	pt.WriteStringLines("	Cashier : "+s.CreateBy)

	pt.WriteStringLines("       วันที่ :"+s.Created.Format("02-01-2006 15:04:05"))
	pt.WriteStringLines("   เลขที่ : "+s.DocNo)
	pt.WriteStringLines("      Pos Id : "+s.HostCode+"\n")
	pt.LineFeed()
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////

	pt.SetFont("B")
	pt.WriteStringLines("  ")
	pt.WriteStringLines("รายการ" )
	pt.WriteStringLines("                   ")
	pt.WriteStringLines("จำนวน/หน่วย")
	pt.WriteStringLines("  ")
	pt.WriteStringLines("ราคา" )
	pt.WriteStringLines("   ")
	pt.WriteStringLines("รวม\n" )
	pt.FormfeedN(3)
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	var vLineNumber int
	var vDiffEmpty int
	for _, sub := range s.SaleSubs {

		vLen := len(sub.ItemName)
		vDiff := 25- (vLen/3)

		if (vDiff < 0){
			vDiffEmpty = 0
		}else {
			vDiffEmpty = vDiff
		}

		vDiffOld := vDiffEmpty

		fmt.Println("ItemName=",sub.ItemName)
		fmt.Println("Len",vLen/3)
		fmt.Println("Diff ",vDiff)

		if (sub.Line == 0 ) {
			vLineNumber = sub.Line + 1
			pt.SetFont("B")
			pt.WriteStringLines(strconv.Itoa(vLineNumber) + "." + sub.ItemName + strings.Repeat(" ", vDiffEmpty))
			pt.WriteStringLines("  " + strconv.Itoa(sub.Qty) + " " + sub.Unit)
			pt.WriteStringLines("    ")
			pt.WriteStringLines(strconv.FormatFloat(sub.Price, 'f', -1, 64))
			pt.WriteStringLines("    ")
			pt.WriteStringLines(strconv.FormatFloat(sub.Amount, 'f', -1, 64) + "\n")
			pt.FormfeedN(3)
		}else{
			vLineNumber = sub.Line + 1
			pt.SetFont("B")
			pt.WriteStringLines(strconv.Itoa(vLineNumber) + "." + sub.ItemName + strings.Repeat(" ", vDiffOld))
			pt.WriteStringLines("  " + strconv.Itoa(sub.Qty) + " " + sub.Unit)
			pt.WriteStringLines("    ")
			pt.WriteStringLines(strconv.FormatFloat(sub.Price, 'f', -1, 64))
			pt.WriteStringLines("    ")
			pt.WriteStringLines(strconv.FormatFloat(sub.Amount, 'f', -1, 64) + "\n")
			pt.FormfeedN(3)
		}

		//+strings.Repeat(" ",15)
	}
	makeline(pt)
	////////////////////////////////////////////////////////////////////////////////////
	pt.SetFont("B")
	pt.WriteStringLines("รวมเป็นเงิน ")
	pt.WriteStringLines("                                 ")
	pt.WriteStringLines(strconv.FormatFloat(s.TotalAmount, 'f', -1, 64)+" บาท\n")
	makeline(pt)
	// Footer Area
	//pt.SetFont("A")
	//pt.SetAlign("center")
	//pt.WriteStringLines("รหัสผ่าน Wifi : 999999999")
	//pt.Formfeed()
	//pt.Write("*** Completed ***")
	//pt.Formfeed()
	pt.Cut()
	pt.OpenCashBox()
	pt.End()

	return nil
}

//ใบหยิบห้องครัว
func printPickup1(s *Sale, c *Config, db *sqlx.DB)error{
	fmt.Println("c.Printer2Port",c.Printer2Port)

	f, err := net.Dial("tcp", c.Printer1Port)

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
	//pt.WriteRaw([]byte{29,	40,	76,	6,	0,	48,	85,	32,	32,10,10 })

	//////////////////////////////////////////////////////////////////////////////////////////////////
	pt.SetCharaterCode(26)
	pt.SetAlign("center")
	pt.SetTextSize(1, 1)
	pt.WriteStringLines("คิวเลขที่ : "+strconv.Itoa(s.QueId))
	pt.LineFeed()
	pt.SetTextSize(0, 1)
	pt.WriteStringLines("Kitchen Slip"+"\n")
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	//pt.LineFeed()
	pt.SetTextSize(0, 0)
	pt.SetFont("A")
	pt.SetAlign("left")
	pt.WriteStringLines("   ")
	pt.WriteStringLines("รายการ" )
	pt.WriteStringLines("               ")
	pt.WriteStringLines("จำนวน")
	pt.WriteStringLines("     ")
	pt.WriteStringLines("สถานะ\n" )
	//pt.FormfeedN(3)
	pt.SetAlign("center")
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////

	pt.SetAlign("left")
	for _, sub := range s.SaleSubs {
		var vAtHome string

		if (sub.IsAtHome==1){
			vAtHome = "H"
		}else{
			vAtHome = ""
		}
		if (sub.IsKitchen==1) {
			pt.SetTextSize(1, 1)
			pt.SetFont("A")
			pt.SetAlign("left")
			if (sub.Description != ""){
				pt.WriteStringLines(sub.ShortName+"("+sub.Description+")")
			}else{
				pt.WriteStringLines(sub.ShortName)
			}
			pt.WriteStringLines(" " + strconv.Itoa(sub.Qty))
			pt.WriteStringLines(" " + vAtHome + "\n")
			//pt.FormfeedN(3)
			pt.SetTextSize(1, 1)
		}
	}
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	//pt.SetFont("B")
	//pt.SetAlign("center")
	//pt.Formfeed()
	//pt.Write("*** Completed ***")
	//pt.Formfeed()
	pt.Cut()
	//pt.OpenCashBox()
	pt.End()

	return nil
}

//ใบหยิบ บาร์น้ำ
func printPickup2(s *Sale, c *Config, db *sqlx.DB)error{
	fmt.Println("c.Printer2Port",c.Printer3Port)

	f, err := net.Dial("tcp", c.Printer2Port)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	p := escpos.New(w)

	pt := hw.PosPrinter{p,w}
	pt.Init()
	pt.SetLeftMargin(10)
	//pt.PrintRegistrationBitImage(0, 0)
	//pt.WriteRaw([]byte{29,	40,	76,	6,	0,	48,	85,	32,	32,10,10 })

	//////////////////////////////////////////////////////////////////////////////////////////////////
	pt.SetCharaterCode(26)
	pt.SetAlign("center")
	pt.SetTextSize(1, 1)
	pt.WriteStringLines("คิวเลขที่ : "+strconv.Itoa(s.QueId))
	pt.LineFeed()
	pt.SetTextSize(0, 1)
	pt.WriteStringLines("Bar Slip"+"\n")
	pt.SetAlign("center")
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	//pt.LineFeed()
	pt.SetTextSize(0, 0)
	pt.SetFont("A")
	pt.SetAlign("left")
	pt.WriteStringLines("   ")
	pt.WriteStringLines("รายการ" )
	pt.WriteStringLines("               ")
	pt.WriteStringLines("จำนวน")
	pt.WriteStringLines("     ")
	pt.WriteStringLines("สถานะ\n" )
	//pt.FormfeedN(3)
	pt.SetAlign("center")
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////

	for _, sub := range s.SaleSubs {
		var vAtHome string

		if (sub.IsAtHome==1){
			vAtHome = "H"
		}else{
			vAtHome = ""
		}
		if (sub.IsKitchen==0) {
			pt.SetTextSize(1, 0)
			pt.SetFont("A")
			pt.SetAlign("left")
			if (sub.Description != ""){
				pt.WriteStringLines(sub.ItemName+"("+sub.Description+")")
			}else{
				pt.WriteStringLines(sub.ItemName)
			}
			pt.WriteStringLines(" " + strconv.Itoa(sub.Qty))
			pt.WriteStringLines(" " + vAtHome + "\n")
			//pt.FormfeedN(3)
			pt.SetTextSize(1, 1)
		}
	}
	pt.SetAlign("center")
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	//pt.SetFont("B")
	//pt.SetAlign("center")
	//pt.Formfeed()
	//pt.Write("*** Completed ***")
	//pt.Formfeed()
	pt.Cut()
	//pt.OpenCashBox()
	pt.End()

	return nil
}

func (s *Sale)PrintSaleDailyTotal(db *sqlx.DB, host_code string, doc_date string)(sales []*Sale, err error){
	var sql string

	s.DocDate = doc_date
	s.HostCode = host_code

	fmt.Println("DOCDATE = ",doc_date,s.DocDate)

	config := new(Config)
	config = GetConfig(db)

	f, err := net.Dial("tcp", config.Printer1Port)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	p := escpos.New(w)

	pt := hw.PosPrinter{p,w}
	pt.Init()
	pt.SetLeftMargin(20)

	//////////////////////////////////////////////////////////////////////////////////////
	pt.WriteRaw([]byte{28, 112, 1, 0})
	pt.SetCharaterCode(26)
	pt.SetAlign("center")
	pt.SetTextSize(0, 0)
	pt.WriteStringLines("สรุปยอดขายประจำวัน : "+s.DocDate)
	pt.LineFeed()
	pt.SetTextSize(0, 0)
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	if(s.HostCode == ""){
		sql = `select distinct host_code,doc_date,
			(select count(doc_no) from sale where doc_date = a.doc_date and is_cancel = 0) as bill_count_all,
			(select count(doc_no) from sale where host_code = a.host_code and doc_date = a.doc_date and is_cancel = 0) as bill_count,
			(select sum(pay_amount) from sale where doc_date = a.doc_date and is_cancel = 0) as sum_cash_amount_all,
			(select sum(change_amount) from sale where doc_date = a.doc_date and is_cancel = 0) as sum_change_amount_all,
			(select sum(pay_amount) from sale where doc_date = a.doc_date and is_cancel = 0)- (select sum(change_amount) from sale where doc_date = a.doc_date and is_cancel = 0) as net_amount_all,
			(select sum(pay_amount) from sale where doc_date = a.doc_date and host_code = a.host_code and is_cancel = 0 group by host_code,doc_date) as sum_cash_amount,
			(select sum(change_amount) from sale where doc_date = a.doc_date and host_code = a.host_code and is_cancel = 0 group by host_code,doc_date) as sum_change_amount,
			(select sum(pay_amount) from sale where doc_date = a.doc_date and host_code = a.host_code and is_cancel = 0 group by host_code,doc_date) - (select sum(change_amount) from sale where doc_date = a.doc_date and host_code = a.host_code and is_cancel = 0 group by host_code,doc_date) as net_amount
			from sale a where doc_date = ? and is_cancel = 0 order by host_code`
		err = db.Select(&sales, sql, doc_date)
	}else{
		sql = `	select distinct host_code,doc_date,
		    (select count(doc_no) from sale where doc_date = a.doc_date and is_cancel = 0) as bill_count_all,
			(select count(doc_no) from sale where host_code = a.host_code and doc_date = a.doc_date and is_cancel = 0) as bill_count,
			(select sum(pay_amount) from sale where host_code = a.host_code and doc_date = a.doc_date and is_cancel = 0) as sum_cash_amount_all,
			(select sum(change_amount) from sale where host_code = a.host_code and doc_date = a.doc_date and is_cancel = 0) as sum_change_amount_all,
			(select sum(pay_amount) from sale where host_code = a.host_code and doc_date = a.doc_date and is_cancel = 0)- (select sum(change_amount) from sale where host_code = a.host_code and doc_date = a.doc_date and is_cancel = 0) as net_amount_all,
			(select sum(pay_amount) from sale where host_code = a.host_code and doc_date = a.doc_date and host_code = a.host_code and is_cancel = 0 group by host_code,doc_date) as sum_cash_amount,
			(select sum(change_amount) from sale where host_code = a.host_code and doc_date = a.doc_date and host_code = a.host_code and is_cancel = 0 group by host_code,doc_date) as sum_change_amount,
			(select sum(pay_amount) from sale where host_code = a.host_code and doc_date = a.doc_date and host_code = a.host_code and is_cancel = 0 group by host_code,doc_date) - (select sum(change_amount) from sale where doc_date = a.doc_date and host_code = a.host_code and is_cancel = 0 group by host_code,doc_date) as net_amount
			from sale a where host_code = ? and doc_date = ? and is_cancel = 0 order by host_code`
		err = db.Select(&sales, sql, host_code, doc_date)
	}

	fmt.Println("sql = ",sql, host_code, doc_date)
	if err != nil {
		return nil, err
	}
	fmt.Println("Sale Data ",sales[0].SumCashAmount)
	pt.SetFont("B")
	pt.WriteStringLines("   จุดขาย")
	pt.WriteStringLines("  ")
	pt.WriteStringLines("   มูลค่าเงินสด" )
	pt.WriteStringLines("  ")
	pt.WriteStringLines("   มูลค่าเงินทอน" )
	pt.WriteStringLines("   ")
	pt.WriteStringLines("   มูลค่าขายสุทธิ\n" )
	pt.FormfeedN(3)
	makeline(pt)
	/////////////////////////////////////////////////////////////////////////////////
	var vLineNumber int
	for _, s := range sales{
		pt.SetAlign("left")
		vLineNumber = vLineNumber+1
		pt.SetFont("B")
		pt.WriteStringLines("    "+strconv.Itoa(vLineNumber)+"."+s.HostCode)
		pt.WriteStringLines("      "+CommaFloat(s.SumCashAmount))
		pt.WriteStringLines("        "+CommaFloat(s.SumChangeAmount))
		pt.WriteStringLines("        "+CommaFloat(s.NetAmount)+"\n")
		pt.FormfeedN(3)
	}
	makeline(pt)
	pt.SetAlign("left")
	pt.SetFont("B")
	if(s.HostCode==""){
		pt.WriteStringLines("จำนวนบิลทั้งหมด "+strconv.Itoa(sales[0].BillCountAll)+" บิล\n")
	}else{
		pt.WriteStringLines("จำนวนบิลทั้งหมด "+strconv.Itoa(sales[0].BillCount)+" บิล\n")
	}

	//////////////////////////////////////////////////////////////////////////////////
	makeline(pt)
	fmt.Println("SumCashAmount = ",sales[0].SumCashAmount )
	pt.SetFont("B")
	pt.WriteStringLines("รวมเป็นเงิน ")
	pt.WriteStringLines("    ")
	pt.WriteStringLines(CommaFloat(sales[0].SumCashAmountAll))
	pt.WriteStringLines("        ")
	pt.WriteStringLines(CommaFloat(sales[0].SumChangeAmountAll))
	pt.WriteStringLines("        ")
	pt.WriteStringLines(CommaFloat(sales[0].NetAmountAll)+"\n")
	makeline(pt)
	pt.Cut()
	pt.End()

	return nil, nil
}
