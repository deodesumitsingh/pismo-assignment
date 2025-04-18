package model

import "time"

type Transaction struct {
	ID              int       `json:"id"`
	AccountId       int       `json:"account_id"`
	OperationTypeId int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
}
