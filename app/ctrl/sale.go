package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
	"fmt"
	"net/http"
	"strconv"
	"log"
)

func ShowChangeAmount(c *gin.Context){
	log.Println("call GET ShowChangeAmount")
	c.Keys = headerKeys

	newSale := &model.Sale{}
	err := c.BindJSON(newSale)
	if err != nil {
		fmt.Println(err.Error())
	}
	amount, msg, err := newSale.ShowChangeAmount()
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content: " + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = msg + strconv.FormatFloat(amount, 'f', 6, 64)
		c.JSON(http.StatusOK, rs)
	}
	fmt.Println(rs.Data)

}

func SaleSave(c *gin.Context) {
	log.Println("call POST SaleSave")
	c.Keys = headerKeys

	NewSale := &model.Sale{}
	err := c.BindJSON(NewSale)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Start Controller Create Quotation")
	s, _ := NewSale.SaleSave(dbc)

	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content: " + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = s
		c.JSON(http.StatusOK, rs)
	}
}

func SearchSales(c *gin.Context){
	log.Println("call Get SearchSales")
	c.Keys = headerKeys

	host_code := c.Request.URL.Query().Get("host_code")
	doc_date := c.Request.URL.Query().Get("doc_date")
	keyword := c.Request.URL.Query().Get("keyword")

	NewSale := new(model.Sale)
	sales, err := NewSale.SearchSales(dbc,host_code,doc_date,keyword)
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


func PrintSaleDaily(c *gin.Context){
	log.Println("call Get SearchSales")
	c.Keys = headerKeys

	doc_date := c.Request.URL.Query().Get("doc_date")


	NewSale := new(model.Sale)
	sales, err := NewSale.PrintSaleDaily(dbc,doc_date,model.H)
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

func SearchSaleById(c *gin.Context){
	log.Println("call Get SearchSaleById")
	c.Keys = headerKeys

	strId := c.Param("id")
	id, _ := strconv.ParseInt(strId, 10, 64)

	NewSale := new(model.Sale)
	err := NewSale.SearchSaleById(dbc, id)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = NewSale
		c.JSON(http.StatusOK, rs)
	}

}

