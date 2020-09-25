package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithNonexistentUser(t *testing.T) {
	userID := int64(1)
	carcass := functional.NewCarcass(t)

	transactions := carcass.API.TransactionsExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object()

	transactions.Value("transactions").Array().Length().Equal(0)
	transactions.Value("pages").Equal(0)
}
