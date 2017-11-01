package ctrl

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/itnopadol/hapos_api/app/model"
	"github.com/itnopadol/bc_api/bc_api/bean/resp"
	"net/http"
)

func SaveShift(c *gin.Context){
	fmt.Println("Call POST SaveShift")
	c.Keys = headerKeys

	newShift := &model.Shift{}
	err := c.BindJSON(newShift)
	if err != nil {
		fmt.Println(err.Error())
	}
	s := newShift.SaveShift(dbc)

	rs := Resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = s
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

	s := newShift.UpdateShift(dbc)

	rs := Resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = s
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

	s := newShift.ClosedShift(dbc)

	rs := Resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = s
		c.JSON(http.StatusOK, rs)
	}
}

func ShiftDetails(c *gin.Context){
	fmt.Println("Call GET ShiftDetails")
	c.Keys = headerKeys

	hostid := c.Param("host_id")
	host_id := hostid//strconv.ParseInt(strId, 10, 64)

	s := new(model.Shift)

	err := s.ShiftDetails(dbc,host_id)
	rs := Resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = s
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

	s := new(model.Shift)

	shifts, err := s.SearchShiftByKeyword(dbc,host_id,doc_date)
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