package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gin-contrib/cors.v1"
	"github.com/itnopadol/api_pos/app/ctrl"
	_ "github.com/itnopadol/api_pos/app/model"
)

type item struct {
	name string
	qty int
}

type items struct {
	items []item
}

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

	r.POST("/shift/open", ctrl.SaveShift)
	r.PUT("/shift/update", ctrl.UpdateShift)
	r.PUT("/shift/closed", ctrl.ClosedShift)
	r.GET("/shift/search", ctrl.ShiftDetails)
	r.GET("/shift/senddaily", ctrl.PrintSendDailyTotal)

	r.GET("/host", ctrl.SearchHost)

	r.GET("/user/login", ctrl.LogIn)

	r.Run(":8888")

}



