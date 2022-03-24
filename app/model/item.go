package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	// "bufio"
	// "github.com/knq/escpos"
	// "net"
	"time"
)

type Item struct {
	Id        int64
	Code      string       `json:"code" db:"code"`
	Name      string       `json:"name" db:"name"`
	ShortName string       `json:"short_name" db:"short_name"`
	NameEn    string       `json:"name_en,omitempty" db:"name_en"`
	NameCn    string       `json:"name_cn,omitempty" db:"name_cn"`
	Price     float32      `json:"price" db:"price"`
	Unit      string       `json:"unit"`
	UnitEn    string       `json:"unit_en,omitempty" db:"unit_en"`
	UnitCn    string       `json:"unit_cn,omitempty" db:"unit_cn"`
	MenuId    uint64       `json:"menu_id,omitempty" db:"menu_id"`
	MenuSeq   int          `json:"menu_seq,omitempty" db:"menu_seq"`
	Image     string       `json:"image" db:"image"`
	IsKitchen int          `json:"is_kitchen" db:"is_kitchen"`
	Active    int          `json:"active" db:"active"`
	CreatedBy string       `json:"created_by" db:"created_by"`
	Created   *time.Time   `json:"created" db:"created"`
	EditedBy  string       `json:"edited_by" db:"edited_by"`
	Edited    *time.Time   `json:"edited" db:"edited"`
	Prices    []*PricesSub `json:"prices_sub"`
}

type PricesSub struct {
	Id     int64     `json:"id"`
	ItemId int64   `json:"item_id" db:"item_id"`
	Name   string  `json:"name" db:"name"`
	NameEn string  `json:"name_en" db:"name_en"`
	NameCn string  `json:"name_cn" db:"name_cn"`
	Price1 float32 `json:"price" db:"price"`
	Active int     `json:"active" db:"active"`
}

func (i *Item) Get(db *sqlx.DB, id int64) (err error) {
	sql := `SELECT id,code,ifnull(short_name,'') as short_name,ifnull(name,'') as name, ifnull(name_en,'') as name_en, ifnull(name_cn,'') as name_cn, ifnull(unit,'') as unit, ifnull(unit_en,'') as unit_en, ifnull(unit_cn,'') as unit_cn, ifnull(menu_seq,0) as menu_seq, ifnull(menu_id,0) as menu_id, ifnull(image,'') as image, ifnull(price,0) as price, ifnull(active,0) as active, ifnull(is_kitchen,0) as is_kitchen, ifnull(created_by,'') as created_by FROM item WHERE  id = ?`
	fmt.Println("Item = ", sql, id)
	err = db.Get(i, sql, id)

	//err = db.QueryRowx(sql, id).StructScan(i)
	if err != nil {
		return err
	}

	// ดึงข้อมูลราคาทั้งหมดของสินค้ารายการนี้
	sizes := []*PricesSub{}
	sql = `SELECT * FROM price_sub WHERE  item_id = ?`
	fmt.Println("Price = ", sql)
	err = db.Select(&sizes, sql, id)
	if err != nil {
		return err
	}

	i.Prices = sizes
	return nil
}

func (i *Item) ByMenuId(db *sqlx.DB, id int64) ([]*Lang, error) {
	fmt.Println("call method: Item.ByMenuId::lang:", langs)
	var sql string
	langInit()
	for _, l := range langs {
		fmt.Println("Lang1 = ", l.Name)
		items := []*Item{}
		switch l.Id {
		case 1:
			sql = `SELECT id,code,short_name,name, unit, menu_seq, menu_id, image, price, active, is_kitchen FROM item WHERE active = 1 and  menu_id = ? order by code`
		case 2:
			sql = `SELECT id,code,short_name,name_en as name, unit_en as unit, menu_seq, menu_id, image, price, active, is_kitchen FROM item WHERE active = 1 and  menu_id = ? order by code`
		case 3:
			sql = `SELECT id,code,short_name,name_cn as name, unit_cn as unit, menu_seq, menu_id, image, price, active, is_kitchen FROM item WHERE active = 1 and  menu_id = ? order by code`
		}
		fmt.Println("case:", l.Id, l.Name, sql, id)
		err := db.Select(&items, sql, id)
		if err != nil {
			fmt.Println("error = ", err.Error())
			log.Println("Error select Items")
			return nil, err
		}
		fmt.Println("Items:", items)
		// query Size{}
		for _, i := range items {
			prices := []*PricesSub{}
			sql = `SELECT * FROM price_sub WHERE active = 1 and  item_id = ?`
			item_id := int(i.Id)
			fmt.Println("item_id =", item_id)
			err = db.Select(&prices, sql, item_id)
			if err != nil {
				return nil, err
			}
			i.Prices = prices
		}
		l.MenuId = int(id)
		l.Items = items
	}
	return langs, nil
}

