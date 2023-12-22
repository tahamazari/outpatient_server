package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tahamazari/outpatient_server/api/db"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key")

func CheckRecordExistence(errorMessage string, model interface{}, query string, args []interface{}) error {
	db := db.DB()

	if err := db.Where(query, args...).First(model).Error; err != nil {
		return fmt.Errorf(errorMessage)
	}

	return nil
}

func GenerateJWT(employeeID int) (string, error) {
	claims := jwt.MapClaims{
		"employee_id": employeeID,
		"exp":         time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(employeeID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id":     employeeID,
		"employee_id": time.Now().Add(time.Hour * 24 * 30).Unix(), // Refresh token expires in 30 days
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
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

const ErrEmployeeWithEmailAlreadyExists = "Employee with this email/employeeId/certificateId already exists"
const ErrDependentNotFound = "Dependent not found or does not belong to the requesting employee"
const ErrBillingClaimNotFound = "Billing Claim not found or does not belong to the requesting employee"
const ErrMedicalBillNotFound = "Medical Bill not found or does not belong to the specified Billing Claim"
