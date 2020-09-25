package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithNegativeLimit(t *testing.T) {
	userID := int64(1)
	page := 1
	limit := -1
	carcass := functional.NewCarcass(t)

	carcass.API.TransactionsExpectWithPageAndLimit(userID, limit, page).Status(http.StatusBadRequest)
}
