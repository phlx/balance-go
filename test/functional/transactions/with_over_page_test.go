package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithOverPage(t *testing.T) {
	userID := int64(1)
	limit := 10
	page := 100100100
	carcass := functional.NewCarcass(t)

	transactions := carcass.API.TransactionsExpectWithPageAndLimit(userID, limit, page).
		Status(http.StatusOK).
		JSON().
		Object()

	transactions.Value("transactions").Array().Length().Equal(0)
	transactions.Value("pages").Equal(0)
}
