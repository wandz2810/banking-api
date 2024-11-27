package dto

type TransactionResponse struct {
	Amount        float64 `json:"amount"`
	TransactionId string  `json:"transaction_id"`
	Balance       float64 `json:"balance"`
}