//func (i *Item) SaveItem(db *sqlx.DB) error {
//	var checkCount int
//	sqlCheckExist := `select count(id) as vCount from item where code = ?`
//	err := db.Get(&checkCount, sqlCheckExist, i.Code)
//	if err != nil {
//		fmt.Println(err.Error())
//		return err
//	}
//
//	if (checkCount == 0) {
//		sql := `INSERT INTO item(code,short_name,name, name_en, unit, unit_en, menu_id, menu_seq, image, price, active, is_kitchen, created_by, created) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP())`
//		rs, err := db.Exec(sql, i.Code, i.ShortName, i.Name, i.NameEn, i.Unit, i.UnitEn, i.MenuId, i.MenuSeq, i.Image, i.Price, 1, i.IsKitchen, i.CreatedBy)
//		if err != nil {
//			return err
//		}
//		id, _ := rs.LastInsertId()
//		i.Id = id
//		fmt.Println("Item Id : ", id)
//	}
//
//	return nil
//}

func (i *Item) SaveItem(db *sqlx.DB) error {
	var checkCount int
	var checkCountSub int

	sqlCheckExist := `select count(id) as vCount from item where code = ?`
	err := db.Get(&checkCount, sqlCheckExist, i.Code)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if (checkCount == 0) {
		sql := `INSERT INTO item(code,short_name,name, name_en, unit, unit_en, menu_id, menu_seq, image, price, active, is_kitchen, created_by, created) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP())`
		rs, err := db.Exec(sql, i.Code, i.ShortName, i.Name, i.NameEn, i.Unit, i.UnitEn, i.MenuId, i.MenuSeq, i.Image, i.Price, 1, i.IsKitchen, i.CreatedBy)
		if err != nil {
			return err
		}
		id, _ := rs.LastInsertId()
		i.Id = id
		fmt.Println("Item Id : ", id)


		for _, sub := range i.Prices {

			sqlCheckSubExist := `select count(id) as vCount from price_sub where item_id = ? and name = ?`
			err := db.Get(&checkCountSub, sqlCheckSubExist, i.Id, sub.Name)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			if (checkCountSub == 0) {
				sql_sub := `INSERT INTO price_sub(item_id, name, name_en, name_cn, price, active) VALUES(?, ?, ?, ?, ?, 1)`
				rs_sub, err := db.Exec(sql_sub, i.Id, sub.Name, sub.NameEn, sub.NameCn, sub.Price1)
				if err != nil {
					return err
				}

				item_id, _ := rs_sub.LastInsertId()
				sub.Id = item_id
				fmt.Println("Item_sub Id : ", item_id)
			}

		}
	}


	return nil
}

//func (i *Item) UpdateItem(db *sqlx.DB) error {
//	var checkCount int
//	sqlCheckExist := `select count(id) as vCount from item where code = ?`
//	err := db.Get(&checkCount, sqlCheckExist, i.Code)
//	if err != nil {
//		fmt.Println(err.Error())
//		return err
//	}
//	fmt.Println("Count : ", checkCount)
//	fmt.Println("ID = ", i.Id)
//	if (checkCount != 0) {
//		sql := `UPDATE item set code=?, short_name=?, name=?, name_en=?, unit=?, unit_en=?, menu_id=?, menu_seq=?, image=?, price=?, active=?, is_kitchen=?, edited_by=?, edited = CURRENT_TIMESTAMP() where id = ?`
//		_, err := db.Exec(sql, i.Code, i.ShortName, i.Name, i.NameEn, i.Unit, i.UnitEn, i.MenuId, i.MenuSeq, i.Image, i.Price, i.Active, i.IsKitchen, i.EditedBy, i.Id)
//		fmt.Println("sql = ", sql)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}


