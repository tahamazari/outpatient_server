package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/tahamazari/outpatient_server/api/db"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key")

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
)

// ExtractEmployeeIDFromToken extracts the employee ID from the Authorization header in the request
func ExtractEmployeeIDFromToken(c echo.Context) (int, error) {
	// Get the Authorization header value
	authHeader := c.Request().Header.Get(authorizationHeader)
	if authHeader == "" {
		return 0, errors.New("Authorization header is missing")
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return 0, errors.New("Invalid authorization format")
	}

	// Extract the token without the "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// TODO: Provide your JWT secret key here
		// This key should match the key used during token generation
		return []byte("your_secret_key"), nil
	})

	if err != nil {
		return 0, err
	}

	// Check if the token is valid
	if !token.Valid {
		return 0, errors.New("Invalid token")
	}

	// Extract employee ID from the token claims and convert to int
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Failed to extract claims from token")
	}

	employeeIDFloat, ok := claims["employee_id"].(float64)
	if !ok {
		return 0, errors.New("Employee ID not found or not a valid number in token claims")
	}

	employeeID := int(employeeIDFloat)

	return employeeID, nil
}

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
		"employee_id": employeeID,
		"exp":         time.Now().Add(time.Hour * 24 * 30).Unix(), // Refresh token expires in 30 days
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

// CheckPasswordHash compares a hashed password with its possible plaintext equivalent
func CheckPasswordHash(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

const ErrEmployeeWithEmailAlreadyExists = "Employee with this email/employeeId/certificateId already exists"
const ErrDependentNotFound = "Dependent not found or does not belong to the requesting employee"
const ErrBillingClaimNotFound = "Billing Claim not found or does not belong to the requesting employee"
const ErrMedicalBillNotFound = "Medical Bill not found or does not belong to the specified Billing Claim"
