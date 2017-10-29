package main

import (
	//"bufio"
	//"os"

	//"github.com/knq/escpos"
	//"github.com/itnopadol/hapos_api/hw"
	//"net"
	"github.com/gin-gonic/gin"
	"gopkg.in/gin-contrib/cors.v1"
	"github.com/itnopadol/hapos_api/app/ctrl"
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
	//f, err := os.Open("/dev/usb/lp3")
	//f, err :=os.OpenFile("/dev/ttyUSB0", os.O_WRONLY, os.ModeDevice)

	r := gin.New()
	r.Use(cors.Default())

	r.GET("/item/:id", ctrl.GetItemById)
	r.GET("/menu/:id/", ctrl.GetItemsByMenuId)
	r.GET("/menu", ctrl.GetMenu)

	r.POST("/sale/change", ctrl.ShowChangeAmount)
	r.POST("/sale", ctrl.SaleSave)
	r.GET("/sales", ctrl.SearchSales)
	r.GET("/sale/:id", ctrl.SearchSaleById)


	r.Run(":8888")

	//f, err := net.Dial("tcp", printerIP)
	//
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//
	//w := bufio.NewWriter(f)
	//p := escpos.New(w)
	//
	//pt := hw.PosPrinter{p,w}
	//pt.Init()
	//pt.SetLeftMargin(20)
	////pt.PrintRegistrationBitImage(0, 0)
	//pt.WriteRaw([]byte{29,	40,	76,	6,	0,	48,	85,	32,	32,10,10 })
	//printHeader(pt)
	//printCompanyInfo(pt)
	//printDetail(pt)
	//printFooter(pt)
	//
	//printKitchenHeader(pt)
	//printKitchenDetail(pt)
	//printKitchenFooter(pt)


	//pt.WriteRaw([]byte{27,112,0,25,250})
	//pt.Pulse()
}

//
//func printHeader(pt hw.PosPrinter) {
//	pt.SetCharaterCode(26)
//	pt.SetAlign("center")
//	pt.SetTextSize(1, 1)
//	pt.WriteStringLines("คิวเลขที่ 35")
//	pt.LineFeed()
//	pt.SetTextSize(0, 0)
//}
//
//func printKitchenHeader(pt hw.PosPrinter) {
//	pt.SetCharaterCode(26)
//	pt.SetAlign("center")
//	pt.SetTextSize(1, 1)
//	pt.WriteStringLines("คิวเลขที่ 35")
//	pt.LineFeed()
//	pt.WriteStringLines("Kitchen Slip")
//
//	pt.LineFeed()
//	pt.SetTextSize(0, 0)
//	makeline(pt)
//}
//
//
//func printCompanyInfo(pt hw.PosPrinter) {
//	pt.SetFont("B")
//	pt.WriteStringLines("===== ร้านมนต์นมสด =====\n")
//	pt.SetAlign("left")
//	pt.WriteStringLines("เลขประจำตัวผู้เสียภาษี 999999999999")
//	pt.SetAlign("right")
//
//	pt.WriteStringLines("	Cashier : XXXX\n")
//
//	pt.WriteStringLines("วันที่ : 01/09/2017 09:34น.")
//	pt.WriteStringLines("   เลขที่ : 0120171005-001")
//	pt.LineFeed()
//	makeline(pt)
//
//}
//func printFooter(pt hw.PosPrinter) {
//	pt.SetFont("B")
//	pt.WriteStringLines("รวมเป็นเงิน ")
//	pt.WriteStringLines("				")
//	pt.WriteStringLines("30 บาท\n")
//	makeline(pt)
//	// Footer Area
//	pt.SetFont("A")
//	pt.SetAlign("center")
//	pt.WriteStringLines("รหัสผ่าน Wifi : 999999999")
//	pt.Formfeed()
//	pt.Write("*** Completed ***")
//	pt.Formfeed()
//	pt.Cut()
//	pt.End()
//}
//func printDetail(pt hw.PosPrinter) {
//	pt.SetFont("B")
//	pt.WriteStringLines("    1. คาปูชิโน่")
//	pt.WriteStringLines("		")
//	pt.WriteStringLines("	1 ชิ้น\n")
//	//pt.LineFeed()
//	pt.WriteStringLines("    2. เอสเพรสโซ่")
//	pt.WriteStringLines("		")
//	pt.WriteStringLines("	1 ชิ้น\n")
//	pt.WriteStringLines("    3. น้ำเปล่าตราช้าง")
//	pt.WriteStringLines("		")
//	pt.WriteStringLines("	1 ชิ้น\n")
//	pt.FormfeedN(3)
//	makeline(pt)
//}
//
//func printKitchenDetail(pt hw.PosPrinter) {
//	pt.SetTextSize(1,1)
//	pt.SetAlign("left")
//	pt.WriteStringLines("1. CAPU")
//	pt.WriteStringLines(" x ")
//	pt.WriteStringLines(" 1 \n")
//	pt.LineFeed()
//	pt.WriteStringLines("2. ESP ")
//	pt.WriteStringLines(" x ")
//	pt.WriteStringLines(" 1 \n")
//	pt.FormfeedN(3)
//	pt.SetTextSize(1,1)
//	makeline(pt)
//
//}
//
//
//func printKitchenFooter(pt hw.PosPrinter) {
//	// Footer Area
//	pt.SetFont("B")
//	pt.SetAlign("center")
//	pt.Formfeed()
//	pt.Write("*** Completed ***")
//	pt.Formfeed()
//	pt.Cut()
//	pt.End()
//}
//func makeline(pt hw.PosPrinter) {
//	pt.SetTextSize(0,0)
//	pt.SetFont("A")
//	pt.WriteStringLines("-----------------------------------------\n")
//}
