package service

import (
	"banking/domain"
	"banking/dto"
	"github.com/wandz2810/banking-lib/errs"
)

//go:generate mockgen -destination=../mocks/service/mockCustomerService.go -package=service banking/service CustomerService
type CustomerService interface {
	GetAllCustomers(string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomerById(string) (*dto.CustomerResponse, *errs.AppError)
	UpdateCustomer(request dto.UpdateCustomerResquest) (*dto.UpdateCustomerResponse, *errs.AppError)
	DeleteUserById(string) *errs.AppError
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomers(status string) ([]dto.CustomerResponse, *errs.AppError) {
	if status == "active" {
		status = "1"
	} else if status == "inactive" {
		status = "0"
	} else {
		status = ""
	}

	customers, err := s.repo.FindAll(status)
	if err != nil {
		return nil, err
	}
	response := make([]dto.CustomerResponse, 0)
	for _, c := range customers {
		response = append(response, c.ToDto())
	}
	return response, nil
}

func (s DefaultCustomerService) GetCustomerById(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}

	response := c.ToDto()

	return &response, nil
}

func (s DefaultCustomerService) UpdateCustomer(req dto.UpdateCustomerResquest) (*dto.UpdateCustomerResponse, *errs.AppError) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	customer := domain.NewCustomer(req.CustomerId, req.Name, req.City, req.Zipcode, req.DateofBirth)
	if newAccount, err := s.repo.UpdateById(customer); err != nil {
		return nil, err
	} else {

		return newAccount.ToNewUpdateCustomerResponseDto(), nil
	}
}

func (s DefaultCustomerService) DeleteUserById(id string) *errs.AppError {
	err := s.repo.DeleteById(id)
	if err != nil {
		return err
	}

	return nil
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository}
}
