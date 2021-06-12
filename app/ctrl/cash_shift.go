package ctrl

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
)

func SaveShift(c *gin.Context) {
	fmt.Println("Call POST SaveShift")
	c.Keys = headerKeys

	newShift := &model.Shift{}
	err := c.BindJSON(newShift)
	if err != nil {
		fmt.Println(err.Error())
	}
	ch, err := newShift.SaveShift(dbc)
	if err != nil {
		fmt.Println(err.Error())
	}
	rs := resp.Response{}
	fmt.Println("ch =", ch)
	if ch != "" || err != nil {
		rs.Status = "error"
		rs.Message = "No Content :" + ch
		c.JSON(http.StatusNotFound, rs)
	} else if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = ch
		c.JSON(http.StatusOK, rs)
	}
}

func UpdateShift(c *gin.Context) {
	fmt.Println("Call PUT UpdateShift")
	c.Keys = headerKeys

	newShift := &model.Shift{}
	err := c.BindJSON(newShift)
	if err != nil {
		fmt.Println(err.Error())
	}

	ch,err := newShift.UpdateShift(dbc)
	if err != nil {
		fmt.Println(err.Error())
	}
	rs := resp.Response{}
	fmt.Println("ch =", ch)
	if ch != "" || err != nil {
		rs.Status = "error"
		rs.Message = "No Content :" + ch
		c.JSON(http.StatusNotFound, rs)
	} else if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = ch
		c.JSON(http.StatusOK, rs)
	}
	// if err != nil {
	// 	rs.Status = "error"
	// 	rs.Message = "No Content :" + err.Error()
	// 	c.JSON(http.StatusNotFound, rs)
	// } else {
	// 	rs.Status = "success"
	// 	rs.Data = ch
	// 	c.JSON(http.StatusOK, rs)
	// }
}

func ClosedShift(c *gin.Context) {
	fmt.Println("Call PUT ClosedShift")
	c.Keys = headerKeys

	newShift := &model.Shift{}
	err := c.BindJSON(newShift)
	if err != nil {
		fmt.Println(err.Error())
	}

	ch,err := newShift.ClosedShift(dbc)
	if err != nil {
		fmt.Println(err.Error())
	}
	rs := resp.Response{}
	fmt.Println("ch =", ch)
	if ch != "" || err != nil {
		rs.Status = "error"
		rs.Message = "No Content :" + ch
		c.JSON(http.StatusNotFound, rs)
	} else if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = ch
		c.JSON(http.StatusOK, rs)
	}
	// rs := resp.Response{}
	// if err != nil {
	// 	rs.Status = "error"
	// 	rs.Message = "No Content :" + err.Error()
	// 	c.JSON(http.StatusNotFound, rs)
	// } else {
	// 	rs.Status = "success"
	// 	rs.Data = ch
	// 	c.JSON(http.StatusOK, rs)
	// }
}

func ShiftDetails(c *gin.Context) {
	fmt.Println("Call GET ShiftDetails")
	c.Keys = headerKeys

	host_code := c.Request.URL.Query().Get("host_code")
	doc_date := c.Request.URL.Query().Get("doc_date")

	//host_code := hostcode//strconv.ParseInt(strId, 10, 64)

	ch := new(model.Shift)

	err := ch.ShiftDetails(dbc, host_code, doc_date)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = ch
		c.JSON(http.StatusOK, rs)
	}
}

func ShiftLastID(c *gin.Context) {
	fmt.Println("Call GET ShiftLastID")
	c.Keys = headerKeys

	host_code := c.Request.URL.Query().Get("host_code")

	//host_code := hostcode//strconv.ParseInt(strId, 10, 64)

	ch := new(model.Shift)

	err := ch.ShiftLastID(dbc, host_code)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = ch
		c.JSON(http.StatusOK, rs)
	}
}

func SearchShiftByKeyword(c *gin.Context) {
	fmt.Println("Call GET ShiftDetails")
	c.Keys = headerKeys

	hostid := c.Param("host_id")
	host_id := hostid //strconv.ParseInt(strId, 10, 64)

	docdate := c.Param("doc_date")
	doc_date := docdate //strconv.ParseInt(strId, 10, 64)

	ch := new(model.Shift)

	shifts, err := ch.SearchShiftByKeyword(dbc, host_id, doc_date)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = shifts
		c.JSON(http.StatusOK, rs)
	}
}

func PrintSendDailyTotal(c *gin.Context) {
	log.Println("call Get SendDaily")
	c.Keys = headerKeys

	host_code := c.Request.URL.Query().Get("host_code")
	doc_date := c.Request.URL.Query().Get("doc_date")

	NewShift := new(model.Shift)
	shifts, err := NewShift.PrintSendDailyTotal(dbc, host_code, doc_date)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = shifts
		c.JSON(http.StatusOK, rs)
	}

}
