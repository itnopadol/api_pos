package ctrl

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
	"net/http"
	"fmt"
)

func GenWifiPassword(c *gin.Context){
	log.Println("call Get Wifi")
	c.Keys = headerKeys

	cf := new(model.Config)
	err := cf.GenWifiPassword(dbc)
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


func SaveConfig(c *gin.Context){
	log.Println("call POST Save Config")
	c.Keys = headerKeys

	newConfig := &model.Config{}
	err := c.BindJSON(newConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	cf := newConfig.Save(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content : "+err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = cf
		c.JSON(http.StatusOK, rs)
	}
}

func Search(c *gin.Context){
	log.Println("call Get Config")
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


func UpdateConfig(c *gin.Context){
	fmt.Println("Call PUT Update Config")
	c.Keys = headerKeys

	newConfig := &model.Config{}
	err := c.BindJSON(newConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	cf := newConfig.Update(dbc)

	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content : "+err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "Success"
		rs.Data = cf
		c.JSON(http.StatusOK,rs)
	}
}