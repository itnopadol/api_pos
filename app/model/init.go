package model

import (
	"math"
	"github.com/jmoiron/sqlx"
	"fmt"
	"strconv"
	"time"
	"strings"
	"bytes"
	"net/http"
	"io/ioutil"
)

var (
	H *Host
)

func init() {

}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func GenDocno(db *sqlx.DB, host_code string) (doc_no string) {
	var last_number1 int
	var last_number string
	var snumber string
	var intyear int
	var vHeader string
	var vyear string

	var intmonth int
	var intmonth1 int
	var vmonth string
	var vmonth1 string
	var lenmonth int

	var intday int
	var intday1 int
	var vday string
	var vday1 string
	var lenday int

	last_number1, _ = GetLastDocNo(db, host_code)
	last_number = strconv.Itoa(last_number1)
	fmt.Println("Last No = ", last_number)
	if (time.Now().Year() >= 2560) {
		intyear = time.Now().Year()
	} else {
		intyear = time.Now().Year() + 543
	}

	vyear = strconv.Itoa(intyear)
	vyear1 := vyear[2:len(vyear)]

	fmt.Println("year = ", vyear1)

	intmonth = int(time.Now().Month())
	intmonth1 = int(intmonth)
	vmonth = strconv.Itoa(intmonth1)

	fmt.Println("month =", vmonth)

	lenmonth = len(vmonth)

	if (lenmonth == 1) {
		vmonth1 = "0" + vmonth
	} else {
		vmonth1 = vmonth
	}

	intday = int(time.Now().Day())
	intday1 = int(intday)
	vday = strconv.Itoa(intday1)

	fmt.Println("day =", vday)

	lenday = len(vday)

	if (lenday == 1) {
		vday1 = "0" + vday
	} else {
		vday1 = vday
	}

	if (len(string(last_number)) == 1) {
		fmt.Println("Last_number =", last_number)
		snumber = "000" + last_number
	}
	if (len(string(last_number)) == 2) {
		snumber = "00" + last_number
	}
	if (len(string(last_number)) == 3) {
		snumber = "0" + last_number
	}
	if (len(string(last_number)) == 4) {
		snumber = last_number
	}

	fmt.Println(snumber)
	fmt.Println(vHeader)

	doc_no = host_code + vyear1 + vmonth1 + vday1 + "-" + snumber
	fmt.Println(snumber)
	fmt.Println(vHeader)

	fmt.Println("NewDocNo = ", doc_no)

	return doc_no
}

func GetLastDocNo(db *sqlx.DB, host_code string) (last_no int, err error) {
	sql := `select CONVERT(right(ifnull(max(doc_no),0),4),UNSIGNED INTEGER)+1 as maxno from sale where host_code = ? and year(doc_date) = year(CURDATE()) and month(doc_date) = month(CURDATE()) and day(doc_date) = day(CURDATE())`

	fmt.Println("Query = ", sql)
	err = db.Get(&last_no, sql, host_code)
	if err != nil {
		fmt.Println(err)
		return 1, err
	}

	fmt.Println("Last No = ", last_no)
	return last_no, nil
}

func LastQueId(db *sqlx.DB) (que_id int) {
	sql := `select ifnull(max(que_id),0)+1 maxno from sale where year(doc_date) = year(CURDATE()) and month(doc_date) = month(CURDATE()) and day(doc_date) = day(CURDATE())`
	fmt.Println("Query = ", sql)
	err := db.Get(&que_id, sql)
	if err != nil {
		fmt.Println(err)
		return que_id
	}

	fmt.Println("Last No = ", que_id)
	return que_id
}

func CommaFloat(v float64) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}

	comma := []byte{','}

	parts := strings.Split(strconv.FormatFloat(v, 'f', -1, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos: pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

func genMikrotikPassword(c *Config) (password string) {
	res, err := http.Get(c.LinkMikrotik)
	if err != nil {
		return ""
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	password = string(robots)

	fmt.Println("robots wifi = ", password)

	return password
}

func GetConfig(db *sqlx.DB) (config *Config) {
	cf := new(Config)
	sql := `select ifnull(company_name,'') as company_name,ifnull(address,'') as address,ifnull(tax_id,'') as tax_id,ifnull(tax_rate,0) as tax_rate,ifnull(printer1_port,'') as printer1_port,ifnull(printer2_port,'') as printer2_port,ifnull(printer3_port,'') as printer3_port,ifnull(printer4_port,'') as printer4_port, ifnull(link_mikrotik,'') as link_mikrotik from config`
	fmt.Println("Config = ", sql)
	err := db.Get(cf, sql)
	if err != nil {
		fmt.Println(err.Error())
	}

	config = cf
	return config
}

func GetHostPrinter(db *sqlx.DB, host_code string) (host *Host) {
	h := new(Host)
	sql := `select ifnull(host_code,'') as host_code,ifnull(host_name,'') as host_name,ifnull(printer_port,'') as printer_port,ifnull(status,0) as status,active from host where host_code = ? and active = 1`
	fmt.Println("Host = ", sql)
	err := db.Get(h, sql, host_code)
	if err != nil {
		fmt.Println(err.Error())
	}

	host = h
	return host
}
func Commaf(v float64) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}

	comma := []byte{','}

	parts := strings.Split(strconv.FormatFloat(v, 'f', 2, 64), ".")
	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos: pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return buf.String()
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}