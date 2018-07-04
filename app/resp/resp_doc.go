package resp


type ResponseDoc struct {
	Status  string `json:"status"`
	Message string  `json:"message,omitempty"`
	PrintBill string `json:"print_bill" db:"print_bill"`
	PrintKitchecn string `json:"print_kitchecn" db:"print_kitchecn"`
	PrintBar string `json:"print_bar" db:"print_bar"`
	Data    interface{} `json:"data,omitempty"`
}
