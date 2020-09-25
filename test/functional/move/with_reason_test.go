package move

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestMoveWithReason(t *testing.T) {
	fromUserID := int64(1)
	toUserID := int64(2)
	reason := "Move like Jaeger 🕺"

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(fromUserID)
	test.carcass.ResetModels(toUserID)

	test.carcass.API.GiveExpect(fromUserID, 100, nil, "").Status(http.StatusOK)

	test.carcass.API.MoveExpect(fromUserID, toUserID, 100, reason).Status(http.StatusOK)

	test.carcass.API.TransactionsExpect(fromUserID).Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array().
		First().
		Object().
		ValueEqual("reason", reason)

	test.carcass.API.TransactionsExpect(toUserID).Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array().
		First().
		Object().
		ValueEqual("reason", reason)

	test.carcass.ResetModels(fromUserID)
	test.carcass.ResetModels(toUserID)
}
