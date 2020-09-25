package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeNormalValue(t *testing.T) {
	userID := int64(1)

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)

	amount := float64(100100100)
	taken := float64(100100)
	expected := float64(100000000)

	test.carcass.API.GiveExpect(userID, amount, nil, "").Status(http.StatusOK)

	test.carcass.API.TakeExpect(userID, taken, nil, "").Status(http.StatusOK)

	test.carcass.API.BalanceExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", expected)

	test.carcass.ResetModels(userID)
}
