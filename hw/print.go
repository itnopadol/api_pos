package hw

import (
	"time"
	"fmt"
	"errors"
	"strconv"
	"strings"
	"log"
	"os"
)

type Printer struct {
	machineId string `json:"machine_id"`
	Status    string
	Online    bool
	PosPrinter
	NoChange float64
	QueueNum int
}

func (p *Printer) makeTicket(s *Sale) doGroup {
	var g doGroup
	g.setTextSize(0)
	g.printLine("=========== ร้านกาแฟ สุขใจขายได้สบายดี ===========")
	var today = time.Now()
	date := fmt.Sprintf("%s", today.Format(time.RFC3339Nano))
	g.printLine(date)
	ss := s.SaleSubs
	for _, sub := range ss {
		item := fmt.Sprintf("%2dx%-17s%4s", sub.Qty, sub.ItemName, sub.PriceName)
		//detail := fmt.Sprintf("@%3.2f %4.2f฿", float64(sub.Qty), sub.Price)
		detail := fmt.Sprintf("@%3.2d %4.2f฿", sub.Qty, sub.Price)
		g.setTextSize(1)
		g.printLine(item)
		g.setTextSize(0)
		g.print(detail)
		g.newLine()
	}
	sum := fmt.Sprintf("%6s%4.2f%6s%6.2f%6s%6.2f", "Total:", s.Total, " Payment:", s.Pay, " Change:", s.Change)
	g.print(sum)
	g.newLine()
	g.printLine("==============================================")
	//g.printBarcode("CODE39", "12345678")
	g.paperCut("full_cut", 90)
	fmt.Println(&g.actions)
	return g
}

// Print() รับค่า Sale แล้วพิมพ์ฟอร์มที่กำหนด ให้
func (p *Printer) PrintTicket(s *Sale, h *Host) error {
	fmt.Println("p.Print() run")
	p.CheckStatus(h)
	p.GetStatus()
	if !p.Online {
		return errors.New("Printer offline!")
	}
	p.Init()
	p.SetLeftMargin(0)
	p.SetTextSize(0, 0)
	p.SetCharaterCode(26) // 26 = Thai code 18
	// Print logo image
	p.SetLeftMargin((576 - h.LogoImageWidth) / 2)
	p.PrintRegistrationBitImage(byte(h.LogoImageId), 0)
	p.LineFeed()
	// Print queue number
	h.LastTicketNumber = p.NextQueueNumber()
	qStr := strconv.Itoa(h.LastTicketNumber)
	qLen := len(qStr)
	lMargin := 288 - int(16*(qLen/2)) // Set text center
	p.SetLeftMargin(lMargin)
	p.SetTextSize(1, 1)
	p.WriteString(qStr)
	p.SetTextSize(0, 0)
	p.LineFeed()

	_, midLine, _ := p.ConvertUnicodeToThaiAscii3Lines(h.ShopTitle)
	lMargin = 288 - int(12*(len(midLine)/2)) // Set text center
	p.SetLeftMargin(lMargin)
	p.WriteString3Lines(h.ShopTitle)
	p.SetLeftMargin(0)
	var today = time.Now()
	date := fmt.Sprintf("Date %s Time %s", today.Format("02/01/2006"), today.Format("15:04:05"))
	p.WriteString(date)
	p.LineFeed()
	//p.SetLineSpacing(5)
	ss := s.SaleSubs
	for _, sub := range ss {
		//item := fmt.Sprintf("%2dx%-17s%4s", sub.Qty, sub.ItemName, sub.PriceName)
		item := fmt.Sprintf("%2dx%s/%s", sub.Qty, sub.ItemName, sub.PriceName)
		//detail := fmt.Sprintf("@%3.2f %4.2f฿", float64(sub.Qty), sub.Price)
		detail := fmt.Sprintf("%4.2f฿", float64(sub.Qty)*sub.Price)
		p.SetTextSize(0, 0)
		_, itemMidLine, _ := p.ConvertUnicodeToThaiAscii3Lines(item)
		tab := len(itemMidLine) / 8                                      // หาจำนวนแทบของข้อความ โดยนำข้อความ/8
		p.WriteString3Lines(item + strings.Repeat("\t", 5-tab) + detail) // พิมพ์ราคาที่ตำแหน่ง tab ที่ 5
	}
	sum := fmt.Sprintf("%6s%4.2f%6s%6.2f%6s%6.2f", "Total:", s.Total, " Payment:", s.Pay, " Change:", s.Change)
	p.WriteString(sum)
	p.LineFeed()
	// พิมพ์เงินทอนค้างจ่าย
	if p.NoChange > 0 {
		p.SetTextSize(1, 1)
		p.WriteString(fmt.Sprintf("%6s%6.2f", "Unpaid Change:", p.NoChange))
		p.SetTextSize(0, 0)
		p.LineFeed()
	}
	p.WriteString("==============================================")
	p.LineFeed()
	p.PrintStringQRCode(h.ShopQRCode)
	p.WriteString(h.ShopWebsite)
	p.ForwardLinesFeed(10)
	p.PaperFullCut(0)

	fmt.Println("1. สั่งพิมพ์ รอ Priner ตอบสนอง...")
	err := p.End()
	if err != nil {
		fmt.Println("พิมพ์ไม่สำเร็จ Print error!")
		m2 := &Message{
			Device:  "ui",
			Command: "print",
			Type:    "event",
			Data:    "error",
		}
		h.Web.Send <- m2
		return err
	}
	//data := p.makeTicket(s)
	//data := gin.H{"action": "printline", "action_data": "นี่คือคูปอง"}
	//m := &Message{
	//	Device:  "printer",
	//	Command: "do_group",
	//	Type:    "request",
	//	Data:    data.actions,
	//}
	fmt.Println("พิมพ์สำเร็จ Print success!")
	m2 := &Message{
		Device:  "ui",
		Command: "print",
		Type:    "event",
		Data:    "success",
	}
	h.Web.Send <- m2
	return nil
}

