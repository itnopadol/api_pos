package model

import (
	"time"
)

type User struct {
	Id int `json:"id" db:"id"`
	UserCode string `json:"user_code" db:"user_code"`
	UserName string `json:"user_name" db:"user_name"`
	Password string `json:"password" db:"password"`
	ConfirmPassword string `json:"confirm_password" db:"confirm_password"`
	RoleId int `json:"role_id" db:"role_id"`
	Active int `json:"active" db:"active"`
	CreatedBy string `json:"created_by" db:"created_by"`
	Created time.Time `json:"created" db:"created"`
	EditedBy string `json:"edited_by" db:"edited_by"`
	Edited time.Time `json:"edited" db:"edited"`
}
