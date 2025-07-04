package models

import "time"

type Payment struct {
	IBAN      string    `json:"iban"`
	CVV       string    `json:"cvv"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Nonce     string    `json:"nonce"`
	Timestamp time.Time `json:"timestamp"`
}
