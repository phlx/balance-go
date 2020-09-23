package move

type Request struct {
	FromUserID int64   `json:"from_user_id" binding:"required,numeric,gt=0"`
	ToUserID   int64   `json:"to_user_id" binding:"required,numeric,gt=0"`
	Amount     float64 `json:"amount" binding:"required,numeric,gte=0.01"`
	Reason     string  `json:"reason" binding:"omitempty"`
}
