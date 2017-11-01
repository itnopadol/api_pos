package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gin-contrib/cors.v1"
	"github.com/itnopadol/hapos_api/app/ctrl"
	_ "github.com/itnopadol/hapos_api/app/model"
)

type item struct {
	name string
	qty int
}

type items struct {
	items []item
}

const (
	printerIP = "192.168.0.206:9100"
	dbPort = "5432"
	dbHost = "localhost"
	dbUser = "paybox"
	dbPass = "paybox"
	dbName = "paybox_vending"
	sslMode = "disable"
)

func main() {

	r := gin.New()
	r.Use(cors.Default())

	r.GET("/item/:id", ctrl.GetItemById)
	r.GET("/menu/:id", ctrl.GetItemsByMenuId)
	r.GET("/menu", ctrl.GetMenu)

	r.POST("/sale/change", ctrl.ShowChangeAmount)
	r.POST("/sale", ctrl.SaleSave)
	r.GET("/sales", ctrl.SearchSales)
	r.GET("/sale/:id", ctrl.SearchSaleById)

	r.POST("/shift", ctrl.SaveShift)
	r.PUT("/shift", ctrl.UpdateShift)
	r.PUT("/shift/closed", ctrl.ClosedShift)
	r.GET("/shift/:host_id", ctrl.ShiftDetails)

	r.Run(":8888")




}



