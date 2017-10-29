package model

import (
	"time"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Seller interface {
	Post()
	Save()
}

type SaleMock struct {
}

func (sm *SaleMock) Post() error {
	return nil
}

func (sm *SaleMock) Save() error {
	return nil
}

// Sale เป็นหัวเอกสารขายแต่ละครั้ง
type Sale struct {
	Id       uint64     `json:"id" db:"id"`
	HostId   string   	`json:"host_id" db:"host_id"`
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
	SaleId    uint64  `json:"-" db:"sale_id"`
	Line      uint64  `json:"line"`
	ItemId    uint64  `json:"item_id" db:"item_id"`
	ItemName  string  `json:"item_name" db:"item_name"`
	Price     float64 `json:"price" db:"price"`
	Qty       int     `json:"qty" db:"qty"`
	Unit      string  `json:"unit" db:"unit"`
	Amount	float64 `json:"amount" db:"amount"`
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

	t := time.Now()
	s.Created = &t

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

	fmt.Println("*Sale.Save() start")
	sql1 := `INSERT INTO sale(host_id,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted,create_by,created) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`
		fmt.Println("*Sale.Save()",sql1)
		rs, err := db.Exec(sql1,
		s.HostId,
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

	sql2 := `INSERT INTO sale_sub(sale_id,line,item_id,item_name,price,qty,unit,amount) VALUES(?,?,?,?,?,?,?,?)`
	for _, ss := range s.SaleSubs {
		fmt.Println("start for range s.SaleSubs")
		rs, err = db.Exec(sql2,
			s.Id,
			ss.Line,
			ss.ItemId,
			ss.ItemName,
			ss.Price,
			ss.Qty,
			ss.Unit,
			ss.Amount)
		if err != nil {
			fmt.Printf("Error when db.Exec(sql2) %v\n", err.Error())
			return "", err
		}
		fmt.Println("Insert sale_sub line ", ss)
	}

	}else{
		return "ลูกค้าชำระเงิน ยังไม่ครบกรุณาตรวจสอบ", nil
	}
	fmt.Println("Save data sucess: sale =", s)

	return docno, nil
}

func (s *Sale)SearchSales(db *sqlx.DB) (sales []*Sale, err error){

	sql := `select id,host_id,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted from sale order by created desc limit 10`
	err = db.Select(&sales, sql)
	if err != nil {
		return nil, err
	}

	for _, sub := range sales{
		fmt.Println("SaleID = ",sub.Id)
		sqlsub := `select sale_id,line,item_id,item_name,price,qty,unit,amount from sale_sub where sale_id = ?`
		err = db.Select(&sub.SaleSubs,sqlsub,sub.Id)
		if err != nil {
			return nil, err
		}
	}
	return sales, nil
}


func (s *Sale)SearchSaleById(db *sqlx.DB, id int64) error{
	fmt.Println("ID = ",id)
	sql := `select id,host_id,total_amount,pay_amount,change_amount,type,tax_rate,item_amount,before_tax_amount,tax_amount,is_cancel,is_posted from sale where id = ?`
	err := db.Get(s, sql, id)
	if err != nil {
		return err
	}

	fmt.Println("SaleID = ",s.Id)
	sqlsub := `select sale_id,line,item_id,item_name,price,qty,unit,amount from sale_sub where sale_id = ?`
	err = db.Select(&s.SaleSubs,sqlsub,id)
	if err != nil {
		return  err
	}

	return nil
}