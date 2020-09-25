package idempotence

import (
	"net/http"
	"testing"

	"github.com/google/uuid"

	"balance/test/functional"
)

func TestGiveMethodIsIdempotent(t *testing.T) {
	userID := int64(1)
	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)

	idempotencyKey := uuid.New().String()
	amount := 753.0
	given := test.carcass.API.GiveWithIdempotencyKeyExpect(idempotencyKey, userID, amount, nil, "").
		Status(http.StatusOK).
		JSON().
		Object()

	transaction := int64(given.Value("transaction").Number().Raw())
	time := given.Value("time").String().Raw()

	for attempts := 0; attempts < 5; attempts++ {
		test.carcass.API.GiveWithIdempotencyKeyExpect(idempotencyKey, userID, amount, nil, "").
			Status(http.StatusOK).
			JSON().
			Object().
			ValueEqual("transaction", transaction).
			ValueEqual("time", time)
	}

	test.carcass.API.BalanceExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amount)

	transactions := test.carcass.API.TransactionsExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array()

	transactions.Length().Equal(1)

	transactions.First().
		Object().
		ValueEqual("amount", amount).
		ValueEqual("time", time)

	test.carcass.ResetModels(userID)
}
