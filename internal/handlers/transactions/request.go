package transactions

type Request struct {
	UserID     int64  `form:"user_id" binding:"required,numeric,gt=0"`
	SortTime   string `form:"sort[time]" binding:"omitempty,oneof=asc desc"`
	SortID     string `form:"sort[id]" binding:"omitempty,oneof=asc desc"`
	SortAmount string `form:"sort[amount]" binding:"omitempty,oneof=asc desc"`
	Limit      int    `form:"limit" binding:"omitempty,numeric,min=10,max=100"`
	Page       int    `form:"page" binding:"omitempty,numeric,min=1"`
}
