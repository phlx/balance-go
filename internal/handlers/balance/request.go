package balance

type Request struct {
	UserID   int64  `form:"user_id" binding:"required,numeric,gt=0"`
	Currency string `form:"currency" binding:"omitempty,currency"`
}
