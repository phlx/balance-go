package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithOverLimit(t *testing.T) {
	userID := int64(1)
	page := 1
	limit := 100100100
	carcass := functional.NewCarcass(t)

	carcass.API.TransactionsExpectWithPageAndLimit(userID, limit, page).Status(http.StatusBadRequest)
}
