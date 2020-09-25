package balance

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"balance/internal/controllers"
	"balance/internal/exchangerates"
	"balance/internal/metrics"
	"balance/internal/services/core"
)

func Controller(coreService *core.Service, stats metrics.Client) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer stats.NewTiming().Send(metrics.ControllersBalanceTiming)
		stats.Increment(metrics.ControllersBalanceCount)

		var in Request
		if err := controllers.Validate(context, &in, binding.Query); err != nil {
			context.JSON(http.StatusBadRequest, err.Response)
			stats.Increment(metrics.Responses400AllCount)
			return
		}

		userID := in.UserID
		currency := exchangerates.Currency(in.Currency)
		if currency == "" {
			currency = exchangerates.RUB
		}
		result, err := coreService.Get(userID, currency)

		if err == core.ErrorUserNotFound {
			context.JSON(http.StatusBadRequest, controllers.ErrorBadRequest(
				controllers.ErrorCodeBalanceNotFound,
				fmt.Sprintf("Not found balance for user %d", in.UserID),
			))
			stats.Increment(metrics.Responses400AllCount)
			return
		}

		if err != nil {
			context.JSON(http.StatusInternalServerError, controllers.ErrorInternal())
			stats.Increment(metrics.Responses500AllCount)
			return
		}

		context.JSON(http.StatusOK, Response{
			Currency: currency,
			Balance:  result,
		})
		stats.Increment(metrics.Responses200AllCount)
	}
}
