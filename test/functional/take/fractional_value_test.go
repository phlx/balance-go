package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeFractionalValue(t *testing.T) {
	userID := int64(1)
	var amountTaken, amountExpected float64

	carcass := functional.NewCarcass(t)
	carcass.ResetModels(userID)

	amountStarted := 1.0

	carcass.API.GiveExpect(userID, amountStarted, nil, "").Status(http.StatusOK)

	amountTaken = 0.125
	amountExpected = 0.88

	carcass.API.TakeExpect(userID, amountTaken, nil, "").Status(http.StatusOK)

	carcass.API.BalanceExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amountExpected)

	carcass.ResetModels(userID)

	carcass.API.GiveExpect(userID, amountStarted, nil, "").Status(http.StatusOK)

	amountTaken = 0.135
	amountExpected = 0.86

	carcass.API.TakeExpect(userID, amountTaken, nil, "").Status(http.StatusOK)

	carcass.API.BalanceExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amountExpected)

	carcass.ResetModels(userID)
}
