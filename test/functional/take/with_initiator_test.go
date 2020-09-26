package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeWithInitiator(t *testing.T) {
	userID := int64(1)
	initiatorID := int64(2)

	carcass := functional.NewCarcass(t)
	carcass.ResetModels(userID)

	carcass.API.GiveExpect(userID, 100, nil, "").Status(http.StatusOK)

	carcass.API.TakeExpect(userID, 100, &initiatorID, "").Status(http.StatusOK)

	carcass.API.TransactionsExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array().
		First().
		Object().
		ValueEqual("from_user_id", initiatorID)

	carcass.ResetModels(userID)
}
