package domain

import "banking_auth/dto"

type Register struct {
	CustomerId  string `bson:"customer_id"`
	Username    string `bson:"username"`
	Password    string `bson:"password"`
	Role        string `bson:"role"`
	Name        string `bson:"name"`
	City        string `bson:"city"`
	Zipcode     string `bson:"zip_code"`
	DateofBirth string `bson:"date_of_birth"`
	Status      string `bson:"status"`
}

func (r Register) ToNewRegisterResponseDto() *dto.RegisterResponse {
	return &dto.RegisterResponse{r.CustomerId}
}

func NewRegister(username string, password string, name string, city string, zipcode string, dateofbirth string) Register {
	return Register{
		Username:    username,
		Password:    password,
		Role:        "user",
		Name:        name,
		City:        city,
		Zipcode:     zipcode,
		DateofBirth: dateofbirth,
		Status:      "1",
	}
}
