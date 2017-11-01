package model

type Config struct {
	Id int `json:"id" db:"Id"`
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
