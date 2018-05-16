package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
	"log"
	"net/http"
	"strconv"
)

func GenTaxData(c *gin.Context){
	log.Println("call Get SearchSales")
	c.Keys = headerKeys

	tax_amount := c.Request.URL.Query().Get("tax_amount")
	begin_date := c.Request.URL.Query().Get("begin_date")
	end_date := c.Request.URL.Query().Get("end_date")

	tax_amount1, err := strconv.ParseFloat(tax_amount, 64)

	Tax := new(model.TaxData)
	err = Tax.GenTaxData(dbc,begin_date,end_date,tax_amount1)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = nil
		c.JSON(http.StatusOK, rs)
	}

}

