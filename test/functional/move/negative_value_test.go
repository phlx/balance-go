package move

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestMoveNegativeValue(t *testing.T) {
	fromUserID := int64(1)
	toUserID := int64(2)

	carcass := functional.NewCarcass(t)

	carcass.API.MoveExpect(fromUserID, toUserID, -0.01, "").Status(http.StatusBadRequest)
}
