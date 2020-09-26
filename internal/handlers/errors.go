package handlers

const (
	ErrorCodeInternal          = "internal_error"
	ErrorCodeValidation        = "validation_error:%s"
	ErrorCodeBalanceNotFound   = "balance_not_found"
	ErrorCodeUserNotFound      = "user_not_found"
	ErrorCodeInsufficientFunds = "insufficient_funds"
)

type ErrorResponseItem struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type ErrorResponse struct {
	Errors []ErrorResponseItem `json:"errors"`
}

func ErrorBadRequest(code string, detail string) ErrorResponse {
	return ErrorResponse{
		Errors: []ErrorResponseItem{
			{
				Code:   code,
				Detail: detail,
			},
		},
	}
}

func ErrorInternal() ErrorResponse {
	return ErrorResponse{
		Errors: []ErrorResponseItem{
			{
				Code:   ErrorCodeInternal,
				Detail: "Internal Server Error",
			},
		},
	}
}
