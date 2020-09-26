package take

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTakeNegativeValue(t *testing.T) {
	userID := int64(1)

	carcass := functional.NewCarcass(t)

	carcass.API.TakeExpect(userID, -0.01, nil, "").Status(http.StatusBadRequest)
}
