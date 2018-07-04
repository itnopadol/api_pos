package ctrl

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
)

func ReportTax(c *gin.Context){
	log.Println("call Get ReportTax")
	c.Keys = headerKeys

	report_month := c.Request.URL.Query().Get("report_month")
	report_year := c.Request.URL.Query().Get("report_year")

	r := new(model.Report)
	err := r.ReportTax(dbc, report_month, report_year)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = r
		c.JSON(http.StatusOK, rs)
	}
}