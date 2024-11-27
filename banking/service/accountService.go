package service

import (
	"banking/domain"
	"banking/dto"
	"github.com/wandz2810/banking-lib/errs"
)

const dbTSLayout = "2006-01-02 15:04:05"

type AccountService interface {
	CreateAccount(dto.CreateAccountRequest) (*dto.CreateAccountResponse, *errs.AppError)
	MakeTransaction(dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError)
}
type DefaultAccountService struct {
	repo domain.AccountRepository
}

func (s DefaultAccountService) CreateAccount(req dto.CreateAccountRequest) (*dto.CreateAccountResponse, *errs.AppError) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	account := domain.NewAccount(req.CustomerId, req.AccountType, req.Amount)
	if newAccount, err := s.repo.SaveAccount(account); err != nil {
		return nil, err
	} else {

		return newAccount.ToNewAccountResponseDto(), nil
	}
}

func (s DefaultAccountService) MakeTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	account, err := s.repo.FindByCustomerIdAndAccountId(req.CustomerId, req.AccountId)
	if err != nil {
		return nil, err
	}
	if req.IsTransactionTypeWithdrawal() {

		if !account.CanWithdraw(req.Amount) {
			return nil, errs.NewValidationError("insufficient balance")
		}
	}
	a := domain.Transaction{
		AccountId: req.AccountId,
		Amount:    req.Amount,
		//Balance:         account.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: dbTSLayout,
	}
	transaction, appError := s.repo.SaveTransaction(a)
	if appError != nil {
		return nil, appError
	}
	response := transaction.ToNewTransactionDto()
	return &response, nil

}

func NewAccountService(repository domain.AccountRepository) DefaultAccountService {

	return DefaultAccountService{repository}
}
