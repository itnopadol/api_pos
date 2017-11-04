package ctrl

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
	"net/http"
	"log"
	"github.com/itnopadol/bc_api/bc_api/bean/resp"
)

func SaveShift(c *gin.Context){
	fmt.Println("Call POST SaveShift")
	c.Keys = headerKeys

	newShift := &model.Shift{}
	err := c.BindJSON(newShift)
	if err != nil {
		fmt.Println(err.Error())
	}
	ch := newShift.SaveShift(dbc)

	rs := Resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = ch
		c.JSON(http.StatusOK, rs)
	}
}

func UpdateShift(c *gin.Context){
	fmt.Println("Call PUT UpdateShift")
	c.Keys = headerKeys

	newShift := &model.Shift{}
	err := c.BindJSON(newShift)
	if err != nil {
		fmt.Println(err.Error())
	}

	ch := newShift.UpdateShift(dbc)

	rs := Resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = ch
		c.JSON(http.StatusOK, rs)
	}
}

func ClosedShift(c *gin.Context){
	fmt.Println("Call PUT ClosedShift")
	c.Keys = headerKeys

	newShift := &model.Shift{}
	err := c.BindJSON(newShift)
	if err != nil {
		fmt.Println(err.Error())
	}

	ch := newShift.ClosedShift(dbc)

	rs := Resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = ch
		c.JSON(http.StatusOK, rs)
	}
}

func ShiftDetails(c *gin.Context){
	fmt.Println("Call GET ShiftDetails")
	c.Keys = headerKeys

	host_code := c.Request.URL.Query().Get("host_code")
	doc_date := c.Request.URL.Query().Get("doc_date")

	//host_code := hostcode//strconv.ParseInt(strId, 10, 64)

	ch := new(model.Shift)

	 err := ch.ShiftDetails(dbc,host_code,doc_date)
	rs := Resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = ch
		c.JSON(http.StatusOK, rs)
	}
}

func SearchShiftByKeyword(c *gin.Context){
	fmt.Println("Call GET ShiftDetails")
	c.Keys = headerKeys

	hostid := c.Param("host_id")
	host_id := hostid//strconv.ParseInt(strId, 10, 64)

	docdate := c.Param("doc_date")
	doc_date := docdate//strconv.ParseInt(strId, 10, 64)

	ch := new(model.Shift)

	shifts, err := ch.SearchShiftByKeyword(dbc,host_id,doc_date)
	rs := Resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = shifts
		c.JSON(http.StatusOK, rs)
	}
}


func PrintSendDailyTotal(c *gin.Context){
	log.Println("call Get SearchSales")
	c.Keys = headerKeys

	doc_date := c.Request.URL.Query().Get("doc_date")

	NewSale := new(model.Shift)
	sales, err := NewSale.PrintSendDailyTotal(dbc,doc_date)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = sales
		c.JSON(http.StatusOK, rs)
	}

}