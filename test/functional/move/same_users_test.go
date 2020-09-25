package move

import (
	"fmt"
	"net/http"
	"testing"

	"balance/internal/controllers"
	"balance/test/functional"
)

func TestMoveBetweenSameUsers(t *testing.T) {
	fromUserID := int64(1)
	toUserID := int64(1)

	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(fromUserID)

	test.carcass.API.GiveExpect(fromUserID, 100, nil, "").Status(http.StatusOK)

	test.carcass.API.MoveExpect(fromUserID, toUserID, 100, "").
		Status(http.StatusBadRequest).
		JSON().
		Object().
		Value("errors").
		Array().
		First().
		Object().
		ValueEqual("code", fmt.Sprintf(controllers.ErrorCodeValidation, "to_user_id"))

	test.carcass.ResetModels(fromUserID)
}
