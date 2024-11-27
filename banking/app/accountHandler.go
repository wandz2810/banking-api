package app

import (
	"banking/dto"
	"banking/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type AccountHandler struct {
	service service.AccountService
}

func (h AccountHandler) NewAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customerId := vars["customer_id"]
	var request dto.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, err.Error(), http.StatusBadRequest)
	} else {
		request.CustomerId = customerId
		account, appError := h.service.CreateAccount(request)
		if appError != nil {
			writeResponse(w, appError.AsMessage(), appError.Code)
		} else {
			writeResponse(w, account, http.StatusCreated)
		}
	}
}

func (h AccountHandler) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountId := vars["account_id"]
	customerId := vars["customer_id"]
	var request dto.TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, err.Error(), http.StatusBadRequest)
	} else {
		request.AccountId = accountId
		request.CustomerId = customerId
		account, appError := h.service.MakeTransaction(request)
		if appError != nil {
			writeResponse(w, appError.Code, appError.Code)
		} else {
			writeResponse(w, account, http.StatusCreated)
		}
	}
}
