package model

import (
	"time"
)

type Host struct {
	Id                    string  // รหัสเมนบอร์ดตู้
	LastTicketNumber      int     // เลขคิวตั๋วล่าสุด ของแต่ละวัน ปิดเครื่องต้องยังอยู่ ขึ้นวันใหม่ต้อง Reset
	Status                HostStatus
	LogoImageId           int
	LogoImageWidth        int
	ApiServerUrl          string
	FirmwareUpdateDone    chan bool
	StatusPollingInterval time.Duration // ส่งสถานะทุกๆกี่วินาที , 0 = ไม่ส่ง
	XToken				  string // token ที่ได้จากการ authen api server
	ShopTitle			  string
	ShopQRCode			  string
	ShopWebsite			  string
	Email			  	  string
	Password			  string
}

type HostStatus struct {
	UUID                string `json:"host_uuid"`
	BillAccOnline       bool   `json:"bill_acc_online"`
	BillAccStatus       string `json:"bill_acc_status"`
	CoinAccOnline       bool   `json:"coin_acc_online"`
	CoinAccStatus       string `json:"coin_acc_status"`
	CoinHopperOnline    bool   `json:"coin_hopper_online"`
	NoteDispenserOnline bool   `json:"note_dispenser_online"`
	PrinterOnline       bool   `json:"printer_online"`
	PrinterStatus       string `json:"printer_status"`
	GsmOnline           bool   `json:"gsm_online"`    // สถานะ GSM ปัจจุบัน (Real time)
	ServerOnline        bool   `json:"server_online"` // สถานะเซิร์ฟเวอร์ครั้งสุดท้ายที่สื่อสาร
	FrontDoorStatus     string `json:"front_door_status"`
	CashDoorStatus      string `json:"cash_door_status"`
}


