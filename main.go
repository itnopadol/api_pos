package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gin-contrib/cors.v1"
	"github.com/itnopadol/api_pos/app/ctrl"
	_ "github.com/itnopadol/api_pos/app/model"
)

type item struct {
	name string
	qty  int
}

type items struct {
	items []item
}

func main() {

	r := gin.New()
	r.Use(cors.Default())

	r.POST("/menu", ctrl.SaveMenu)
	r.PUT("/menu", ctrl.UpdateMenu)
	r.GET("/menu", ctrl.GetMenu)

	r.GET("/menu/:id", ctrl.GetItemsByMenuId)
	r.GET("/item/:id", ctrl.GetItemById)
	r.GET("/items/print", ctrl.PrintTest)
	r.POST("/item", ctrl.SaveItem)
	r.PUT("/item", ctrl.UpdateItem)

	//r.POST("/menu", ctrl.SaveMenu)

	r.POST("/sale/change", ctrl.ShowChangeAmount)
	r.POST("/sale", ctrl.SaleSave)
	r.GET("/sales", ctrl.SearchSales)
	r.GET("/sale/:id", ctrl.SearchSaleById)
	r.POST("/sale/docno", ctrl.SearchSaleByDocNo)
	r.GET("/sales/saledaily", ctrl.PrintSaleDailyTotal)
	r.GET("/sales/salenetdaily", ctrl.PrintSaleNetDaily)
	r.PUT("/sale/void", ctrl.SaleVoid)
	r.GET("/sales/reportsale", ctrl.ReportSaleDaily)
	r.GET("/sales/reportsalebymenu", ctrl.ReportSaleDailyByMenu)
	
	r.GET("/saletype/search",ctrl.SearchSaleType)

	r.POST("/shift/open", ctrl.SaveShift)
	r.PUT("/shift/update", ctrl.UpdateShift)
	r.PUT("/shift/closed", ctrl.ClosedShift)
	r.GET("/shift/search", ctrl.ShiftDetails)
	r.GET("/shift/list", ctrl.ShiftList)
	r.GET("/shift/last/id", ctrl.ShiftLastID)
	r.GET("/shift/senddaily", ctrl.PrintSendDailyTotal)

	r.GET("/report/tax", ctrl.ReportTax)

	r.GET("/host", ctrl.SearchHost)
	r.POST("/host", ctrl.SaveHost)
	r.PUT ("/host" ,ctrl.UpdateHost)

	r.GET("/user/login", ctrl.LogIn)
	//r.GET("/user", ctrl.SearchUser)
	r.GET("/users", ctrl.ListUser)
	r.POST("/user", ctrl.SaveUser)
	r.PUT("/user", ctrl.UpdateUser)

	r.GET("/config", ctrl.GenWifiPassword)
	r.GET("/config/search", ctrl.Search)
	r.POST("/config", ctrl.SaveConfig)
	r.PUT("/config", ctrl.UpdateConfig)

	r.GET("/gentax", ctrl.GenTaxData)
	r.GET("/gentaxall", ctrl.GenTaxWithNoVatData)

	r.Run(":8888")

}
