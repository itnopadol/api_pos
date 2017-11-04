package model

import (
)

type Host struct {
	Id int `json:"id" db:"id"`
	HostCode string `json:"host_code" db:"host_code"`
	Status int `json:"status" db:"status"`
	Active int `json:"active" db:"active"`
	LogoImageId           int
	LogoImageWidth        int
}


