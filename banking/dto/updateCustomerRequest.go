package dto

import (
	"github.com/wandz2810/banking-lib/errs"
	"regexp"
	"strings"
	"time"
)

type UpdateCustomerResquest struct {
	CustomerId  string `json:"customer_id"`
	Name        string `json:"full_name"`
	City        string `json:"city"`
	Zipcode     string `json:"zipcode"`
	DateofBirth string `json:"date_of_Birth"`
}

const dbTSLayout = "2006-01-02"

func (r UpdateCustomerResquest) Validate() *errs.AppError {

	if r.Name == "" {
		return errs.NewValidationError("Name is required")
	}
	if r.City == "" {
		return errs.NewValidationError("City is required")
	}
	if r.Zipcode != "" {
		// Loại bỏ các khoảng trắng nếu có
		zipcode := strings.TrimSpace(r.Zipcode)
		// Kiểm tra mã bưu điện có đúng 5 chữ số hay không
		if match, _ := regexp.MatchString(`^\d{5}$`, zipcode); !match {
			return errs.NewValidationError("zipcode must be a 5-digit number")
		}
	}
	if r.DateofBirth != "" {
		// Kiểm tra định dạng ngày tháng
		if _, err := time.Parse(dbTSLayout, r.DateofBirth); err != nil {
			return errs.NewValidationError("date of birth must be in YYYY-MM-DD format")
		}
	}
	return nil
}
