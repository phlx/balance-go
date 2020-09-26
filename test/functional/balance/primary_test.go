package balance

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"balance/internal/exchangerates"
	"balance/internal/redis"
	"balance/test/functional"
)

func TestBalanceOnBucks(t *testing.T) {
	ctx := context.Background()

	carcass := functional.NewCarcass(t)

	userID := int64(1)
	_, _, err := carcass.App.CoreService.Give(userID, 100, nil, "")
	assert.Nil(t, err, "error on give 100 RUB to user ID %d", userID)

	carcass.App.Container.Redis = redis.NewStub()

	rates := map[exchangerates.Currency]float64{
		exchangerates.RUB: 1,
		exchangerates.USD: 0.03,
	}

	err = carcass.App.ExchangeRatesService.SetCachedRates(ctx, rates)
	assert.Nil(t, err, "error on set cached rates %+v", rates)

	result := carcass.Expectations.GET("/v1/balance").
		WithQuery("user_id", userID).
		WithQuery("currency", exchangerates.RUB).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	result.ValueEqual("currency", exchangerates.RUB)
	result.ValueEqual("balance", 100)

	result = carcass.Expectations.GET("/v1/balance").
		WithQuery("user_id", userID).
		WithQuery("currency", exchangerates.USD).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	result.ValueEqual("currency", exchangerates.USD)
	result.ValueEqual("balance", 3)

	carcass.ResetModels(userID)
}
