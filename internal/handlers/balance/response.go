package balance

import (
	"balance/internal/exchangerates"
)

type Response struct {
	Currency exchangerates.Currency `json:"currency"`
	Balance  float64                `json:"balance"`
}
