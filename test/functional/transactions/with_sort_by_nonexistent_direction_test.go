package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithSortByNonexistentDirection(t *testing.T) {
	userID := int64(1)
	sort := map[string]string{
		"id": "nonexistent",
	}
	carcass := functional.NewCarcass(t)

	carcass.API.TransactionsExpectWithSort(userID, sort).Status(http.StatusBadRequest)
}
