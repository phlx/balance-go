package transactions

import (
	"math"
	"math/rand"
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithSortByAmountAsc(t *testing.T) {
	userID := int64(1)
	carcass := functional.NewCarcass(t)
	carcass.ResetModels(userID)

	models := 5
	shift := 500.0
	maxAmount := 1000.0 + shift
	max := shift
	min := maxAmount

	amounts := make([]float64, models)

	for i := 0; i < models; i++ {
		amount := math.Round(rand.Float64()*maxAmount + shift)
		amounts[i] = amount
		if min > amount {
			min = amount
		}
		if max < amount {
			max = amount
		}
	}

	sort := map[string]string{
		"amount": "asc",
	}

	for i := 0; i < models; i++ {
		carcass.API.GiveExpect(userID, amounts[i], nil, "").Status(http.StatusOK)
	}

	transactions := carcass.API.TransactionsExpectWithSort(userID, sort).
		Status(http.StatusOK).
		JSON().
		Object()

	transactions.Value("transactions").Array().Length().Equal(models)
	transactions.Value("transactions").Array().First().Object().ValueEqual("amount", min)
	transactions.Value("transactions").Array().Last().Object().ValueEqual("amount", max)
	transactions.ValueEqual("pages", 1)

	carcass.ResetModels(userID)
}
