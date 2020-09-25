package move

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestMoveFromNonexistentUser(t *testing.T) {
	fromUserID := int64(3)
	toUserID := int64(2)

	test := &test{carcass: functional.NewCarcass(t)}

	test.carcass.API.MoveExpect(fromUserID, toUserID, 100, "").Status(http.StatusBadRequest)
}
