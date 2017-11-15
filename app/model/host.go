package model

import (
	"github.com/jmoiron/sqlx"
)

type Host struct {
	Id int `json:"id" db:"id"`
	HostCode string `json:"host_code" db:"host_code"`
	HostName string `json:"host_name" db:"host_name"`
	PrinterPort string `json:"printer_port" db:"printer_port"`
	Status int `json:"status" db:"status"`
	Active int `json:"active" db:"active"`
	LogoImageId           int
	LogoImageWidth        int
}

func (h *Host)SearchHost(db *sqlx.DB)(hosts []*Host, err error) {
	sql := `select host_code,ifnull(host_name,'') as host_name,ifnull(printer_port,'') as printer_port,status,active from host  order by id`
	err = db.Select(&hosts, sql)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}


func (h *Host)GetHostPrinter(db *sqlx.DB,hostcode string)(host *Host, err error) {
	sql := `select host_code,ifnull(host_name,'') as host_name,ifnull(printer_port,'') as printer_port,status,active from host where host_code= ? order by id`
	err = db.Select(&host, sql,hostcode)
	if err != nil {
		return nil, err
	}
	return host, nil
}

