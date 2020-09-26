package transactions

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"balance/internal/handlers"
	"balance/internal/metrics"
	"balance/internal/models"
	"balance/internal/pkg/time"
	"balance/internal/services/core"
)

var (
	DefaultSort  = map[string]string{"time": "desc"}
	DefaultLimit = 10
	DefaultPage  = 1
)

func Controller(coreService *core.Service, stats metrics.Client) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer stats.NewTiming().Send(metrics.ControllersTransactionsTiming)
		stats.Increment(metrics.ControllersTransactionsCount)

		var in Request
		if err := handlers.Validate(context, &in, binding.Query); err != nil {
			context.JSON(http.StatusBadRequest, err.Response)
			return
		}

		transactions, pages, err := coreService.List(in.UserID, sort(in), limit(in), page(in))

		if err == core.ErrorUserNotFound {
			context.JSON(http.StatusBadRequest, handlers.ErrorBadRequest(
				handlers.ErrorCodeBalanceNotFound,
				fmt.Sprintf("Not found balance for user %d", in.UserID),
			))
			stats.Increment(metrics.Responses400AllCount)
			return
		}

		if err == core.ErrorInvalidSort {
			context.JSON(http.StatusBadRequest, handlers.ErrorBadRequest(
				fmt.Sprintf(handlers.ErrorCodeValidation, "sort"),
				"Invalid sort value",
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
			Time:       time.Format(tx.CreatedAt),
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
