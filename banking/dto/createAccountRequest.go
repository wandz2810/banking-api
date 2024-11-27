package dto

import (
	"github.com/wandz2810/banking-lib/errs"
	"strings"
)

type CreateAccountRequest struct {
	CustomerId  string  `json:"customer_id"`
	AccountType string  `json:"account_type"`
	Amount      float64 `json:"amount"`
}

func (r CreateAccountRequest) Validate() *errs.AppError {
	if r.CustomerId == "" {
		return errs.NewValidationError("Customer id is required")
	}
	if r.Amount < 5000 {
		return errs.NewValidationError("To open a new account you need to deposit at least 5000")
	}
	if strings.ToLower(r.AccountType) != "saving" && r.AccountType != "checking" {
		return errs.NewValidationError("Account type should be either saving or checking")
	}
	return nil
}
