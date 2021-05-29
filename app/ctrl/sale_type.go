package ctrl

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
)

func SearchSaleType(c *gin.Context) {
	log.Println("call Get SearchSaleType")
	c.Keys = headerKeys

	NewSaleType := new(model.SaleType)
	s, err := NewSaleType.SearchSaleType(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, s)
	} else {
		rs.Status = "success"
		rs.Data = NewSaleType
		c.JSON(http.StatusOK, s)
	}

}