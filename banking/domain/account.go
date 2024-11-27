package domain

import (
	"banking/dto"
	"github.com/wandz2810/banking-lib/errs"
	"time"
)

type Account struct {
	AccountId   string  `bson:"account_id"`
	CustomerId  string  `bson:"customer_id"`
	OpeningDate string  `bson:"opening_date"`
	AccountType string  `bson:"account_type"`
	Amount      float64 `bson:"amount"`
	Status      string  `bson:"status"`
}

const dbTSLayout = "2006-01-02 15:04:05"

var now = time.Now()

func (a Account) ToNewAccountResponseDto() *dto.CreateAccountResponse {
	return &dto.CreateAccountResponse{a.AccountId}
}

//go:generate mockgen -destination=../mocks/domain/mockAccountRepository.go -package=domain banking/domain AccountRepository
type AccountRepository interface {
	SaveAccount(Account) (*Account, *errs.AppError)
	SaveTransaction(Transaction) (*Transaction, *errs.AppError)
	FindByAccountId(string) (*Account, *errs.AppError)
	FindByCustomerIdAndAccountId(string, string) (*Account, *errs.AppError)
}

func (a Account) CanWithdraw(amount float64) bool {
	if a.Amount < amount {
		return false
	}
	return true
}

func NewAccount(customerId, accountType string, amount float64) Account {
	return Account{
		CustomerId:  customerId,
		OpeningDate: now.Format(dbTSLayout),
		AccountType: accountType,
		Amount:      amount,
		Status:      "1",
	}
}
