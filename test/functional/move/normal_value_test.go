package move

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestMoveNormalValue(t *testing.T) {
	fromUserID := int64(1)
	toUserID := int64(2)
	var amountTaken, amountExpectedFrom, amountExpectedTo float64

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(fromUserID)
	test.carcass.ResetModels(toUserID)

	amountStarted := 1000.0

	test.carcass.API.GiveExpect(fromUserID, amountStarted, nil, "").Status(http.StatusOK)

	amountTaken = 250.0
	amountExpectedFrom = 750.0
	amountExpectedTo = 250.0

	test.carcass.API.MoveExpect(fromUserID, toUserID, amountTaken, "").Status(http.StatusOK)

	test.carcass.API.BalanceExpect(fromUserID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amountExpectedFrom)

	test.carcass.API.BalanceExpect(toUserID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amountExpectedTo)

	test.carcass.ResetModels(fromUserID)
	test.carcass.ResetModels(toUserID)
}
