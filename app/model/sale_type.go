package model

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SaleType struct {
	SaleTypeID          int    `json:"sale_type_id" db:"sale_type_id"`
	SaleTypeName        string `json:"sale_type_name" db:"sale_type_name"`
	SaleTypeDescription string `json:"sale_type_description" db:"sale_type_description"`
	IsPay               int    `json:"is_pay" db:"is_pay"`
}

func (s *SaleType) SearchSaleType(db *sqlx.DB) (saleType []*SaleType, err error) {
	query := `SELECT sale_type_id, ifnull(sale_type_name,'') as sale_type_name, ifnull(sale_type_description,'') as sale_type_description,ifnull(is_pay,1) as is_pay FROM sale_type where is_cancel = 0 order by sale_type_id	`
	err = db.Select(&saleType, query)
	if err != nil {
		return nil, err
	}

	fmt.Println("s=", saleType)

	return saleType, err
}
