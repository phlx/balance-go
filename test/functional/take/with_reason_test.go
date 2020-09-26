package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeWithReason(t *testing.T) {
	userID := int64(1)
	reason := "I'm нашёл твои money, Lebowski! 😈"

	carcass := functional.NewCarcass(t)
	carcass.ResetModels(userID)

	carcass.API.GiveExpect(userID, 100, nil, "").Status(http.StatusOK)

	carcass.API.TakeExpect(userID, 100, nil, reason).Status(http.StatusOK)

	carcass.API.TransactionsExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array().
		First().
		Object().
		ValueEqual("reason", reason)

	carcass.ResetModels(userID)
}
