package move

type Response struct {
	TransactionFromId   int64  `json:"transaction_from_id"`
	TransactionFromTime string `json:"transaction_from_time"`
	TransactionToId     int64  `json:"transaction_to_id"`
	TransactionToTime   string `json:"transaction_to_time"`
}
