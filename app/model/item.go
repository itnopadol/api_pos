package model

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type Item struct {
	Id      int
	Code 	string `json:"code" db:"code"`
	Name    string   `json:"name" db:"name"`
	ShortName	string `json:"short_name" db:"short_name"`
	NameEn  string   `json:"name_en,omitempty" db:"name_en"`
	NameCn  string   `json:"name_cn,omitempty" db:"name_cn"`
	Price   float32  `json:"price" db:"price"`
	Unit    string   `json:"unit"`
	UnitEn  string   `json:"unit_en,omitempty" db:"unit_en"`
	UnitCn  string   `json:"unit_cn,omitempty" db:"unit_cn"`
	MenuId  uint64   `json:"menu_id,omitempty" db:"menu_id"`
	MenuSeq int      `json:"menu_seq,omitempty" db:"menu_seq"`
	Image   string   `json:"image" db:"image"`
	IsKitchen int `json:"is_kitchen" db:"is_kitchen"`
	Active  int 	 `json:"active" db:"active"`
	Prices []*PricesSub `json:"prices_sub"`
}

type PricesSub struct {
	Id     int     `json:"id"`
	ItemId int64   `json:"-" db:"item_id"`
	Name   string  `json:"name" db:"name"`
	NameEn string  `json:"name_en" db:"name_en"`
	NameCn string  `json:"name_cn" db:"name_cn"`
	Price1  float32 `json:"price" db:"price"`
	Active  int `json:"active" db:"active"`
}

func (i *Item) Get(db *sqlx.DB, id int64) (err error) {
	sql := `SELECT * FROM item WHERE active = 1 and id = ?`
	fmt.Println("Item = ",sql, id)
	err = db.Get(i, sql, id)

	//err = db.QueryRowx(sql, id).StructScan(i)
	if err != nil {
		return err
	}

	fmt.Println("Item Name = ",i.Name )

	vLen := len(i.Name)

	fmt.Println("Lenght = ", vLen/3)
	// ดึงข้อมูลราคาทั้งหมดของสินค้ารายการนี้
	sizes := []*PricesSub{}
	sql = `SELECT * FROM price_sub WHERE active = 1 and item_id = ?`
	fmt.Println("Price = ", sql)
	fmt.Println("Lenght = ", vLen)
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
		fmt.Println("Lang1 = ",l.Name)
		items := []*Item{}
		switch l.Id {
		case 1:
			sql = `SELECT id,code,short_name,name, unit, menu_seq, image, price, active, is_kitchen FROM item WHERE active = 1 and  menu_id = ?`
		case 2:
			sql = `SELECT id,code,short_name,name_en as name, unit_en as unit, menu_seq, image, price, active, is_kitchen FROM item WHERE active = 1 and  menu_id = ?`
		case 3:
			sql = `SELECT id,code,short_name,name_cn as name, unit_cn as unit, menu_seq, image, price, active, is_kitchen FROM item WHERE active = 1 and  menu_id = ?`
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
