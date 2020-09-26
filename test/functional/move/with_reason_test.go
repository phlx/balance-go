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

	carcass := functional.NewCarcass(t)
	carcass.ResetModels(fromUserID)
	carcass.ResetModels(toUserID)

	carcass.API.GiveExpect(fromUserID, 100, nil, "").Status(http.StatusOK)

	carcass.API.MoveExpect(fromUserID, toUserID, 100, reason).Status(http.StatusOK)

	carcass.API.TransactionsExpect(fromUserID).Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array().
		First().
		Object().
		ValueEqual("reason", reason)

	carcass.API.TransactionsExpect(toUserID).Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array().
		First().
		Object().
		ValueEqual("reason", reason)

	carcass.ResetModels(fromUserID)
	carcass.ResetModels(toUserID)
}
