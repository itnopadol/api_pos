package ctrl

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"fmt"
	"net/http"
	"log"
)

func GenUser(c *gin.Context){
	res, err := http.Get("http://hapos.dyndns.org:9003/wifi/genuser.php")
	if err != nil {
		log.Println(err.Error())
		return
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	password := string(robots)

	fmt.Println("robots = >",password)
	c.JSON(http.StatusOK, gin.H{"password":password})

}
