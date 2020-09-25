package idempotence

import (
	"net/http"
	"testing"

	"github.com/google/uuid"

	"balance/test/functional"
)

func TestMoveMethodIsIdempotent(t *testing.T) {
	fromUserID := int64(1)
	toUserID := int64(2)
	test := &test{carcass: functional.NewCarcass(t)}
	test.carcass.ResetModels(fromUserID)
	test.carcass.ResetModels(toUserID)
	amountStarted := 789.0
	amountMoved := 321.0
	amountExpected := 468.0

	test.carcass.API.GiveExpect(fromUserID, amountStarted, nil, "").
		Status(http.StatusOK).
		JSON().
		Object()

	idempotencyKey := uuid.New().String()
	moved := test.carcass.API.MoveWithIdempotencyKeyExpect(idempotencyKey, fromUserID, toUserID, amountMoved, "").
		Status(http.StatusOK).
		JSON().
		Object()

	transactionFrom := int64(moved.Value("transaction_from_id").Number().Raw())
	timeFrom := moved.Value("transaction_from_time").String().Raw()
	transactionTo := int64(moved.Value("transaction_to_id").Number().Raw())
	timeTo := moved.Value("transaction_to_time").String().Raw()

	for attempts := 0; attempts < 5; attempts++ {
		test.carcass.API.MoveWithIdempotencyKeyExpect(idempotencyKey, fromUserID, toUserID, amountMoved, "").
			Status(http.StatusOK).
			JSON().
			Object().
			ValueEqual("transaction_from_id", transactionFrom).
			ValueEqual("transaction_from_time", timeFrom).
			ValueEqual("transaction_to_id", transactionTo).
			ValueEqual("transaction_to_time", timeTo)
	}

	test.carcass.API.BalanceExpect(fromUserID).
		Status(http.StatusOK).
		JSON().
		Object().
		ValueEqual("balance", amountExpected)

	transactionsFrom := test.carcass.API.TransactionsExpect(fromUserID).
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array()

	transactionsFrom.Length().Equal(2)

	transactionsFrom.First().
		Object().
		ValueEqual("amount", -amountMoved).
		ValueEqual("time", timeFrom)

	transactionsFrom.Last().Object().ValueEqual("amount", amountStarted)

	transactionsTo := test.carcass.API.TransactionsExpect(toUserID).
		Status(http.StatusOK).
		JSON().
		Object().
		Value("transactions").
		Array()

	transactionsTo.Length().Equal(1)

	transactionsTo.First().
		Object().
		ValueEqual("amount", amountMoved).
		ValueEqual("time", timeTo)

	test.carcass.ResetModels(fromUserID)
	test.carcass.ResetModels(toUserID)
}
