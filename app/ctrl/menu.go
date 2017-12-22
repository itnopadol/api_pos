package ctrl

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
	"net/http"
	"fmt"
)

func GetMenu(c *gin.Context) {
	log.Println("call GET Menu")
	c.Keys = headerKeys

	var menu model.Menu

	langs, err := menu.Index(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = langs
		c.JSON(http.StatusOK, rs)
	}

}

func SaveMenu(c *gin.Context) {
	log.Println("Call POST Menu")
	c.Keys = headerKeys

	menu := &model.Menu{}
	err := c.BindJSON(menu)
	if err != nil {
		fmt.Println(err.Error())
	}

	m := menu.Save(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :" + err.Error()
		c.JSON(http.StatusNotFound, rs)
	} else {
		rs.Status = "success"
		rs.Data = m
		c.JSON(http.StatusOK, rs)
	}
}

func UpdateMenu(c *gin.Context) {
	log.Println("Call PUT Menu")
	c.Keys = headerKeys

	menu := &model.Menu{}
	err := c.BindJSON(menu)
	if err != nil {
		fmt.Println(err.Error())
	}

	m := menu.Update(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = m
		c.JSON(http.StatusOK, rs)
	}
}
