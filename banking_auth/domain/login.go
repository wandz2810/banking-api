package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/wandz2810/banking-lib/errs"
	"github.com/wandz2810/banking-lib/logger"
	"strings"
	"time"
)

const TOKEN_DURATION = 15 * time.Minute

type Login struct {
	Username   string `bson:"username"`
	CustomerId string `bson:"customer_id"`
	Accounts   string `bson:"account_numbers"`
	Role       string `bson:"role"`
}

func (l Login) GenerateToken() (*string, *errs.AppError) {
	var claims jwt.MapClaims
	if l.Accounts == "" && l.CustomerId == "" {
		claims = l.claimsForAdmin()
	} else {
		claims = l.claimsForUser()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenAsString, err := token.SignedString([]byte(HMAC_SAMPlE_SECRET))
	if err != nil {
		logger.Error("Failed while signing token:" + err.Error())
		return nil, errs.NewUnexpectedError("Cannot generate token")
	}
	return &signedTokenAsString, nil
}

func (l Login) claimsForUser() jwt.MapClaims {

	accounts := strings.Split(l.Accounts, ", ") // Adjust based on your concatenation logic

	return jwt.MapClaims{
		"customer_id":     l.CustomerId,
		"account_numbers": accounts,
		"username":        l.Username,
		"role":            l.Role,
		"exp":             time.Now().Add(TOKEN_DURATION).Unix(),
	}
}
func (l Login) claimsForAdmin() jwt.MapClaims {
	return jwt.MapClaims{
		"username": l.Username,
		"role":     l.Role,
		"exp":      time.Now().Add(TOKEN_DURATION).Unix(),
	}
}
