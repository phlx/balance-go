package transactions

import (
	"net/http"
	"testing"

	"balance/test/functional"
)

func TestTransactionsWithSortByNonexistentField(t *testing.T) {
	userID := int64(1)
	sort := map[string]string{
		"nonexistent": "asc",
	}
	carcass := functional.NewCarcass(t)

	carcass.API.TransactionsExpectWithSort(userID, sort).Status(http.StatusOK)
	// TODO: It must be BadRequest, but it's need to correct validation for request query structs
}
