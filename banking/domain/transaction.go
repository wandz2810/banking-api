package domain

import "banking/dto"

const WITHDRAWAL = "withdrawal"

type Transaction struct {
	TransactionId   string  `bson:"transaction_id""`
	AccountId       string  `bson:"account_id"`
	Amount          float64 `bson:"amount"`
	Balance         float64 `bson:"balance"`
	TransactionType string  `bson:"transaction_type"`
	TransactionDate string  `bson:"transaction_date"`
}

func (t Transaction) ToNewTransactionDto() dto.TransactionResponse {
	return dto.TransactionResponse{t.Amount, t.TransactionId, t.Balance}
}

func (t Transaction) IsWithdrawal() bool {
	if t.TransactionType == WITHDRAWAL {
		return true
	} else {
		return false
	}
}

type TransactionRepository interface {
}
