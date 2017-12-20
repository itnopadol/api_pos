package model

import (
	"time"
	"github.com/jmoiron/sqlx"
	"fmt"
)

type User struct {
	Id              int64     `json:"id" db:"id"`
	UserCode        string    `json:"user_code" db:"user_code"`
	UserName        string    `json:"user_name" db:"user_name"`
	Password        string    `json:"password" db:"password"`
	ConfirmPassword string    `json:"confirm_password" db:"confirm_password"`
	RoleId          int       `json:"role_id" db:"role_id"`
	Active          int       `json:"active" db:"active"`
	CreatedBy       string    `json:"created_by" db:"created_by"`
	Created         time.Time `json:"created" db:"created"`
	EditedBy        string    `json:"edited_by" db:"edited_by"`
	Edited          time.Time `json:"edited" db:"edited"`
}

func (u *User) LogIn(db *sqlx.DB, user_code string, password string) error {
	sql := `select user_code,user_name,role_id,active from user where user_code = ? and password = ? and active = 1`
	err := db.Get(u, sql, user_code, password)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) SearchUser(db *sqlx.DB, keyword string) (users []*User, err error) {

	if (keyword == "") {
		sql := `select id,user_code,user_name,password,confirm_password,role_id,active from user where active = 1 order by id`
		err = db.Select(&users, sql)
	} else {
		sql := `select id,user_code,user_name,password,confirm_password,role_id,active from user where active = 1 and user_code like CONCAT("%",?,"%") or user_name like CONCAT("%",?,"%")  order by id`
		err = db.Select(&users, sql, keyword, keyword)
	}
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) ListUser(db *sqlx.DB, keyword string) (users []*User, err error) {
	//var sql string

	fmt.Println("keyword = ", keyword)
	if (keyword == "") {
		sql := `select id,user_code,user_name,password,confirm_password,role_id,active from user where active = 1 order by id`
		err = db.Select(&users, sql)
	} else {
		sql := `select id,user_code,user_name,password,confirm_password,role_id,active from user where active = 1 and (user_code like CONCAT("%",?,"%") or user_name like CONCAT("%",?,"%")) order by id`
		err = db.Select(&users, sql, keyword, keyword)
	}

	//fmt.Println("keyword", sql)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) Save(db *sqlx.DB) error {

	var checkCount int
	sqlCheckExist := `select count(id) as vCount from user where user_code = ?`
	err := db.Get(&checkCount, sqlCheckExist, u.UserCode)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if (checkCount == 0) {
		sql := `INSERT INTO user(user_code,user_name,password,confirm_password,role_id,active,created_by,created) VALUES(?,?,?,?,?,?,?,CURRENT_TIMESTAMP())`
		rs, err := db.Exec(sql, u.UserCode, u.UserName, u.Password, u.ConfirmPassword, u.RoleId, 1, u.CreatedBy)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		Id, _ := rs.LastInsertId()
		u.Id = Id
	}

	return nil
}

func (u *User) Update(db *sqlx.DB) error {

	var checkCount int
	sqlCheckExist := `select count(id) as vCount from user where id = ? and user_code = ?`
	err := db.Get(&checkCount, sqlCheckExist, u.Id, u.UserCode)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("Count : ", checkCount)
	fmt.Println("ID = ", u.Id)

	if (u.Password == u.ConfirmPassword) {
		if (checkCount != 0) {
			sql := `Update user set user_name=?,password=?,confirm_password=?,role_id=?,active=?,edited_by=?, edited=CURRENT_TIMESTAMP() where id = ? and user_code = ?`
			rs, err := db.Exec(sql, u.UserName, u.Password, u.ConfirmPassword, u.RoleId, u.Active, u.EditedBy, u.Id, u.UserCode)
			if err != nil {
				fmt.Println(err.Error())
				return err
			}

			Id, _ := rs.LastInsertId()
			u.Id = Id
		}
	} else {
		return err
	}

	return nil
}
