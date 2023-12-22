package utils

import (
	"fmt"

	"github.com/tahamazari/outpatient_server/api/db"
	"golang.org/x/crypto/bcrypt"
)

type ErrorMessageArgs struct {
	ErrorMessage string
	Model        interface{}
	Query        string
	Args         []interface{}
}

func CheckRecordExistence(errorMessage string, model interface{}, query string, args []interface{}) error {
	db := db.DB()

	if err := db.Where(query, args...).First(model).Error; err != nil {
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func HashPassword(password string) (string, error) {
	// Use a secure password hashing library like bcrypt
	// Example using bcrypt:
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

const ErrEmployeeWithEmailAlreadyExists = "User with this email already exists"
const ErrDependentNotFound = "Dependent not found or does not belong to the requesting employee"
const ErrBillingClaimNotFound = "Billing Claim not found or does not belong to the requesting employee"
const ErrMedicalBillNotFound = "Medical Bill not found or does not belong to the specified Billing Claim"
