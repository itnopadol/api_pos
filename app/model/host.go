package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Host struct {
	Id             int64  `json:"id" db:"id"`
	HostCode       string `json:"host_code" db:"host_code"`
	HostName       string `json:"host_name" db:"host_name"`
	PrinterPort    string `json:"printer_port" db:"printer_port"`
	Status         int    `json:"status" db:"status"`
	Active         int    `json:"active" db:"active"`
	BranchID       int    `json:"branch_id" db:"branch_id"`
	BranchCode     string `json:"branch_code" db:"branch_code"`
	HeadDocSale    string `json:"head_doc_sale" db:"head_doc_sale"`
	LogoImageId    int
	LogoImageWidth int
}

func (h *Host) Save(db *sqlx.DB) error {
	var vCountHost int
	fmt.Println("Save Host")

	sqlCheck := `select count(id) as vCount from host where host_code = ?`
	err := db.Get(&vCountHost, sqlCheck, h.HostCode)
	if err != nil {
		return err
	}

	fmt.Println("Count Host", vCountHost)
	if vCountHost == 0 {
		sql := `Insert Into host(host_code, host_name, status, printer_port, active) Values(?,?,0,?,1)`
		rs, err := db.Exec(sql, h.HostCode, h.HostName, h.PrinterPort)
		if err != nil {
			return err
		}

		id, _ := rs.LastInsertId()
		h.Id = id
		fmt.Println("Insert ID : ", id)
	}

	return nil
}

func (h *Host) Update(db *sqlx.DB) error {
	var vCountHost int

	sqlCheck := `select count(id) as vCount from host where host_code = ?`
	err := db.Get(&vCountHost, sqlCheck, h.HostCode)
	if err != nil {
		return err
	}

	if vCountHost != 0 {
		sql := `Update host set host_code=?, host_name=?, printer_port=?, active=? where id =?`
		_, err := db.Exec(sql, h.HostCode, h.HostName, h.PrinterPort, h.Active, h.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Host) SearchHost(db *sqlx.DB) (hosts []*Host, err error) {
	//sql := `select host_code,ifnull(host_name,'') as host_name,ifnull(printer_port,'') as printer_port,status,active from host  order by id`
	sql := `select host_code,ifnull(host_name,'') as host_name,ifnull(printer_port,'') as printer_port,status,active,ifnull(b.branch_id,1) as branch_id,ifnull(b.branch_code,'') as branch_code, ifnull(b.head_doc_sale,'') as head_doc_sale from host a left join config b  on b.id = 1  order by a.id`
	err = db.Select(&hosts, sql)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func (h *Host) GetHostPrinter(db *sqlx.DB, hostcode string) (host *Host, err error) {
	//sql := `select host_code,ifnull(host_name,'') as host_name,ifnull(printer_port,'') as printer_port,status,active from host where host_code= ? order by id`
	sql := `select host_code,ifnull(host_name,'') as host_name,ifnull(printer_port,'') as printer_port,status,active,ifnull(b.branch_id,1) as branch_id,ifnull(b.branch_code,'') as branch_code, ifnull(b.head_doc_sale,'') as head_doc_sale from host a left join config b  on b.id = 1  order by a.id`
	err = db.Select(&host, sql, hostcode)
	if err != nil {
		return nil, err
	}
	return host, nil
}
