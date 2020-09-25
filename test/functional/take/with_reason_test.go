package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeWithReason(t *testing.T) {
	userID := int64(1)
	reason := "I'm нашёл твои money, Lebowski! 😈"

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)

	test.carcass.API.GiveExpect(userID, 100, nil, "").Status(http.StatusOK)

	test.carcass.API.TakeExpect(userID, 100, nil, reason).Status(http.StatusOK)

	test.carcass.API.TransactionsExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array().
		First().
		Object().
		ValueEqual("reason", reason)

	test.carcass.ResetModels(userID)
}
