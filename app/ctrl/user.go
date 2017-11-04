package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"net/http"
	"github.com/itnopadol/api_pos/app/resp"
	"log"
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
