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

	carcass := functional.NewCarcass(t)
	carcass.ResetModels(fromUserID)
	carcass.ResetModels(toUserID)

	amountStarted := 1000.0

	carcass.API.GiveExpect(fromUserID, amountStarted, nil, "").Status(http.StatusOK)

	amountTaken = 250.0
	amountExpectedFrom = 750.0
	amountExpectedTo = 250.0

	carcass.API.MoveExpect(fromUserID, toUserID, amountTaken, "").Status(http.StatusOK)

	carcass.API.BalanceExpect(fromUserID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amountExpectedFrom)

	carcass.API.BalanceExpect(toUserID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amountExpectedTo)

	carcass.ResetModels(fromUserID)
	carcass.ResetModels(toUserID)
}
