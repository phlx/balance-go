package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithNegativePage(t *testing.T) {
	userID := int64(1)
	page := -1
	limit := 10
	carcass := functional.NewCarcass(t)

	carcass.API.TransactionsExpectWithPageAndLimit(userID, limit, page).Status(http.StatusBadRequest)
}
