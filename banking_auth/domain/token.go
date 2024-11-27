package domain

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const HMAC_SAMPlE_SECRET = "hmacSampleSecret"
const ACCESS_TOKEN_DURATION = time.Hour

type Claims struct {
	CustomerId string   `json:"customer_id"`
	Accounts   []string `json:"account_numbers"`
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	jwt.StandardClaims
}

func (c Claims) IsUserRole() bool {
	return c.Role == "user"
}

func (c Claims) IsValidAccountId(accountId string) bool {
	if accountId != "" {
		accountFound := false
		for _, a := range c.Accounts {
			if a == accountId {
				accountFound = true
				break
			}
		}
		return accountFound
	}
	return true
}
func (c Claims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	if c.CustomerId != urlParams["customer_id"] {
		return false
	}

	if !c.IsValidAccountId(urlParams["account_id"]) {
		return false
	}
	return true
}
