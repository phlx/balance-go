package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"balance/internal/metrics"
	"balance/internal/pkg/time"
	"balance/internal/services/core"
)

func Controller(coreService *core.Service, stats metrics.Client) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer stats.NewTiming().Send(metrics.ControllersHealthTiming)
		stats.Increment(metrics.ControllersHealthCount)

		before := time.Now()
		health := coreService.Health()
		after := time.Now()
		duration := after.Sub(before).Milliseconds()
		response := Response{
			Postgres: health.Postgres,
			Redis:    health.Redis,
			Errors:   health.Errors,
			Time:     time.Format(after),
			Latency:  duration,
		}
		context.JSON(http.StatusOK, response)
		stats.Increment(metrics.Responses200AllCount)
	}
}
