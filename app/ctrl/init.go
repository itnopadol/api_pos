package ctrl

import (
	"github.com/jmoiron/sqlx"
	"fmt"
)

var headerKeys = make(map[string]interface{})

func setHeader() {

	headerKeys = map[string]interface{}{
		"Server":                       "smart_pump_invoice",
		"Host":                         "nopadol.net:6000",
		"Content_Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE",
		"Access-Control-Allow-Headers": "Origin, Content-Type, X-Auth-Token",
	}
}

var dbc *sqlx.DB

func init() {
	dbc = ConnectMySql()
}

func ConnectMySql() (mydb *sqlx.DB) {
	//dsn := "sa:[ibdkifu@tcp(192.168.1.250:3306)/" + "pos" + "?parseTime=true&charset=utf8&loc=Local" //ใช้เวลาอัพขึ้น server แม่ริม
	//dsn := "sa:[ibdkifu@tcp(hapos.dyndns.org:9010)/"+ "pos" +"?parseTime=true&charset=utf8&loc=Local"//ลิงค์นอก เรียกข้อมูลจริง แม่ริม
	dsn := "root:pordeeproject88@tcp(68.183.191.228:3306)/"+ "pos" +"?parseTime=true&charset=utf8&loc=Local"
	mydb = sqlx.MustConnect("mysql", dsn)
	if (mydb.Ping() != nil) {
		fmt.Println("Error")
	}
	fmt.Println("mysql = ", mydb.DriverName(), "dsn = ", dsn)
	return mydb
}
