package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"net/http"
	"github.com/itnopadol/api_pos/app/resp"
)

func SearchHost(c *gin.Context){
	h := new(model.Host)

	res, err := h.SearchHost(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = res
		c.JSON(http.StatusOK, rs)
	}
}
