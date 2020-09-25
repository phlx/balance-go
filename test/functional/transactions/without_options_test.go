package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithoutOptions(t *testing.T) {
	userID := int64(1)
	amountFirst := 1.0
	amountSecond := 2.0
	amountThird := 3.0
	carcass := functional.NewCarcass(t)
	carcass.ResetModels(userID)

	carcass.API.GiveExpect(userID, amountFirst, nil, "").Status(http.StatusOK)
	carcass.API.GiveExpect(userID, amountSecond, nil, "").Status(http.StatusOK)
	carcass.API.GiveExpect(userID, amountThird, nil, "").Status(http.StatusOK)

	transactions := carcass.API.TransactionsExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array()

	transactions.Element(0).Object().ValueEqual("amount", amountThird)
	transactions.Element(1).Object().ValueEqual("amount", amountSecond)
	transactions.Element(2).Object().ValueEqual("amount", amountFirst)

	carcass.ResetModels(userID)
}
