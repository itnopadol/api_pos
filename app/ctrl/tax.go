package ctrl

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
)

func GenTaxData(c *gin.Context) {
	log.Println("call Get SearchSales")
	c.Keys = headerKeys

	company_id := c.Request.URL.Query().Get("company_id")
	branch_id := c.Request.URL.Query().Get("branch_id")
	tax_amount := c.Request.URL.Query().Get("tax_amount")
	begin_date := c.Request.URL.Query().Get("begin_date")
	end_date := c.Request.URL.Query().Get("end_date")

	company_id1, err := strconv.ParseInt(company_id, 10, 64)
	branch_id1, err := strconv.ParseInt(branch_id, 10, 64)
	tax_amount1, err := strconv.ParseFloat(tax_amount, 64)

	tax := new(model.TaxData)
	err = tax.GenTaxData(dbc, company_id1, branch_id1, begin_date, end_date, tax_amount1)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = tax
		c.JSON(http.StatusOK, rs)
	}

}

func GenTaxWithNoVatData(c *gin.Context) {
	log.Println("call Get SearchSales")
	c.Keys = headerKeys

	tax_total_amount := c.Request.URL.Query().Get("tax_total_amount")
	no_vat := c.Request.URL.Query().Get("no_vat")
	begin_date := c.Request.URL.Query().Get("begin_date")
	end_date := c.Request.URL.Query().Get("end_date")

	send_tax_total, err := strconv.ParseFloat(tax_total_amount, 64)
	no_vat_amount, err := strconv.ParseFloat(no_vat, 64)

	tax := new(model.TaxData)
	err = tax.GenTaxWithNoVatData(dbc, begin_date, end_date, no_vat_amount, send_tax_total)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = tax
		c.JSON(http.StatusOK, rs)
	}

}
