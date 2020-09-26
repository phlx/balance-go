package move

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestMoveFromNonexistentUser(t *testing.T) {
	fromUserID := int64(3)
	toUserID := int64(2)

	carcass := functional.NewCarcass(t)

	carcass.API.MoveExpect(fromUserID, toUserID, 100, "").Status(http.StatusBadRequest)
}
