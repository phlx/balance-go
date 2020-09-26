package transactions

type Transaction struct {
	Id         int64   `json:"id"`
	Time       string  `json:"time"`
	Amount     float64 `json:"amount"`
	FromUserId *int64  `json:"from_user_id"`
	Reason     string  `json:"reason"`
}

type Response struct {
	Transactions []Transaction `json:"transactions"`
	Pages        int           `json:"pages"`
}
