package app

import (
	"banking_auth/dto"
	"banking_auth/service"
	"encoding/json"
	"github.com/wandz2810/banking-lib/logger"
	"net/http"
)

type AuthHandler struct {
	service service.AuthService
}

func (h AuthHandler) NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, "Handler Not implemented...")
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logger.Error("Error while decoding login request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Login(loginRequest)
		if appErr != nil {
			writeResponse(w, appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, token)
		}
	}
}

// Sample URL string : http://localhost:8181/auth/verify?token=somevalidtokenstring&routeName=GetCustomer&customer_id=2000&account_id=12345
func (h AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	// convert from Query to map type
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}
	if urlParams["token"] != "" {

		appErr := h.service.Verify(urlParams)
		if appErr != nil {
			writeResponse(w, appErr.Code, notAuthorizedResponse(appErr.Message))
		} else {
			writeResponse(w, http.StatusOK, authorizedResponse())
		}
	} else {
		writeResponse(w, http.StatusForbidden, notAuthorizedResponse("missing token"))
	}

}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request dto.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		user, appError := h.service.CreateUser(request)
		if appError != nil {
			writeResponse(w, appError.Code, appError.AsMessage())
		} else {
			writeResponse(w, http.StatusCreated, user)
		}
	}
}
func notAuthorizedResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      msg,
	}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": true}

}
func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
