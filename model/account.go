package model

type Account struct {
	ID           int    `json:"id" db:"id"`
	CustomerName string `json:"customer_name" db:"customer_name"`
	Balance      int    `json:"balance" id:"balance"`
}
