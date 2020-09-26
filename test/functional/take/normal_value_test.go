package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeNormalValue(t *testing.T) {
	userID := int64(1)

	carcass := functional.NewCarcass(t)
	carcass.ResetModels(userID)

	amount := float64(100100100)
	taken := float64(100100)
	expected := float64(100000000)

	carcass.API.GiveExpect(userID, amount, nil, "").Status(http.StatusOK)

	carcass.API.TakeExpect(userID, taken, nil, "").Status(http.StatusOK)

	carcass.API.BalanceExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", expected)

	carcass.ResetModels(userID)
}