func (p *Printer) makeRefund(value float64) error {
	return nil
}

func (p *Printer) CheckStatus(h *Host) {
	status, err := p.GetHwStatus()
	if err != nil {
		log.Printf("Printer CheckStatus() error : %v", err)
		return
	}
	isOnline := status[0] & 0x08 // Check On-line/off-line status
	if isOnline == 0 {
		p.Online = true
		p.event("online", h)
	} else {
		p.Online = false
		p.event("offline", h)
	}
	isPaperNearEnd := status[2] & 0x01 // Check Paper near-end
	if isPaperNearEnd != 0 {
		p.Status = "near_end"
		p.event("near_end", h)
	}
	isPaperOut := status[2] & 0x04 // Check Paper-out
	if isPaperOut != 0 {
		p.Status = "no_paper"
		p.event("no_paper", h)
	}
}

func (p *Printer) GetStatus() {
	status, err := p.GetHwStatus()
	if err != nil {
		p.Online = false
		p.Status = "offline"
		log.Printf("Printer GetHwStatus() error : %v", err)
		return
	}
	isOnline := status[0] & 0x08 // Check On-line/off-line status
	if isOnline == 0 {
		p.Online = true
		p.Status = "ok"
	} else {
		p.Online = false
		p.Status = "offline"
	}
	isPaperNearEnd := status[2] & 0x01 // Check Paper near-end
	if isPaperNearEnd != 0 {
		p.Status = "near_end"
	}
	isPaperOut := status[2] & 0x04 // Check Paper-out
	if isPaperOut != 0 {
		p.Status = "no_paper"
	}
}

//=============== ACTION ====================
type action struct {
	Name string      `json:"action"`
	Data interface{} `json:"action_data"`
}

//=============== DO_GROUP ====================
type doGroup struct {
	actions []*action
}

func (g *doGroup) print(s string) {
	a := &action{
		Name: "print",
		Data: s,
	}
	g.actions = append(g.actions, a)
}

func (g *doGroup) printLine(s string) {
	a := &action{"printline", s}
	g.actions = append(g.actions, a)
}

func (g *doGroup) setTextSize(size int) {
	a := &action{"set_text_size", size}
	g.actions = append(g.actions, a)
}

func (g *doGroup) newLine() {
	a := &action{Name: "newline"}
	g.actions = append(g.actions, a)
}

//=========== BARCODE =============
type barcode struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func (g *doGroup) printBarcode(t, d string) {
	data := barcode{t, d}
	a := &action{
		Name: "print_barcode",
		Data: data,
	}
	g.actions = append(g.actions, a)
}

//=========== QR_CODE =============
type qrCode struct {
	Mag      int    `json:"mag"`
	Ecl      int    `json:"ect"`
	DataType string `json:"data_type"`
	Data     string `json:"data"`
}

func (g *doGroup) printQr(mag, ecl int, data_type, data string) {
	d := qrCode{mag, ecl, data_type, data}
	a := &action{
		Name: "print_qr",
		Data: d,
	}
	g.actions = append(g.actions, a)
}

//=========== PAPER_CUT =============
type paperCut struct {
	Type string `json:"type"`
	Feed int    `json:"feed"`
}

func (g *doGroup) paperCut(t string, f int) {
	data := paperCut{
		Type: t,
		Feed: f,
	}
	a := &action{
		Name: "paper_cut",
		Data: data,
	}
	g.actions = append(g.actions, a)
}

func (p *Printer) QueueNumber() int {
	today := time.Now()
	qNumber := 1
	qf := fmt.Sprintf("%s", today.Format("2006-01-02"))
	lp, err := os.Open(qf)
	if err != nil {
		log.Println(err)
		return qNumber
	}
	lp.Close()
	q := make([]byte, 5)
	n, err := lp.Read(q)
	if n > 0 {
		qNumber, _ = strconv.Atoi(string(q[:n]))
	}
	return qNumber
}

func (p *Printer) NextQueueNumber() int {
	today := time.Now()
	qNumber := 1
	qf := fmt.Sprintf("%s", today.Format("2006-01-02"))
	lp, err := os.OpenFile(qf, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		log.Println(err)
		return qNumber
	}
	defer lp.Close()
	q := make([]byte, 5)
	n, err := lp.Read(q)
	if n > 0 {
		qNumber, _ = strconv.Atoi(string(q[:n]))
		qNumber += 1
	}
	lp.WriteAt([]byte(strconv.Itoa(qNumber)), 0)
	return qNumber
}

