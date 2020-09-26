package move

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"balance/internal/handlers"
	"balance/internal/metrics"
	"balance/internal/pkg/time"
	"balance/internal/services/core"
)

func Controller(coreService *core.Service, stats metrics.Client) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer stats.NewTiming().Send(metrics.ControllersMoveTiming)
		stats.Increment(metrics.ControllersMoveCount)

		var in Request
		if err := handlers.Validate(context, &in, binding.JSON); err != nil {
			context.JSON(http.StatusBadRequest, err.Response)
			stats.Increment(metrics.Responses400AllCount)
			return
		}

		if in.FromUserID == in.ToUserID {
			context.JSON(http.StatusBadRequest, handlers.ErrorBadRequest(
				fmt.Sprintf(handlers.ErrorCodeValidation, "to_user_id"),
				"Move allowed only between different users",
			))
			stats.Increment(metrics.Responses400AllCount)
			return
		}

		sleep := sleep(context.GetHeader("X-Sleep"))

		txFrom, txTo, err := coreService.Move(in.FromUserID, in.ToUserID, in.Amount, in.Reason, sleep)

		if err == core.ErrorUserNotFound {
			context.JSON(http.StatusBadRequest, handlers.ErrorBadRequest(
				handlers.ErrorCodeUserNotFound,
				"User was not found",
			))
			stats.Increment(metrics.Responses400AllCount)
			return
		}

		if err == core.ErrorInsufficientFunds {
			context.JSON(http.StatusBadRequest, handlers.ErrorBadRequest(
				handlers.ErrorCodeInsufficientFunds,
				"User has insufficient funds",
			))
			stats.Increment(metrics.Responses400AllCount)
			return
		}

		if err != nil {
			context.JSON(http.StatusInternalServerError, handlers.ErrorInternal())
			stats.Increment(metrics.Responses500AllCount)
			return
		}

		context.JSON(http.StatusOK, Response{
			TransactionFromId:   txFrom.Id,
			TransactionFromTime: time.Format(txFrom.CreatedAt),
			TransactionToId:     txTo.Id,
			TransactionToTime:   time.Format(txTo.CreatedAt),
		})
		stats.Increment(metrics.Responses200AllCount)
	}
}

func sleep(source string) int64 {
	var sleep int64
	if source != "" {
		if parsed, e := strconv.ParseInt(source, 10, 64); e == nil {
			sleep = parsed
		}
	}
	return sleep
}
