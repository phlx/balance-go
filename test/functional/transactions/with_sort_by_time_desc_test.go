package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithSortByTimeDesc(t *testing.T) {
	userID := int64(1)
	carcass := functional.NewCarcass(t)
	carcass.ResetModels(userID)

	models := 10

	sort := map[string]string{
		"time": "desc",
	}

	for amount := 1; amount <= models; amount++ {
		carcass.API.GiveExpect(userID, float64(amount), nil, "").Status(http.StatusOK)
	}

	transactions := carcass.API.TransactionsExpectWithSort(userID, sort).
		Status(http.StatusOK).
		JSON().
		Object()

	transactions.Value("transactions").Array().Length().Equal(models)
	transactions.Value("transactions").Array().First().Object().ValueEqual("amount", float64(models))
	transactions.ValueEqual("pages", 1)

	carcass.ResetModels(userID)
}
