package move

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestMoveNegativeValue(t *testing.T) {
	fromUserID := int64(1)
	toUserID := int64(2)

	test := &test{carcass: functional.NewCarcass(t)}

	test.carcass.API.MoveExpect(fromUserID, toUserID, -0.01, "").Status(http.StatusBadRequest)
}
