package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"net/http"
	"github.com/itnopadol/api_pos/app/resp"
	"log"
	"fmt"
)

func LogIn(c *gin.Context){
	log.Println("call Get User")
	c.Keys = headerKeys

	user_code := c.Request.URL.Query().Get("user_code")
	password := c.Request.URL.Query().Get("password")


	u := new(model.User)
	err := u.LogIn(dbc, user_code, password)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = u
		c.JSON(http.StatusOK, rs)
	}
}


func ListUser(c *gin.Context){
	log.Println("call Get User")
	c.Keys = headerKeys

	//keyword := c.Request.URL.Query().Get("keyword")

	u := new(model.User)
	users, err := u.ListUser(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = users
		c.JSON(http.StatusOK, rs)
	}
}


func SaveUser(c *gin.Context){
	log.Println("call POST Save User")
	c.Keys = headerKeys

	newUser := &model.User{}
	err := c.BindJSON(newUser)
	if err != nil {
		fmt.Println(err.Error())
	}
	cf := newUser.Save(dbc)
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


func UpdateUser(c *gin.Context){
	log.Println("call PUT Save User")
	c.Keys = headerKeys

	newUser := &model.User{}
	err := c.BindJSON(newUser)
	if err != nil {
		fmt.Println(err.Error())
	}
	cf := newUser.Update(dbc)
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