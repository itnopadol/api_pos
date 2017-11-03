package model

import (
	"github.com/jmoiron/sqlx"
	"fmt"
)

type Config struct {
	Id int64 `json:"id" db:"Id"`
	CompanyName string `json:"company_name" db:"company_name"`
	Address string `json:"address" db:"address"`
	Telephone string `json:"telephone" db:"telephone"`
	Fax string `json:"fax" db:"fax"`
	LineId string `json:"line_id" db:"line_id"`
	Facebook string `json:"facebook" db:"facebook"`
	TaxId string `json:"tax_id" db:"tax_id"`
	TaxRate int `json:"tax_rate" db:"tax_rate"`
	Printer1Port string `json:"printer1_port" db:"printer1_port"`
	Printer2Port string `json:"printer2_port" db:"printer2_port"`
	Printer3Port string `json:"printer3_port" db:"printer3_port"`
}

func (c *Config)Save(db *sqlx.DB)error{
	sql := `INSERT INTO config(company_name,address,telephone,fax,line_id,facebook,tax_id,tax_rate,printer1_port,printer2_port,printer3_port) VALUES(?,?,?,?,?,?,?,?,?,?,?)`
	rs, err := db.Exec(sql, c.CompanyName, c.Address, c.Telephone, c.Fax, c.LineId, c.Facebook, c.TaxId, c.TaxRate, c.Printer1Port, c.Printer2Port, c.Printer3Port)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Id , _:= rs.LastInsertId()
	c.Id = Id

	return nil
}

func (c *Config)Search(db *sqlx.DB)error{

	sql := `select ifnull(company_name,'') as company_name,ifnull(address,'') as address,ifnull(telephone,'') as telephone,ifnull(fax,'') as fax,ifnull(line_id,'') as line_id,ifnull(facebook,'') as facebook,ifnull(tax_id,'') as tax_id,ifnull(tax_rate,0) as tax_rate,ifnull(printer1_port,'') as printer1_port,ifnull(printer2_port,'') as printer2_port,ifnull(printer3_port,'') as printer3_port from config`
	err := db.Get(&c, sql)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

