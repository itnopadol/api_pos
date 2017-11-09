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

	//s := strconv.FormatFloat(3.00, 'E', -1, 64)
	//fmt.Println("SSS=",s)
	//
	//f := 64.0000
	//s1 := strconv.FormatFloat(f, 'g', 2, 64)
	//fmt.Println(s1)
	//s4 := strconv.FormatFloat(f, 'f', 2, 64)
	//fmt.Println(s4)
	//
	//s2 := strconv.FormatFloat(f, 'g', 2, 64)
	//fmt.Println(s2)
	//s3 := strconv.FormatFloat(f, 'f', 2, 64)
	//fmt.Println(s3)
	//
	//i := 5
	//fl:= float64(i)
	//fmt.Printf("f is %f\n", fl)


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