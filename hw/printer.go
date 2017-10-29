package hw


import (
"bufio"
//"github.com/gotmc/libusb"
"github.com/knq/escpos"
"log"
//"os"
"unicode"
	//"github.com/gotmc/libusb"
	//"github.com/gotmc/libusb"
)

type PosPrinter struct {
	*escpos.Escpos
	Writer *bufio.Writer
}


func (pos *PosPrinter) SetCharaterCode(code byte) {
	pos.WriteRaw([]byte{0x1B, 0x74, code})
}

func (pos *PosPrinter) LineFeed() {
	pos.WriteRaw([]byte{0xA})
}

func (pos *PosPrinter) FormFeed() {
	pos.WriteRaw([]byte{0xC})
}

func (pos *PosPrinter) PaperFullCut(n byte) {
	if n == 0 {
		pos.WriteRaw([]byte{0x1D, 0x56, 0})
	} else {
		pos.WriteRaw([]byte{0x1D, 0x56, 65, n})
	}
}

func (pos *PosPrinter) PaperPartialCut(n byte) {
	if n == 0 {
		pos.WriteRaw([]byte{0x1D, 0x56, 1})
	} else {
		pos.WriteRaw([]byte{0x1D, 0x56, 66, n})
	}
}

func (pos *PosPrinter) WriteString(str string) {
	pos.WriteRaw([]byte(str))
}

func (pos *PosPrinter) WriteString3Lines(str string) {
	upperLine, midLine, underLine := pos.ConvertUnicodeToThaiAscii3Lines(str)
	pos.Write(string(upperLine))
	pos.LineFeed()
	pos.Write(string(midLine))
	pos.LineFeed()
	pos.Write(string(underLine))
	pos.LineFeed()
}

func (pos *PosPrinter) BitImagePrintingSingleDensity(n1 byte, n2 byte, bitImageData []byte) {
	escposCmd := []byte{0x1B, 0x2A, 0x60, n1, n2}
	escposCmd = append(escposCmd, bitImageData...)
	pos.WriteRaw(escposCmd)
}

func (pos *PosPrinter) BitImagePrintingDoubleDensity(n1 byte, n2 byte, bitImageData []byte) {
	escposCmd := []byte{0x1B, 0x2A, 0x61, n1, n2}
	escposCmd = append(escposCmd, bitImageData...)
	pos.WriteRaw(escposCmd)
}

func (pos *PosPrinter) SetTextAlignment(align string) {
	switch align {
	case "L", "l":
		pos.WriteRaw([]byte{0x1B, 0x61, 0})
	case "C", "c":
		pos.WriteRaw([]byte{0x1B, 0x61, 1})
	case "R", "r":
		pos.WriteRaw([]byte{0x1B, 0x61, 2})
	}
}

func (pos *PosPrinter) SetTextSize(width byte, height byte) {
	size := height
	size = size | (width << 4)
	pos.WriteRaw([]byte{0x1D, 0x21, size})
}

func (pos *PosPrinter) SetPrintingAreaWidth(n1 byte, n2 byte) {
	escposCmd := []byte{0x1D, 0x57, n1, n2}
	pos.WriteRaw(escposCmd)
}

func (pos *PosPrinter) RegistrationBitImage(id byte, horizontalByte byte, y1 byte, y2 byte, bitImageData []byte) {
	escposCmd := []byte{0x1D, 0x26, id, horizontalByte, y1, y2}
	escposCmd = append(escposCmd, bitImageData...)
	pos.WriteRaw(escposCmd)
}

func (pos *PosPrinter) PrintRegistrationBitImage(id, mode byte) {
	escposCmd := []byte{0x1D, 0x27, id, mode}
	pos.WriteRaw(escposCmd)
}

