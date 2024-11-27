package app

import (
	"banking/dto"
	"banking/mocks/service"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/wandz2810/banking-lib/errs"
	"net/http"
	"net/http/httptest"
	"testing"
)

var router *mux.Router

var mockService *service.MockCustomerService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockCustomerService(ctrl)
	ch := CustomerHandler{mockService}
	router = mux.NewRouter()
	router.HandleFunc("/customer", ch.getAllCustomer)
	return func() {
		router = nil
		defer ctrl.Finish()
	}
}

func Test_should_return_customers_with_status_code_200(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()
	dummyCustomers := []dto.CustomerResponse{
		{"1001", "Quang", "HCM", "2011", "2003-10-28", "1"},
		{"1002", "Huy", "HN", "2011", "2003-03-03", "1"},
	}
	mockService.EXPECT().GetAllCustomers("").Return(dummyCustomers, nil)
	request, _ := http.NewRequest(http.MethodGet, "/customer", nil)
	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	// Assert
	fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusOK {
		t.Errorf("Failed while testing the status code")
	}
}

func Test_should_return_status_code_500_with_error_message(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	mockService.EXPECT().GetAllCustomers("").Return(nil, errs.NewUnexpectedError("some database error"))
	request, _ := http.NewRequest(http.MethodGet, "/customer", nil)
	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	// Assert
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Failed while testing the status code")
	}
}
