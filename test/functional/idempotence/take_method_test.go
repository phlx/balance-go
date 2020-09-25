package idempotence

import (
	"net/http"
	"testing"

	"github.com/google/uuid"

	"balance/test/functional"
)

func TestTakeMethodIsIdempotent(t *testing.T) {
	userID := int64(1)
	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(userID)
	amountStarted := 789.0
	amountTaken := 321.0
	amountExpected := 468.0

	test.carcass.API.GiveExpect(userID, amountStarted, nil, "").
		Status(http.StatusOK).
		JSON().
		Object()

	idempotencyKey := uuid.New().String()
	taken := test.carcass.API.TakeWithIdempotencyKeyExpect(idempotencyKey, userID, amountTaken, nil, "").
		Status(http.StatusOK).
		JSON().
		Object()

	transaction := int64(taken.Value("transaction").Number().Raw())
	time := taken.Value("time").String().Raw()

	for attempts := 0; attempts < 5; attempts++ {
		test.carcass.API.TakeWithIdempotencyKeyExpect(idempotencyKey, userID, amountTaken, nil, "").
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
		ValueEqual("balance", amountExpected)

	transactions := test.carcass.API.TransactionsExpect(userID).
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array()

	transactions.Length().Equal(2)

	transactions.First().
		Object().
		ValueEqual("amount", -amountTaken).
		ValueEqual("time", time)

	test.carcass.ResetModels(userID)
}
