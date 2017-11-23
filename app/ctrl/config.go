package ctrl

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
	"net/http"
)

func GenWifiPassword(c *gin.Context){
	log.Println("call Get Wifi")
	c.Keys = headerKeys

	cf := new(model.Config)
	err := cf.Search(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = cf
		c.JSON(http.StatusOK, rs)
	}
}
