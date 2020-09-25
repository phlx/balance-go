package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeFractionalValue(t *testing.T) {
	userID := int64(1)
	var amountTaken, amountExpected float64

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)

	amountStarted := 1.0

	test.carcass.API.GiveExpect(userID, amountStarted, nil, "").Status(http.StatusOK)

	amountTaken = 0.125
	amountExpected = 0.88

	test.carcass.API.TakeExpect(userID, amountTaken, nil, "").Status(http.StatusOK)

	test.carcass.API.BalanceExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amountExpected)

	test.carcass.ResetModels(userID)

	test.carcass.API.GiveExpect(userID, amountStarted, nil, "").Status(http.StatusOK)

	amountTaken = 0.135
	amountExpected = 0.86

	test.carcass.API.TakeExpect(userID, amountTaken, nil, "").Status(http.StatusOK)

	test.carcass.API.BalanceExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amountExpected)

	test.carcass.ResetModels(userID)
}
