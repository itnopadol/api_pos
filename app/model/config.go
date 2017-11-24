package model

import (
	"github.com/jmoiron/sqlx"
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
)

type Config struct {
	Id           int64      `json:"id" db:"Id"`
	CompanyName  string     `json:"company_name" db:"company_name"`
	Address      string     `json:"address" db:"address"`
	Telephone    string     `json:"telephone" db:"telephone"`
	Fax          string     `json:"fax" db:"fax"`
	LineId       string     `json:"line_id" db:"line_id"`
	Facebook     string     `json:"facebook" db:"facebook"`
	TaxId        string     `json:"tax_id" db:"tax_id"`
	TaxRate      int        `json:"tax_rate" db:"tax_rate"`
	Printer1Port string     `json:"printer1_port" db:"printer1_port"`
	Printer2Port string     `json:"printer2_port" db:"printer2_port"`
	Printer3Port string     `json:"printer3_port" db:"printer3_port"`
	Printer4Port string     `json:"printer4_port" db:"printer4_port"`
	LinkMikrotik string     `json:"link_mikrotik" db:"link_mikrotik"`
	CreatedBy    string     `json:"created_by" db:"created_by"`
	Created      *time.Time `json:"created" db:"created"`
	EditedBy     string     `json:"edited_by" db:"edited_by"`
	Edited       *time.Time `json:"edited" db:"edited"`
}


func (c *Config) Save(db *sqlx.DB) error {

	var checkCount int
	sqlCheckExist := `select count(id) as vCount from config`
	err := db.Get(&checkCount, sqlCheckExist)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if (checkCount == 0) {
		sql := `INSERT INTO config(company_name,address,telephone,fax,line_id,facebook,tax_id,tax_rate,printer1_port,printer2_port,printer3_port,printer4_port,link_mikrotik,created_by,created) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
		rs, err := db.Exec(sql, c.CompanyName, c.Address, c.Telephone, c.Fax, c.LineId, c.Facebook, c.TaxId, c.TaxRate, c.Printer1Port, c.Printer2Port, c.Printer3Port, c.Printer4Port, c.LinkMikrotik, c.CreatedBy)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		Id, _ := rs.LastInsertId()
		c.Id = Id
	}

	return nil
}

func (c *Config) Update(db *sqlx.DB) error {

	var checkCount int
	sqlCheckExist := `select count(id) as vCount from config`
	err := db.Get(&checkCount, sqlCheckExist)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if (checkCount != 0) {
		sql := `Update config set company_name=?, address=?, telephone=?, fax=?, line_id=?, facebook=?, tax_id=?, tax_rate=?, printer1_port=?, printer2_port=?, printer3_port=?, printer4_port=?, link_mikrotik=?, edited_by=?, edited=CURRENT_TIMESTAMP()`
		_, err := db.Exec(sql, c.CompanyName, c.Address, c.Telephone, c.Fax, c.LineId, c.Facebook, c.TaxId, c.TaxRate, c.Printer1Port, c.Printer2Port, c.Printer3Port, c.Printer4Port, c.LinkMikrotik, c.CreatedBy)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func (c *Config) Search(db *sqlx.DB) error {

	sql := `select ifnull(company_name,'') as company_name,ifnull(address,'') as address,ifnull(telephone,'') as telephone,ifnull(fax,'') as fax,ifnull(line_id,'') as line_id,ifnull(facebook,'') as facebook,ifnull(tax_id,'') as tax_id,ifnull(tax_rate,0) as tax_rate,ifnull(printer1_port,'') as printer1_port,ifnull(printer2_port,'') as printer2_port,ifnull(printer3_port,'') as printer3_port,ifnull(printer4_port,'') as printer4_port,ifnull(link_mikrotik,'') as link_mikrotik from config`
	err := db.Get(c, sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (c *Config) GenWifiPassword(db *sqlx.DB) error {
	sql := `select ifnull(company_name,'') as company_name,ifnull(address,'') as address,ifnull(tax_id,'') as tax_id,ifnull(tax_rate,0) as tax_rate,ifnull(printer1_port,'') as printer1_port,ifnull(printer2_port,'') as printer2_port,ifnull(printer3_port,'') as printer3_port,ifnull(printer4_port,'') as printer4_port, ifnull(link_mikrotik,'') as link_mikrotik from config`
	fmt.Println("Config = ", sql)
	err := db.Get(c, sql)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Link Wifi : ",c.LinkMikrotik)

	wifi_link := "http://hapos.dyndns.org:9003/wifi/genuser.php"
	res, err := http.Get(wifi_link)
	if err != nil {
		return err
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	password := string(robots)

	fmt.Println("robots wifi = ", password)

	return nil
}
