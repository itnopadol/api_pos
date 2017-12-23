package model

import (
	"fmt"
	"log"
	"github.com/jmoiron/sqlx"
)

type Lang struct {
	Id     int     `json:"lang_id"`
	Name   string  `json:"lang_name"`
	Menus  []*Menu `json:"menus,omitempty"`
	MenuId int     `json:"menu_id,omitempty"`
	Items  []*Item `json:"items,omitempty"`
}

type Menu struct {
	Id	int64 `json:"id" db:"id"`
	ClientId	int `json:"client_id,omitempty" db:"client_id"`
	Name    string `json:"name" db:"name"`
	NameEn  string `json:"name_en,omitempty" db:"name_en"`
	NameCn  string `json:"name_cn,omitempty" db:"name_cn"`
	Short   string `json:"short,omitempty" db:"short"`
	ShortEn string `json:"short_en,omitempty" db:"short_en"`
	ShortCn string `json:"short_cn,omitempty" db:"short_cn"`
	Image   string `json:"image" db:"image"`
	Link    string `json:"link" db:"link"`
	Status		int `json:"status,omitempty" db:"status"`
	Active 		int `json:"active" db:"active"`
}

var langs = make([]*Lang, 3)

func langInit() {
	langs[0] = &Lang{Id: 1, Name: "Thai Female"}
	langs[1] = &Lang{Id: 2, Name: "UK English Female"}
	langs[2] = &Lang{Id: 3, Name: "Chinese Female"}
}

func (m *Menu) Index(db *sqlx.DB) ([]*Lang, error) {
	var sql string
	langInit()
	for _, l := range langs {
		menus := []*Menu{}
		switch l.Id {
		case 1:
			sql = `SELECT id, name, image, link, active FROM menu where active = 1 `
		case 2:
			sql = `SELECT id, name_en as name, image, link, active FROM menu where active = 1`
		case 3:
			sql = `SELECT id, name_cn as name, image, link, active FROM menu where active = 1`
		}
		//log.Println("case:", l.Id, l.Name)
		err := db.Select(&menus, sql)
		if err != nil {
			return nil, err
		}
		l.Menus = menus
		//log.Println(l)
	}
	log.Println(langs)
	return langs, nil
}

func (m *Menu) Save(db *sqlx.DB) error {
	var checkCount int

	sqlCheck := `select count(*) as vCount from menu where id = ?`
	err := db.Get(&checkCount, sqlCheck, m.Id)
	if err != nil {
		return nil
	}

	if (checkCount == 0) {
		fmt.Println("Menu.Save()")
		sql := `INSERT INTO menu(name, name_en, name_cn, image, link, active) VALUES (?,?,?,?,?,1)`
		fmt.Println("sql = ", sql, m.Name, m.NameEn, m.NameCn, m.Image, m.Link)
		rs, err := db.Exec(sql, m.Name, m.NameEn, m.NameCn, m.Image, m.Link)
		if err != nil {
			fmt.Printf("Error when db.Exec(sql1) %v", err.Error())
			return err
		}
		id, _ := rs.LastInsertId()
		m.Id = id
		fmt.Println("c.Id =", m.Id)
		fmt.Printf("Save data success: category = %+v\n", m)
	}
	return nil
}


func (m *Menu) Update(db *sqlx.DB) error {
	var checkCount int

	sqlCheck := `select count(*) as vCount from menu where id = ?`
	err := db.Get(&checkCount, sqlCheck, m.Id)
	if err != nil {
		return nil
	}

	fmt.Println("checkCount = ",checkCount)
	if (checkCount != 0) {
		sql := `Update menu set name = ?, name_en = ?, name_cn = ?, image = ?, link = ?, active = ? where id = ?`
		_, err := db.Exec(sql, m.Name, m.NameEn, m.NameCn, m.Image, m.Link, m.Active, m.Id)
		if err != nil {
			return nil
		}
	}
	return nil
}