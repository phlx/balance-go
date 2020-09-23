package transactions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/alexcesaro/statsd.v2"

	"balance/internal/controllers"
	"balance/internal/metrics"
	"balance/internal/models"
	"balance/internal/services/core"
)

var (
	DefaultSort  = map[string]string{"time": "desc"}
	DefaultLimit = 10
	DefaultPage  = 1
)

func Controller(coreService *core.Service, stats *statsd.Client) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer stats.NewTiming().Send(metrics.ControllersTransactionsTiming)
		stats.Increment(metrics.ControllersTransactionsCount)

		var in Request
		if err := controllers.Validate(context, &in, binding.Query); err != nil {
			context.JSON(http.StatusBadRequest, err.Response)
			return
		}

		transactions, pages, err := coreService.List(in.UserID, sort(in), limit(in), page(in))

		if err == core.ErrorUserNotFound {
			context.JSON(http.StatusBadRequest, controllers.ErrorBadRequest(
				controllers.ErrorCodeBalanceNotFound,
				fmt.Sprintf("Not found balance for user %d", in.UserID),
			))
			stats.Increment(metrics.Responses400AllCount)
			return
		}

		if err == core.ErrorInvalidSort {
			context.JSON(http.StatusBadRequest, controllers.ErrorBadRequest(
				fmt.Sprintf(controllers.ErrorCodeValidation, "sort"),
				"Invalid sort value",
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
			Transactions: mapTransactions(transactions),
			Pages:        pages,
		})
		stats.Increment(metrics.Responses200AllCount)
	}
}

func mapTransactions(transactions []models.Transaction) []Transaction {
	result := make([]Transaction, 0)

	for _, tx := range transactions {
		result = append(result, Transaction{
			Id:         tx.Id,
			Time:       tx.CreatedAt.Format(time.RFC3339Nano),
			Amount:     tx.Amount,
			FromUserId: tx.InitiatorId,
			Reason:     tx.Reason,
		})
	}

	return result
}

func sort(in Request) map[string]string {
	if in.SortTime != "" {
		return map[string]string{"time": in.SortTime}
	}
	if in.SortID != "" {
		return map[string]string{"id": in.SortID}
	}
	if in.SortAmount != "" {
		return map[string]string{"amount": in.SortAmount}
	}
	return DefaultSort
}

func limit(in Request) int {
	if in.Limit != 0 {
		return in.Limit
	}
	return DefaultLimit
}

func page(in Request) int {
	if in.Page != 0 {
		return in.Page
	}
	return DefaultPage
}
