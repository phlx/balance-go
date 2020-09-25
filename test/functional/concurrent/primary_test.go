package concurrent

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"balance/internal/controllers"
	"balance/internal/controllers/take"
	"balance/internal/exchangerates"
	"balance/internal/middlewares"
	"balance/test/functional"
)

func TestConcurrent(t *testing.T) {
	carcass := functional.NewCarcass(t)

	userID := int64(1)
	_, _, err := carcass.App.CoreService.Give(userID, 1000, nil, "")
	assert.Nil(t, err, "error on give 1000 RUB to user ID %d", userID)

	result := carcass.Expectations.GET("/v1/balance").
		WithQuery("user_id", userID).
		WithQuery("currency", exchangerates.RUB).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	result.ValueEqual("currency", exchangerates.RUB)
	result.ValueEqual("balance", 1000)

	for i := 0; i < 10; i++ {
		time.Sleep(10 * time.Millisecond)

		expect := carcass.Expectations.POST("/v1/take").
			WithHeader(middlewares.IdempotencyHeader, uuid.New().String()).
			WithHeader("X-Sleep", "50").
			WithJSON(take.Request{
				UserID: 1,
				Amount: 500,
			}).
			Expect()

		if i < 2 {
			expect.Status(http.StatusOK)
		} else {
			expect.Status(http.StatusBadRequest).
				JSON().
				Object().
				Value("errors").
				Array().
				First().
				Object().
				ValueEqual("code", controllers.ErrorCodeInsufficientFunds)
		}
	}

	carcass.ResetModels(userID)
}