func (i *Item) UpdateItem(db *sqlx.DB) error {
	var checkCount int
	var checkCountSub int

	sqlCheckExist := `select count(id) as vCount from item where id = ?`
	err := db.Get(&checkCount, sqlCheckExist, i.Id)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("Count : ", checkCount)
	fmt.Println("ID = ", i.Id)
	if (checkCount != 0) {
		sql := `UPDATE item set code=?, short_name=?, name=?, name_en=?, unit=?, unit_en=?, menu_id=?, menu_seq=?, image=?, price=?, active=?, is_kitchen=?, edited_by=?, edited = CURRENT_TIMESTAMP() where id = ?`
		_, err := db.Exec(sql, i.Code, i.ShortName, i.Name, i.NameEn, i.Unit, i.UnitEn, i.MenuId, i.MenuSeq, i.Image, i.Price, i.Active, i.IsKitchen, i.EditedBy, i.Id)
		fmt.Println("sql = ", sql)
		if err != nil {
			return err
		}

		fmt.Println("len price",len(i.Prices))
		for _, sub := range i.Prices{

			sqlCheckSubExist := `select count(id) as vCount from price_sub where item_id = ? and name = ?`
			fmt.Println("sqlCheckSubExist : ", sqlCheckSubExist, i.Id, sub.Name)
			err := db.Get(&checkCountSub, sqlCheckSubExist, i.Id, sub.Name)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			fmt.Println("checkCountSub", checkCountSub)
			if (checkCountSub == 0){
				sql_sub := `INSERT INTO price_sub(item_id, name, name_en, name_cn, price, active) VALUES(?, ?, ?, ?, ?, 1)`
				fmt.Println("sql_sub",sql_sub)
				rs_sub, err := db.Exec(sql_sub, i.Id, sub.Name, sub.NameEn, sub.NameCn, sub.Price1)
				if err != nil {
					return err
				}

				item_id, _ := rs_sub.LastInsertId()
				sub.Id = item_id
				fmt.Println("Item_sub Id : ", item_id)
			}else{
				fmt.Println("sub name = ",sub.Active,sub.Name,i.Id)
				sql_sub := `UPDATE price_sub set price = ?, active = ? where item_id = ? and name = ?`
				fmt.Println("sql_sub =",sql_sub,sub.Price1, sub.Active, i.Id, sub.Name)
				_, err := db.Exec(sql_sub, sub.Price1, sub.Active, i.Id, sub.Name)
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}

// func (i *Item) PrintTest(db *sqlx.DB) error {

// 	config := new(Config)
// 	config = GetConfig(db)

// 	//myPassword := genMikrotikPassword(config)
// 	//fmt.Println("password =",myPassword)

// 	fmt.Println(config.Printer3Port)

// 	f, err := net.Dial("tcp", config.Printer3Port)

// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	w := bufio.NewWriter(f)
// 	p := escpos.New(f)

// 	p.Init()
// 	p.SetSmooth(1)
// 	p.SetFontSize(2, 3)
// 	p.SetFont("A")
// 	p.Write("test ")
// 	p.SetFont("B")
// 	p.Write("test2 ")
// 	p.SetFont("C")
// 	p.Write("test3 ")
// 	p.Formfeed()

// 	p.SetFont("B")
// 	p.SetFontSize(1, 1)

// 	p.SetEmphasize(1)
// 	p.Write("halle")
// 	p.Formfeed()

// 	p.SetUnderline(1)
// 	p.SetFontSize(4, 4)
// 	p.Write("halle")

// 	p.SetReverse(1)
// 	p.SetFontSize(2, 4)
// 	p.Write("halle")
// 	p.Formfeed()

// 	p.SetFont("C")
// 	p.SetFontSize(8, 8)
// 	p.Write("halle")
// 	p.FormfeedN(5)

// 	p.Cut()
// 	p.End()

// 	return nil
// }
