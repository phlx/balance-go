package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/alexcesaro/statsd.v2"

	"balance/internal/metrics"
	"balance/internal/services/core"
)

func Controller(coreService *core.Service, stats *statsd.Client) gin.HandlerFunc {
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
			Time:     after,
			Latency:  duration,
		}
		context.JSON(http.StatusOK, response)
		stats.Increment(metrics.Responses200AllCount)
	}
}
