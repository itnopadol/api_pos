package ctrl
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
	"github.com/itnopadol/api_pos/app/model"
	"github.com/itnopadol/api_pos/app/resp"
	"strconv"
	"log"
	"net/http"
	"fmt"
)

func GetItemById(c *gin.Context) {
	var item model.Item

	log.Println("call GET Item")
	c.Keys = headerKeys

	strId := c.Param("id")
	id, _ := strconv.ParseInt(strId, 10, 64)
	err := item.Get(dbc, id)

	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content and Error :"+ err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = item
		c.JSON(http.StatusOK, rs)
	}
}

func GetItemsByMenuId(c *gin.Context) {
	fmt.Println("call GetItemsByMenuId")
	c.Keys = headerKeys

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println("error:", err)
	}
	//var item model.Item
	item := new(model.Item)
	langs, err := item.ByMenuId(dbc, id)
	if err != nil {
		//ctx.HTML(http.StatusNotFound, "error.tpl", err.Error())
		c.JSON(http.StatusNotFound, err.Error())
	}
	c.JSON(http.StatusOK, langs)
}

func SaveItem(c *gin.Context){
	fmt.Println("Call POST SaveItem")
	c.Keys = headerKeys

	newItem := &model.Item{}
	err := c.BindJSON(newItem)
	if err != nil {
		fmt.Println(err.Error())
	}
	i := newItem.SaveItem(dbc)
	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content : "+err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = i
		c.JSON(http.StatusOK, rs)
	}
}

func UpdateItem(c *gin.Context){
	fmt.Println("Call PUT UpdateItem")
	c.Keys = headerKeys

	newItem := &model.Item{}
	err := c.BindJSON(newItem)
	if err != nil {
		fmt.Println(err.Error())
	}
	i := newItem.UpdateItem(dbc)

	rs := resp.Response{}
	if err != nil {
		rs.Status = "error"
		rs.Message = "No Content : "+err.Error()
		c.JSON(http.StatusNotFound, rs)
	}else{
		rs.Status = "success"
		rs.Data = i
		c.JSON(http.StatusOK,rs)
	}
}


func PrintTest(c *gin.Context) {
	fmt.Println("call Print Test")
	c.Keys = headerKeys

	//var item model.Item
	item := new(model.Item)
	err := item.PrintTest(dbc)
	if err != nil {
		//ctx.HTML(http.StatusNotFound, "error.tpl", err.Error())
		c.JSON(http.StatusNotFound, err.Error())
	}
	c.JSON(http.StatusOK, nil)
}