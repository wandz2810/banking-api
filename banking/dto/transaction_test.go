package dto

import (
	"net/http"
	"testing"
)

func Test_should_return_error_when_transaction_type_is_not_deposit_or_withdrawal(t *testing.T) {
	// Arrange
	request := TransactionRequest{
		TransactionType: "invalid transaction type",
	}
	// Act
	appError := request.Validate()
	// Assert
	if appError.Message != "Transaction type must be either 'withdrawal' or 'deposit'" {
		t.Error("Invalid message while testing transaction type")
	}

	if appError.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid status code while testing transaction type")
	}
}

func Test_should_return_error_when_ammount_is_less_than_zero(t *testing.T) {
	// Arrange
	request := TransactionRequest{
		TransactionType: DEPOSIT,
		Amount:          -100,
	}
	// Act
	appError := request.Validate()

	// Assert
	if appError.Message != "Amount must be greater than 0" {
		t.Error("Invalid message while testing transaction type")
	}

	if appError.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid status code while testing transaction type")
	}
}
