package model

import (
	"github.com/jmoiron/sqlx"
)

type Host struct {
	Id int `json:"id" db:"id"`
	HostCode string `json:"host_code" db:"host_code"`
	HostName string `json:"host_name" db:"host_name"`
	Status int `json:"status" db:"status"`
	Active int `json:"active" db:"active"`
	LogoImageId           int
	LogoImageWidth        int
}

func (h *Host)SearchHost(db *sqlx.DB)(hosts []*Host, err error) {
	sql := `select host_code,host_name,status,active from host where active = 1 order by id`
	err = db.Select(&hosts, sql)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}


