package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"net/http"
	"github.com/itnopadol/api_pos/app/resp"
	"log"
	"fmt"
)

func SearchHost(c *gin.Context) {

	log.Println("call GET Host")
	c.Keys = headerKeys

	h := new(model.Host)

	res, err := h.SearchHost(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = res
		c.JSON(http.StatusOK, rs)
	}
}

func SaveHost(c *gin.Context){
	log.Println("POST Save Host")
	c.Keys = headerKeys

	host := &model.Host{}
	err := c.BindJSON(host)
	if err != nil {
		fmt.Println(err.Error())
	}

	h := host.Save(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = h
		c.JSON(http.StatusOK, rs)
	}

}


func UpdateHost(c *gin.Context){
	log.Println("Call PUT Host")
	c.Keys = headerKeys

	host := &model.Host{}
	err := c.BindJSON(host)
	if err != nil {
		fmt.Println(err.Error())
	}

	h := host.Update(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" +err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = h
		c.JSON(http.StatusOK, rs)
	}
}