func (pos *PosPrinter) ConvertUnicodeToThaiAscii3Lines(str string) (upperLines, midLines, underLines []byte) {
	//str := "ฮัลโหลสวัสดีครับ กระผม น้ำใจ ที่นั่น หมูปิ้ง !@#$%^&*()_+"
	for _, c := range str {
		b := c % 256
		//fmt.Printf("%d, ", c)
		if c > 3584 {
			// ถ้ารหัส unicode เป็นตัวอักษรภาษาไทย
			b = 160 + (c % 256)
			if unicode.IsLetter(c) {
				midLines = append(midLines, byte(b))
				upperLines = append(upperLines, byte(32))
				underLines = append(underLines, byte(32))
			}
			if unicode.IsMark(c) {
				//upperLines[len(upperLines)-1] = byte(b)
				switch byte(b) {
				case 216, 217:
					if len(underLines) > 0 {
						underLines[len(underLines)-1] = byte(b)
					} else {
						upperLines = append(upperLines, 32)
						underLines = append(underLines, byte(b))
					}
				default:
					if len(upperLines) > 0 {
						mark := upperLines[len(upperLines)-1]
						if mark == byte(32) {
							upperLines[len(upperLines)-1] = byte(b)
						} else {
							switch mark {
							case 209:
								switch byte(b) {
								case 232:
									upperLines[len(upperLines)-1] = byte(146)
								case 233:
									upperLines[len(upperLines)-1] = byte(147)
								case 234:
									upperLines[len(upperLines)-1] = byte(148)
								case 235:
									upperLines[len(upperLines)-1] = byte(149)
								}
							case 212:
								switch byte(b) {
								case 232:
									upperLines[len(upperLines)-1] = byte(150)
								case 233:
									upperLines[len(upperLines)-1] = byte(151)
								case 234:
									upperLines[len(upperLines)-1] = byte(152)
								case 235:
									upperLines[len(upperLines)-1] = byte(153)
								}
							case 213:
								switch byte(b) {
								case 232:
									upperLines[len(upperLines)-1] = byte(155)
								case 233:
									upperLines[len(upperLines)-1] = byte(156)
								case 234:
									upperLines[len(upperLines)-1] = byte(157)
								case 235:
									upperLines[len(upperLines)-1] = byte(158)
								}
							case 214:
								switch byte(b) {
								case 232:
									upperLines[len(upperLines)-1] = byte(219)
								case 233:
									upperLines[len(upperLines)-1] = byte(220)
								case 234:
									upperLines[len(upperLines)-1] = byte(221)
								case 235:
									upperLines[len(upperLines)-1] = byte(222)
								}
							case 215:
								switch byte(b) {
								case 232:
									upperLines[len(upperLines)-1] = byte(251)
								case 233:
									upperLines[len(upperLines)-1] = byte(252)
								case 234:
									upperLines[len(upperLines)-1] = byte(253)
								case 235:
									upperLines[len(upperLines)-1] = byte(254)
								}
							default:

							}
						}
					} else {
						upperLines = append(upperLines, byte(b))
						underLines = append(underLines, 32)
					}
				}
			}
		} else {
			midLines = append(midLines, byte(b))
			upperLines = append(upperLines, byte(32))
			underLines = append(underLines, byte(32))
		}
	}
	return
}

func (pos *PosPrinter) ForwardLinesFeed(line byte) {
	pos.WriteRaw([]byte{0x1B, 0x64, line})
}

func (pos *PosPrinter) BackwardLinesFeed(line byte) {
	pos.WriteRaw([]byte{0x1B, 0x65, line})
}

func (pos *PosPrinter) PrintQRCode(k1, k2, k3, k4, p byte, data []byte) {
	escposCmd := []byte{0x1D, 0x6B, 0x80, k1, k2, k3, k4, p}
	escposCmd = append(escposCmd, data...)
	pos.WriteRaw(escposCmd)
}

func (pos *PosPrinter) PrintStringQRCode(str string) {
	qrData := []byte(str)
	qrParity := byte(0)
	for i, d := range qrData {
		if i == 0 {
			qrParity = d
		} else {
			qrParity ^= d
		}
	}
	qrData = append(qrData, 0)
	pos.PrintQRCode(0x33, 0xF, qrParity, 15, 0x11, qrData)
}

func (pos *PosPrinter) SetBarcodeWidth(m, n byte) {
	pos.WriteRaw([]byte{0x1D, 0x65, m, n})
}

func (pos *PosPrinter) SetBarcodeHeight(n byte) {
	pos.WriteRaw([]byte{0x1D, 0x68, n})
}

func (pos *PosPrinter) SetLeftMargin(margin int) {
	pos.WriteRaw([]byte{0x1D, 0x4C, byte(margin % 256), byte(margin / 256)})
}

func (pos *PosPrinter) SetHorizontalTabPosition(tabPositions []byte) {
	escposCmd := []byte{0x1B, 0x44}
	escposCmd = append(escposCmd, tabPositions...)
	escposCmd = append(escposCmd, 0x00)
	pos.WriteRaw(escposCmd)
}

func (pos *PosPrinter) ResetHorizontalTabPosition() {
	escposCmd := []byte{0x1B, 0x44, 0x00}
	pos.WriteRaw(escposCmd)
}


func (pos *PosPrinter) End() error {
	err := pos.Writer.Flush()
	if err != nil {
		log.Printf("EscPos writer error : %v", err)
		return err
	}
	return nil
}
func (pos *PosPrinter) WriteStringLines(str string) {
	 midLine := pos.ConvertUnicodeToThaiAscii1Lines(str)
	pos.Write(string(midLine))

	//pos.LineFeed()
}

func (pos *PosPrinter) ConvertUnicodeToThaiAscii1Lines(str string) ( midLines []byte) {
	//str := "ฮัลโหลสวัสดีครับ กระผม น้ำใจ ที่นั่น หมูปิ้ง !@#$%^&*()_+"
	for _, c := range str {
		b := c % 256
		//fmt.Printf("%d, ", c)
		if c > 3584 {
			// ถ้ารหัส unicode เป็นตัวอักษรภาษาไทย
			b = 160 + (c % 256)
			midLines = append(midLines, byte(b))

		} else {
			midLines = append(midLines, byte(b))
			//upperLines = append(upperLines, byte(32))
			//underLines = append(underLines, byte(32))
		}
	}
	return
}