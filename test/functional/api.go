package functional

import (
	"fmt"
	"reflect"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"

	"balance/internal/controllers/give"
	"balance/internal/controllers/move"
	"balance/internal/controllers/take"
	"balance/internal/middlewares"
)

type API struct {
	expectations *httpexpect.Expect
}

func (a *API) BalanceExpect(userID int64) *httpexpect.Response {
	return a.expectations.GET("/v1/balance").
		WithQuery("user_id", userID).
		Expect()
}

func (a *API) TransactionsExpect(userID int64) *httpexpect.Response {
	return a.expectations.GET("/v1/transactions").
		WithQuery("user_id", userID).
		Expect()
}

func (a *API) TransactionsExpectWithPageAndLimit(userID int64, limit, page int) *httpexpect.Response {
	return a.expectations.GET("/v1/transactions").
		WithQuery("user_id", userID).
		WithQuery("limit", limit).
		WithQuery("page", page).
		Expect()
}

func (a *API) TransactionsExpectWithSort(userID int64, sort map[string]string) *httpexpect.Response {
	keys := reflect.ValueOf(sort).MapKeys()
	key := keys[0].String()
	value := sort[key]

	return a.expectations.GET("/v1/transactions").
		WithQuery("user_id", userID).
		WithQuery(fmt.Sprintf("sort[%s]", key), value).
		Expect()
}

func (a *API) GiveExpect(userID int64, amount float64, initiatorID *int64, reason string) *httpexpect.Response {
	idempotencyKey := uuid.New().String()
	return a.GiveWithIdempotencyKeyExpect(idempotencyKey, userID, amount, initiatorID, reason)
}

func (a *API) GiveWithIdempotencyKeyExpect(
	idempotencyKey string,
	userID int64,
	amount float64,
	initiatorID *int64,
	reason string,
) *httpexpect.Response {
	return a.expectations.POST("/v1/give").
		WithHeader(middlewares.IdempotencyHeader, idempotencyKey).
		WithJSON(give.Request{
			UserID:      userID,
			Amount:      amount,
			InitiatorID: initiatorID,
			Reason:      reason,
		}).
		Expect()
}

func (a *API) TakeExpect(userID int64, amount float64, initiatorID *int64, reason string) *httpexpect.Response {
	idempotencyKey := uuid.New().String()
	return a.TakeWithIdempotencyKeyExpect(idempotencyKey, userID, amount, initiatorID, reason)
}

func (a *API) TakeWithIdempotencyKeyExpect(
	idempotencyKey string,
	userID int64,
	amount float64,
	initiatorID *int64,
	reason string,
) *httpexpect.Response {
	return a.expectations.POST("/v1/take").
		WithHeader(middlewares.IdempotencyHeader, idempotencyKey).
		WithJSON(take.Request{
			UserID:      userID,
			Amount:      amount,
			InitiatorID: initiatorID,
			Reason:      reason,
		}).
		Expect()
}

func (a *API) MoveExpect(fromUserID, toUserID int64, amount float64, reason string) *httpexpect.Response {
	idempotencyKey := uuid.New().String()
	return a.MoveWithIdempotencyKeyExpect(idempotencyKey, fromUserID, toUserID, amount, reason)
}

func (a *API) MoveWithIdempotencyKeyExpect(
	idempotencyKey string,
	fromUserID int64,
	toUserID int64,
	amount float64,
	reason string,
) *httpexpect.Response {
	return a.expectations.POST("/v1/move").
		WithHeader(middlewares.IdempotencyHeader, idempotencyKey).
		WithJSON(move.Request{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Amount:     amount,
			Reason:     reason,
		}).
		Expect()
}
