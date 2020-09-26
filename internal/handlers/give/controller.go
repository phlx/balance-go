package give

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"balance/internal/handlers"
	"balance/internal/metrics"
	"balance/internal/pkg/time"
	"balance/internal/services/core"
)

func Controller(coreService *core.Service, stats metrics.Client) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer stats.NewTiming().Send(metrics.ControllersGiveTiming)
		stats.Increment(metrics.ControllersGiveCount)

		var in Request
		if err := handlers.Validate(context, &in, binding.JSON); err != nil {
			context.JSON(http.StatusBadRequest, err.Response)
			stats.Increment(metrics.Responses400AllCount)
			return
		}

		_, tx, err := coreService.Give(in.UserID, in.Amount, in.InitiatorID, in.Reason)

		if err != nil {
			context.JSON(http.StatusInternalServerError, handlers.ErrorInternal())
			stats.Increment(metrics.Responses500AllCount)
			return
		}

		context.JSON(http.StatusOK, Response{
			Transaction: tx.Id,
			Time:        time.Format(tx.CreatedAt),
		})
		stats.Increment(metrics.Responses200AllCount)
	}
}
