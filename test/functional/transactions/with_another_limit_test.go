package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithAnotherLimit(t *testing.T) {
	userID := int64(1)
	carcass := functional.NewCarcass(t)
	carcass.ResetModels(userID)

	models := 15
	limit := 12
	page := 2
	expectedModels := 3

	for amount := 1; amount <= models; amount++ {
		carcass.API.GiveExpect(userID, float64(amount), nil, "").Status(http.StatusOK)
	}

	transactions := carcass.API.TransactionsExpectWithPageAndLimit(userID, limit, page).
		Status(http.StatusOK).
		JSON().
		Object()

	transactions.Value("transactions").Array().Length().Equal(expectedModels)
	transactions.ValueEqual("pages", page)

	carcass.ResetModels(userID)
}
