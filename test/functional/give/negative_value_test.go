package give

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestGiveNegativeValue(t *testing.T) {
	userID := int64(1)

	test := &test{carcass: functional.NewCarcass(t)}

	test.carcass.API.GiveExpect(userID, -0.01, nil, "").Status(http.StatusBadRequest)
}
