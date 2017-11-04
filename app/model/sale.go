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
	CancelDatetime 	*time.Time `json:"cancel_datetime" db:"cancel_datetime"`

	SaleSubs []*SaleSub `json:"sale_subs"`
}

// SaleSub เป็นรายการสินค้าที่ขายใน Sale
type SaleSub struct {
	SaleId    int  `json:"-" db:"sale_id"`
	ItemId    int `json:"item_id" db:"item_id"`
	ItemName  string  `json:"item_name" db:"item_name"`
	ShortName	string `json:"short_name" db:"short_name"`
	Price     float64 `json:"price" db:"price"`
	Qty       int     `json:"qty" db:"qty"`
	Unit      string  `json:"unit" db:"unit"`
	Amount	float64 `json:"amount" db:"amount"`
	IsKitchen  int `json:"is_kitchen" db:"is_kitchen"`
	IsAtHome	int `json:"is_athome" db:"is_athome"`
	Line int `json:"line" db:"line"`

}

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

		s.TotalAmount = toFixed(s.TotalAmount,2)
		vTaxAmount = toFixed(s.TotalAmount-((s.TotalAmount*100)/(100+float64(s.TaxRate))),2)
		vBeforeTaxAmount = toFixed(s.TotalAmount-vTaxAmount,2)

		s.BeforeTaxAmount = vBeforeTaxAmount
		s.TaxAmount = vTaxAmount

		s.CreateBy = "somrod"

		s.DocNo = GenDocno(db,s.HostCode)
		s.QueId = LastQueId(db)

	fmt.Println("*Sale.Save() start")
	sql1 := `INSERT INTO sale(host_code,que_id,doc_no,doc_date,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted,create_by,created) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
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
			s.CreateBy,
			s.Created)
	if err != nil {
		fmt.Printf("Error when db.Exec(sql1) %v", err.Error())
		return "", err
	}
	id, _ := rs.LastInsertId()
	s.Id = uint64(id)
	fmt.Println("s.MachineId =", s.Id)

	sql2 := `INSERT INTO sale_sub(sale_id,line,item_id,item_name,short_name,price,qty,unit,amount,is_kitchen,is_athome) VALUES(?,?,?,?,?,?,?,?,?,?,?)`
	for _, ss := range s.SaleSubs {
		fmt.Println("start for range s.SaleSubs")
		rs, err = db.Exec(sql2,
			s.Id,
			ss.Line,
			ss.ItemId,
			ss.ItemName,
			ss.ShortName,
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
	}
	//พิมพ์ บิล และ ใบจัดสินค้า
	config := new(Config)
	config = GetConfig(db)
	//err = PrintBill(s, config, db)
	err = printPickup(s, config, db)

	}else{
		return "ลูกค้าชำระเงิน ยังไม่ครบกรุณาตรวจสอบ", nil
	}
	fmt.Println("Save data sucess: sale =", s)

	return s.DocNo, nil
}

func (s *Sale)SearchSales(db *sqlx.DB,host_code string,doc_date string,keyword string) (sales []*Sale, err error){

	sql := `select id,host_code,que_id,doc_no,doc_date,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted from sale where host_code = ? and doc_date = ? and (doc_no like CONCAT("%",?,"%")) order by created desc`
	err = db.Select(&sales, sql, host_code, doc_date, keyword)
	if err != nil {
		return nil, err
	}

	for _, sub := range sales{
		fmt.Println("SaleID = ",sub.Id)
		sqlsub := `select sale_id,line,item_id,item_name,ifnull(short_name,'') as short_name,price,qty,unit,amount,is_kitchen,is_athome from sale_sub where sale_id = ?`
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
	sqlsub := `select sale_id,line,item_id,item_name,ifnull(short_name,'') as short_name,price,qty,unit,amount,is_kitchen,is_athome from sale_sub where sale_id = ? `
	err = db.Select(&s.SaleSubs,sqlsub,id)
	if err != nil {
		return  err
	}

	return nil
}

func GetConfig(db *sqlx.DB)(config *Config){
	cf := new(Config)
	sql := `select ifnull(company_name,'') as company_name,ifnull(tax_id,'') as tax_id,ifnull(tax_rate,0) as tax_rate,ifnull(printer1_port,'') as printer1_port,ifnull(printer2_port,'') as printer2_port,ifnull(printer3_port,'') as printer3_port from config`
	fmt.Println("Config = ",sql)
	err := db.Get(cf,sql)
	if err != nil {
		fmt.Println(err.Error())
	}
	//printerIP = config.Printer2Port

	config = cf
	return config
}

var vQueID int
var printerIP string

func PrintBill(s *Sale, c *Config, db *sqlx.DB)error{

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
	for _, sub := range s.SaleSubs {
		var vLineNumber int

		vLineNumber = sub.Line+1
		pt.SetFont("B")
		pt.WriteStringLines("     "+strconv.Itoa(vLineNumber)+"."+sub.ItemName )
		pt.WriteStringLines("     "+strconv.Itoa(sub.Qty)+" "+sub.Unit)
		pt.WriteStringLines("     ")
		pt.WriteStringLines("     "+strconv.FormatFloat(sub.Price, 'f', -1, 64))
		pt.WriteStringLines("     ")
		pt.WriteStringLines("     "+strconv.FormatFloat(sub.Amount, 'f', -1, 64)+"\n")
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

	return nil
}

func printPickup(s *Sale, c *Config, db *sqlx.DB)error{

	f, err := net.Dial("tcp", c.Printer2Port)

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

	//////////////////////////////////////////////////////////////////////////////////////////////////
	pt.SetCharaterCode(26)
	pt.SetAlign("center")
	pt.SetTextSize(1, 1)
	pt.WriteStringLines("คิวเลขที่ : "+strconv.Itoa(s.QueId))
	pt.LineFeed()
	pt.SetTextSize(0, 1)
	pt.WriteStringLines("Kitchen Slip")

	pt.LineFeed()
	pt.SetTextSize(0, 0)
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	pt.SetFont("B")
	pt.SetAlign("left")
	pt.WriteStringLines("  รายการสินค้า" )
	pt.WriteStringLines("   ")
	pt.WriteStringLines("     จำนวน/หน่วย")
	pt.WriteStringLines("    ")
	pt.WriteStringLines("   สถานะ\n" )
	pt.FormfeedN(3)
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////

	for _, sub := range s.SaleSubs {
		var vLineNumber int
		var vAtHome string

		if (sub.IsAtHome==1){
			vAtHome = "Home"
		}else{
			vAtHome = "-"
		}
		if (sub.IsKitchen==1) {
			vLineNumber = sub.Line + 1
			pt.SetTextSize(0, 1)
			pt.SetFont("A")
			pt.SetAlign("left")
			pt.WriteStringLines("   " + strconv.Itoa(vLineNumber) + "." + sub.ShortName)
			pt.WriteStringLines("   ")
			pt.WriteStringLines("     " + strconv.Itoa(sub.Qty) + " " + sub.Unit)
			pt.WriteStringLines("        " + vAtHome + "\n")
			pt.FormfeedN(3)
			pt.SetTextSize(1, 1)
		}
	}
	makeline(pt)
	///////////////////////////////////////////////////////////////////////////////////
	pt.SetFont("B")
	pt.SetAlign("center")
	pt.Formfeed()
	pt.Write("*** Completed ***")
	pt.Formfeed()
	pt.Cut()
	pt.End()

	return nil
}

//func makeline(pt hw.PosPrinter) {
//	pt.SetTextSize(0,0)
//	pt.SetFont("A")
//	pt.WriteStringLines("-----------------------------------------\n")
//}

//func (s *Sale)PrintSaleDaily(db *sqlx.DB, doc_date string, h *Host)(sales []*Sale, err error){
//	s.DocDate = doc_date
//
//	f, err := net.Dial("tcp", "192.168.0.206:9100")
//
//	if err != nil {
//		panic(err)
//	}
//	defer f.Close()
//
//	w := bufio.NewWriter(f)
//	p := escpos.New(w)
//
//	pt := hw.PosPrinter{p,w}
//	pt.Init()
//	pt.SetLeftMargin(20)
//	//pt.PrintRegistrationBitImage(0, 0)
//	pt.WriteRaw([]byte{29,	40,	76,	6,	0,	48,	85,	32,	32,10,10 })
//
//	//////////////////////////////////////////////////////////////////////////////////////
//	pt.SetCharaterCode(26)
//	pt.SetAlign("center")
//	pt.SetTextSize(1, 1)
//	pt.WriteStringLines("สรุปยอดขายประจำวัน : "+s.DocDate)
//	pt.LineFeed()
//	pt.SetTextSize(0, 0)
//	pt.PrintRegistrationBitImage(byte(h.LogoImageId), 0)
//	makeline(pt)
//	///////////////////////////////////////////////////////////////////////////////////
//	sql := `select id,host_code,que_id,doc_no,doc_date,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted from sale where doc_date = ?  order by id`
//	err = db.Select(&sales, sql, doc_date)
//	if err != nil {
//		return nil, err
//	}
//
//	pt.SetFont("B")
//	pt.WriteStringLines("   รายการสินค้า" )
//	pt.WriteStringLines("     ")
//	pt.WriteStringLines("   จำนวน/หน่วย")
//	pt.WriteStringLines("  ")
//	pt.WriteStringLines("   ราคา" )
//	pt.WriteStringLines("    ")
//	pt.WriteStringLines("   รวม\n" )
//	pt.FormfeedN(3)
//	makeline(pt)
//	///////////////////////////////////////////////////////////////////////////////////
//	var vLineNumber int
//	for _, sale := range sales{
//
//		vLineNumber = vLineNumber+1
//		pt.SetFont("B")
//		pt.WriteStringLines("     "+strconv.Itoa(vLineNumber)+"."+sale.DocNo )
//		pt.WriteStringLines("     "+strconv.Itoa(vLineNumber)+"."+sale.DocNo )
//		pt.WriteStringLines("     "+strconv.FormatFloat(sale.TotalAmount, 'f', -1, 64)+"\n")
//		pt.FormfeedN(3)
//	}
//	makeline(pt)
//	////////////////////////////////////////////////////////////////////////////////////
//	pt.SetFont("B")
//	pt.WriteStringLines("รวมเป็นเงิน ")
//	pt.WriteStringLines("                                   ")
//	pt.WriteStringLines(strconv.FormatFloat(s.TotalAmount, 'f', -1, 64)+" บาท\n")
//	makeline(pt)
//	// Footer Area
//	pt.SetFont("A")
//	pt.SetAlign("center")
//	pt.WriteStringLines("รหัสผ่าน Wifi : 999999999")
//	pt.Formfeed()
//	pt.Write("*** Completed ***")
//	pt.Formfeed()
//	pt.Cut()
//	pt.End()
//
//	return nil, nil
//}

//func PrintPickup(s *Sale, db *sqlx.DB)error{
//	const (
//		printerIP = "192.168.0.206:9100"
//	)
//
//	f, err := net.Dial("tcp", printerIP)
//
//	if err != nil {
//		panic(err)
//	}
//	defer f.Close()
//
//	w := bufio.NewWriter(f)
//	p := escpos.New(w)
//
//	pt := hw.PosPrinter{p,w}
//	pt.Init()
//	pt.SetLeftMargin(20)
//	//pt.PrintRegistrationBitImage(0, 0)
//	pt.WriteRaw([]byte{29,	40,	76,	6,	0,	48,	85,	32,	32,10,10 })
//	printBill(s, pt, db)
//	//printPickup(s, pt)
//
//	//pt.WriteRaw([]byte{27,112,0,25,250})
//	//pt.Pulse()
//	return nil
//}



//func printBill(s *Sale, pt hw.PosPrinter, db *sqlx.DB){
//	fmt.Println("Print Bill")
//
//	config := new(Config)
//
//	sql := `select ifnull(company_name,'') as company_name,ifnull(tax_id,'') as tax_id,ifnull(tax_rate,0) as tax_rate,ifnull(printer1_port,'') as printer1_port,ifnull(printer2_port,'') as printer2_port,ifnull(printer3_port,'') as printer3_port from config`
//	fmt.Println("Config = ",sql)
//	err := db.Get(config,sql)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//
//	fmt.Println(config.CompanyName)
//
//
//	pt.SetCharaterCode(26)
//	pt.SetAlign("center")
//	pt.SetTextSize(1, 1)
//	pt.WriteStringLines("คิวเลขที่ : "+strconv.Itoa(s.QueId))
//	pt.LineFeed()
//	pt.SetTextSize(0, 0)
//	/////////////////////////////////////////////////////////////////////////////////////
//	pt.SetFont("B")
//	pt.WriteStringLines(config.CompanyName+"\n")
//	pt.SetAlign("left")
//	pt.WriteStringLines("เลขประจำตัวผู้เสียภาษี : "+config.TaxId)
//	pt.SetAlign("right")
//
//	pt.WriteStringLines("	Cashier : "+s.CreateBy)
//
//	pt.WriteStringLines("     วันที่ :"+s.Created.Format("02-01-2006 15:04:05"))
//	pt.WriteStringLines("   เลขที่ : "+s.DocNo)
//	pt.WriteStringLines("   Pos Id : "+s.HostId)
//	pt.LineFeed()
//	makeline(pt)
//	///////////////////////////////////////////////////////////////////////////////////
//	for _, sub := range s.SaleSubs {
//		var vLineNumber int
//
//		vLineNumber = sub.Line+1
//		pt.SetFont("B")
//		pt.WriteStringLines("   "+strconv.Itoa(vLineNumber)+"."+sub.ItemName )
//		pt.WriteStringLines("		")
//		pt.WriteStringLines("   "+strconv.Itoa(sub.Qty)+" "+sub.Unit+"\n")
//		pt.FormfeedN(3)
//	}
//	makeline(pt)
//	////////////////////////////////////////////////////////////////////////////////////
//	pt.SetFont("B")
//	pt.WriteStringLines("รวมเป็นเงิน ")
//	pt.WriteStringLines("				")
//	pt.WriteStringLines(strconv.FormatFloat(s.TotalAmount, 'f', -1, 64)+" บาท\n")
//	makeline(pt)
//	// Footer Area
//	pt.SetFont("A")
//	pt.SetAlign("center")
//	pt.WriteStringLines("รหัสผ่าน Wifi : 999999999")
//	pt.Formfeed()
//	pt.Write("*** Completed ***")
//	pt.Formfeed()
//	pt.Cut()
//	pt.End()
//}



//func printHeader(pt hw.PosPrinter) {
//	pt.SetCharaterCode(26)
//	pt.SetAlign("center")
//	pt.SetTextSize(1, 1)
//	pt.WriteStringLines("คิวเลขที่ 35")
//	pt.LineFeed()
//	pt.SetTextSize(0, 0)
//}
//
//func printKitchenHeader(pt hw.PosPrinter) {
//	pt.SetCharaterCode(26)
//	pt.SetAlign("center")
//	pt.SetTextSize(1, 1)
//	pt.WriteStringLines("คิวเลขที่ 35")
//	pt.LineFeed()
//	pt.WriteStringLines("Kitchen Slip")
//
//	pt.LineFeed()
//	pt.SetTextSize(0, 0)
//	makeline(pt)
//}


//func printCompanyInfo(s *Sale, pt hw.PosPrinter) {
//	pt.SetFont("B")
//	pt.WriteStringLines("===== นพดลพานิช จำกัด =====\n")
//	pt.SetAlign("left")
//	pt.WriteStringLines("เลขประจำตัวผู้เสียภาษี 999999999999")
//	pt.SetAlign("right")
//
//	pt.WriteStringLines("	Cashier : XXXX\n")
//
//	pt.WriteStringLines("วันที่ : 01/09/2017 09:34น.")
//	pt.WriteStringLines("   เลขที่ : 0120171005-001")
//	pt.LineFeed()
//	makeline(pt)
//}

//func printFooter(s *Sale, pt hw.PosPrinter) {
//	pt.SetFont("B")
//	pt.WriteStringLines("รวมเป็นเงิน ")
//	pt.WriteStringLines("				")
//	pt.WriteStringLines(strconv.FormatFloat(s.TotalAmount, 'f', -1, 64)+" บาท\n")
//	makeline(pt)
//	// Footer Area
//	pt.SetFont("A")
//	pt.SetAlign("center")
//	pt.WriteStringLines("รหัสผ่าน Wifi : 999999999")
//	pt.Formfeed()
//	pt.Write("*** Completed ***")
//	pt.Formfeed()
//	pt.Cut()
//	pt.End()
//}
//func printDetail(s *Sale, pt hw.PosPrinter) {
//	pt.SetFont("B")
//	pt.WriteStringLines("    1. ปูนเสือ")
//	pt.WriteStringLines("		")
//	pt.WriteStringLines("	1 ชิ้น\n")
//	//pt.LineFeed()
//	pt.WriteStringLines("    2. ปูนช้าง")
//	pt.WriteStringLines("		")
//	pt.WriteStringLines("	1 ชิ้น\n")
//	pt.WriteStringLines("    3. น้ำยาเชื่อมท่อ")
//	pt.WriteStringLines("		")
//	pt.WriteStringLines("	1 ชิ้น\n")
//	pt.FormfeedN(3)
//	makeline(pt)
//}

//func printKitchenDetail(pt hw.PosPrinter) {
//	pt.SetTextSize(1,1)
//	pt.SetAlign("left")
//	pt.WriteStringLines("1. CAPU")
//	pt.WriteStringLines(" x ")
//	pt.WriteStringLines(" 1 \n")
//	pt.LineFeed()
//	pt.WriteStringLines("2. ESP ")
//	pt.WriteStringLines(" x ")
//	pt.WriteStringLines(" 1 \n")
//	pt.FormfeedN(3)
//	pt.SetTextSize(1,1)
//	makeline(pt)
//
//}


//func printKitchenFooter(pt hw.PosPrinter) {
//	// Footer Area
//	pt.SetFont("B")
//	pt.SetAlign("center")
//	pt.Formfeed()
//	pt.Write("*** Completed ***")
//	pt.Formfeed()
//	pt.Cut()
//	pt.End()
//}

