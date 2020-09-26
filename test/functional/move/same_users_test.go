package move

import (
	"fmt"
	"net/http"
	"testing"

	"balance/internal/handlers"
	"balance/test/functional"
)

func TestMoveBetweenSameUsers(t *testing.T) {
	fromUserID := int64(1)
	toUserID := int64(1)

	carcass := functional.NewCarcass(t)
	carcass.ResetModels(fromUserID)

	carcass.API.GiveExpect(fromUserID, 100, nil, "").Status(http.StatusOK)

	carcass.API.MoveExpect(fromUserID, toUserID, 100, "").
		Status(http.StatusBadRequest).
		JSON().
		Object().
		Value("errors").
		Array().
		First().
		Object().
		ValueEqual("code", fmt.Sprintf(handlers.ErrorCodeValidation, "to_user_id"))

	carcass.ResetModels(fromUserID)
}
