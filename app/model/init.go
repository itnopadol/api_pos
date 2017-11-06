package model

import (
	"math"
	"github.com/jmoiron/sqlx"
	"fmt"
	"strconv"
	"time"
)


var (
	H        *Host
)

func init(){

}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num * output)) / output
}

func GenDocno(db *sqlx.DB,host_code string)(doc_no string) {
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
	if(time.Now().Year()>=2560){
		intyear = time.Now().Year()
	}else{
		intyear = time.Now().Year()+543
	}

	vyear = strconv.Itoa(intyear)
	vyear1 := vyear[2:len(vyear)]

	fmt.Println("year = ",vyear1)

	intmonth = int(time.Now().Month())
	intmonth1 = int(intmonth)
	vmonth = strconv.Itoa(intmonth1)

	fmt.Println("month =",vmonth)

	lenmonth = len(vmonth)

	if(lenmonth==1){
		vmonth1 = "0"+vmonth
	}else {
		vmonth1 = vmonth
	}

	intday = int(time.Now().Day())
	intday1 = int(intday)
	vday = strconv.Itoa(intday1)

	fmt.Println("day =",vday)

	lenday = len(vday)

	if(lenday==1){
		vday1 = "0"+vday
	}else {
		vday1 = vday
	}

	if(len(string(last_number))==1){
		fmt.Println("Last_number =",last_number)
		snumber = "000"+last_number
	}
	if(len(string(last_number))==2){
		snumber = "00"+last_number
	}
	if(len(string(last_number))==3){
		snumber = "0"+last_number
	}
	if(len(string(last_number))==4) {
		snumber = last_number
	}

	fmt.Println(snumber)
	fmt.Println(vHeader)

	doc_no = host_code+vyear1+vmonth1+vday1+"-"+snumber
	fmt.Println(snumber)
	fmt.Println(vHeader)

	fmt.Println("NewDocNo = ",doc_no)

	return doc_no
}

func GetLastDocNo(db *sqlx.DB, host_code string) (last_no int, err error){
	sql := `select CONVERT(right(ifnull(max(doc_no),0),4),UNSIGNED INTEGER)+1 as maxno from sale where host_code = ? and year(doc_date) = year(CURDATE()) and month(doc_date) = month(CURDATE()) and day(doc_date) = day(CURDATE())`

	fmt.Println("Query = ",sql)
	err = db.Get(&last_no,sql,host_code)
	if err != nil {
		fmt.Println(err)
		return 1, err
	}

	fmt.Println("Last No = ",last_no)
	return last_no, nil
}


func LastQueId(db *sqlx.DB) (que_id int){
	sql := `select ifnull(max(que_id),0)+1 maxno from sale where year(doc_date) = year(CURDATE()) and month(doc_date) = month(CURDATE()) and day(doc_date) = day(CURDATE())`
	fmt.Println("Query = ",sql)
	err := db.Get(&que_id,sql)
	if err != nil {
		fmt.Println(err)
		return  que_id
	}

	fmt.Println("Last No = ",que_id)
	return que_id
}