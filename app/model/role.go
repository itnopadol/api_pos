package model

import "time"

type Role struct {
	Id int `json:"id" db:"id"`
	RoleCode	string `json:"role_code" db:"role_code"`
	RoleName string `json:"role_name" db:"role_name"`
	Active int `json:"active" db:"active"`
	CreatedBy string `json:"created_by" db:"created_by"`
	Created time.Time `json:"created" db:"created"`
	EditedBy string `json:"edited_by" db:"edited_by"`
	Edited time.Time `json:"edited" db:"edited"`
}