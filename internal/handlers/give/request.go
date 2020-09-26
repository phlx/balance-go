package give

type Request struct {
	UserID      int64   `json:"user_id" binding:"required,numeric,gt=0"`
	Amount      float64 `json:"amount" binding:"required,numeric,gte=0.01"`
	Reason      string  `json:"reason" binding:"omitempty"`
	InitiatorID *int64  `json:"initiator_id" binding:"omitempty,numeric"`
}
