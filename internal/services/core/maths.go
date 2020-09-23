package core

import (
	"github.com/shopspring/decimal"
)

func round(number float64) float64 {
	result, _ := decimal.NewFromFloat(number).RoundBank(2).Float64()
	return result
}

func add(a, b float64) float64 {
	result, _ := decimal.NewFromFloat(a).Add(decimal.NewFromFloat(b)).RoundBank(2).Float64()
	return result
}

func sub(a, b float64) float64 {
	result, _ := decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)).RoundBank(2).Float64()
	return result
}

func div(a, b float64) float64 {
	result, _ := decimal.NewFromFloat(a).Div(decimal.NewFromFloat(b)).Round(2).Float64()
	return result
}
