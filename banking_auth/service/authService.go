package service

import (
	"banking_auth/domain"
	"banking_auth/dto"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/wandz2810/banking-lib/errs"
	"github.com/wandz2810/banking-lib/logger"
)

type AuthService interface {
	Login(request dto.LoginRequest) (*string, *errs.AppError)
	//Verify(urlParams map[string]string) error
	Verify(urlParams map[string]string) *errs.AppError
	CreateUser(request dto.RegisterRequest) (*dto.RegisterResponse, *errs.AppError)
}
type DefaultAuthService struct {
	repo            domain.AuthRepositoty
	rolePermissions domain.RolePermissions
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*string, *errs.AppError) {
	login, err := s.repo.FindBy(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := login.GenerateToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	// convert the string token to JWT struct
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.NewAuthorizationError(err.Error())
	} else {
		/*
		   Checking the validity of the token, this verifies the expiry
		   time and the signature of the token
		*/
		if jwtToken.Valid {
			// type cast the token claims to jwt.MapClaims
			claims := jwtToken.Claims.(*domain.Claims)
			/* if Role if user then check if the account_id and customer_id
			   coming in the URL belongs to the same token
			*/
			if claims.IsUserRole() {
				if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
					return errs.NewAuthorizationError("request not verified with the token claims")
				}
			}
			// verify of the role is authorized to use the route
			isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
			if !isAuthorized {
				return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Role))
			}
			return nil
		} else {
			return errs.NewAuthorizationError("Invalid token")
		}
	}
}

func (s DefaultAuthService) CreateUser(req dto.RegisterRequest) (*dto.RegisterResponse, *errs.AppError) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	user := domain.NewRegister(req.Username, req.Password, req.Name, req.City, req.Zipcode, req.DateofBirth)
	if newUser, err := s.repo.CreateUser(user); err != nil {
		return nil, err
	} else {
		return newUser.ToNewRegisterResponseDto(), nil
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPlE_SECRET), nil
	})
	if err != nil {
		logger.Error("Error while parsing token" + err.Error())
		return nil, err
	}
	return token, nil
}

func NewAuthService(repository domain.AuthRepositoty, permissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{repository, permissions}
}
