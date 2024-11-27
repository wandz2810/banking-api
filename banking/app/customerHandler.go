package app

import (
	"banking/dto"
	"banking/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type CustomerHandler struct {
	service service.CustomerService
}

func (ch *CustomerHandler) getAllCustomer(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	customers, err := ch.service.GetAllCustomers(status)
	if err != nil {
		writeResponse(w, err.AsMessage(), err.Code)
	} else {
		writeResponse(w, customers, http.StatusOK)
	}

}
func (ch *CustomerHandler) getCustomerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["customer_id"]

	customer, err := ch.service.GetCustomerById(id)
	if err != nil {
		writeResponse(w, err.AsMessage(), err.Code)
	} else {
		writeResponse(w, customer, http.StatusOK)
	}
}
func (ch CustomerHandler) UpdateCustomerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId := vars["customer_id"]
	var request dto.UpdateCustomerResquest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, err.Error(), http.StatusBadRequest)
	} else {
		request.CustomerId = customerId
		updateCustomer, appError := ch.service.UpdateCustomer(request)
		if appError != nil {
			writeResponse(w, appError.AsMessage(), appError.Code)
		} else {
			writeResponse(w, updateCustomer, http.StatusCreated)
		}
	}
}
func (ch *CustomerHandler) DeleteAccountById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["customer_id"]

	err := ch.service.DeleteUserById(id)
	if err != nil {
		writeResponse(w, err.AsMessage(), err.Code)
	} else {
		writeResponse(w, nil, http.StatusCreated)
	}
}
func writeResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
