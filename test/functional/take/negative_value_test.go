package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeNegativeValue(t *testing.T) {
	userID := int64(1)

	test := &test{carcass: functional.NewCarcass(t)}

	test.carcass.API.TakeExpect(userID, -0.01, nil, "").Status(http.StatusBadRequest)
}
