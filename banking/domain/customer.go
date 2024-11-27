package domain

import (
	"banking/dto"
	"github.com/wandz2810/banking-lib/errs"
)

type Customer struct {
	Id          string `bson:"customer_id"`
	Name        string `bson:"name"`
	City        string `bson:"city"`
	Zipcode     string `bson:"zip_code"`
	DateofBirth string `bson:"date_of_birth"`
	Status      string `bson:"status"`
}

func (c Customer) statusAsText() string {
	statusAsText := "active"
	if c.Status == "0" {
		statusAsText = "inactive"
	}
	return statusAsText
}

func (c Customer) ToDto() dto.CustomerResponse {

	return dto.CustomerResponse{
		Id:          c.Id,
		Name:        c.Name,
		City:        c.City,
		Zipcode:     c.Zipcode,
		DateofBirth: c.DateofBirth,
		Status:      c.statusAsText(),
	}
}

func ToDeleteDto(response string) dto.DeleteUserResponse {
	return dto.DeleteUserResponse{response}
}

func (a Customer) ToNewUpdateCustomerResponseDto() *dto.UpdateCustomerResponse {
	return &dto.UpdateCustomerResponse{a.Id, a.Name, a.City, a.Zipcode, a.Status}
}

type CustomerRepository interface {
	// status == 1 status ==  status == ""
	FindAll(status string) ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
	UpdateById(customer Customer) (*Customer, *errs.AppError)
	DeleteById(string) *errs.AppError
}

func NewCustomer(Id string, Name string, City string, Zipcode string, DateofBirth string) Customer {
	return Customer{
		Id:          Id,
		Name:        Name,
		City:        City,
		Zipcode:     Zipcode,
		DateofBirth: DateofBirth,
		Status:      "1",
	}
}
