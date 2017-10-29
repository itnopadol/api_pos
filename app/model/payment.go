package model

type Payment struct {
	pay        float64 // รับเงินมาทั้งหมด
	total      float64 // มูลค่าเงินพักทั้งหมด
	remain     float64 // เงินคงค้างชำระ
	change     float64 // เงินทอน
}
