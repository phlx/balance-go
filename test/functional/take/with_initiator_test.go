package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeWithInitiator(t *testing.T) {
	userID := int64(1)
	initiatorID := int64(2)

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)

	test.carcass.API.GiveExpect(userID, 100, nil, "").Status(http.StatusOK)

	test.carcass.API.TakeExpect(userID, 100, &initiatorID, "").Status(http.StatusOK)

	test.carcass.API.TransactionsExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array().
		First().
		Object().
		ValueEqual("from_user_id", initiatorID)

	test.carcass.ResetModels(userID)
}
