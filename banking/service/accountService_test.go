package service

import (
	realdomain "banking/domain"
	"banking/dto"
	"banking/mocks/domain"
	"github.com/golang/mock/gomock"
	"github.com/wandz2810/banking-lib/errs"
	"testing"
	"time"
)

var mockRepo *domain.MockAccountRepository
var service AccountService
var now = time.Now()

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo = domain.NewMockAccountRepository(ctrl)
	service = NewAccountService(mockRepo)
	return func() {
		service = nil
		defer ctrl.Finish()
	}

}

func Test_should_return_a_validation_error_response_when_the_service_is_not_validated(t *testing.T) {
	//Arrange
	request := dto.CreateAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      500,
	}
	service := NewAccountService(nil)
	//Act
	_, appError := service.CreateAccount(request)
	//Assert
	if appError == nil {
		t.Error("failed while testing the new account validation")
	}
}

func Test_should_return_an_error_from_the_server_side_if_the_new_account_cannot_be_created(t *testing.T) {
	//Arrange
	teardown := setup(t)
	defer teardown()

	req := dto.CreateAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      500,
	}
	account := realdomain.Account{
		CustomerId:  req.CustomerId,
		OpeningDate: dbTSLayout,
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	mockRepo.EXPECT().SaveAccount(account).Return(nil, errs.NewUnexpectedError("Unexpected database error"))
	//Act
	_, appError := service.CreateAccount(req)
	//Assert
	if appError == nil {
		t.Error("Test failed while validate error for the new account ")
	}
}

func Test_should_return_new_account_response_when_a_new_account_is_save_successfully(t *testing.T) {
	//Arrange
	teardown := setup(t)
	defer teardown()

	req := dto.CreateAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      6000,
	}
	account := realdomain.Account{
		CustomerId:  req.CustomerId,
		OpeningDate: now.Format(dbTSLayout),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	accountWithId := account
	accountWithId.AccountId = "2111"
	mockRepo.EXPECT().SaveAccount(account).Return(&accountWithId, nil)
	//Act
	newAccount, appError := service.CreateAccount(req)
	//Assert
	if appError != nil {
		t.Error("Test failed while creating new account ")
	}
	if newAccount.AccountId != accountWithId.AccountId {
		t.Error("Test failed matching new account ")
	}
}
