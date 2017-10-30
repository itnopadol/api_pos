package model


type Host struct {
	id 	int `json:"id" db:"id"`
	HostCode	string `json:"host_code" db:"host_code"`
	HostName 	string `json:"host_name" db:"host_name"`
	Status	int `json:"status" db:"status"`
	Active int `json:"active" db:"active"`

}
