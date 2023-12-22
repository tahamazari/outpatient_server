package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tahamazari/outpatient_server/api/db"
	"github.com/tahamazari/outpatient_server/api/models"
	"github.com/tahamazari/outpatient_server/utils"
)

type EmployeeResponse struct {
	ID            int    `json:"id"`
	CompanyID     string `json:"company_id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	ContactNumber string `json:"contact_number"`
	CertificateID string `json:"certificate_id"`
}

func SignUp(c echo.Context) error {
	requestSignupEmployee := new(models.Employee)
	db := db.DB()

	// Binding data
	if err := c.Bind(requestSignupEmployee); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	// Check for missing required fields
	if requestSignupEmployee.Name == "" || requestSignupEmployee.Email == "" || requestSignupEmployee.CompanyID == "" ||
		requestSignupEmployee.ContactNumber == "" || requestSignupEmployee.CertificateID == "" || requestSignupEmployee.Password == "" {
		data := map[string]interface{}{
			"message": "Missing required fields",
		}
		return c.JSON(http.StatusBadRequest, data)
	}

	var existingEmployee models.Employee
	if err := db.Where(
		"email = ? OR company_id = ? OR certificate_id = ?",
		requestSignupEmployee.Email,
		requestSignupEmployee.CompanyID,
		requestSignupEmployee.CertificateID,
	).First(&existingEmployee).Error; err == nil {
		data := map[string]interface{}{
			"message": utils.ErrEmployeeWithEmailAlreadyExists,
		}
		return c.JSON(http.StatusConflict, data)
	}

	hashedPassword, err := utils.HashPassword(requestSignupEmployee.Password)
	if err != nil {
		data := map[string]interface{}{
			"message": "Error hashing the password",
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	newEmployee := &models.Employee{
		Name:          requestSignupEmployee.Name,
		Email:         requestSignupEmployee.Email,
		Password:      hashedPassword,
		CompanyID:     requestSignupEmployee.CompanyID,
		ContactNumber: requestSignupEmployee.ContactNumber,
		CertificateID: requestSignupEmployee.CertificateID,
	}

	if err := db.Create(&newEmployee).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	accessToken, err := utils.GenerateJWT(newEmployee.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Error generating access token"})
	}

	// Create and sign the refresh token
	refreshToken, err := utils.GenerateRefreshToken(newEmployee.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Error generating refresh token"})
	}

	// Create response without the password field
	employeeResponse := EmployeeResponse{
		ID:            newEmployee.ID,
		CompanyID:     newEmployee.CompanyID,
		Name:          newEmployee.Name,
		Email:         newEmployee.Email,
		ContactNumber: newEmployee.ContactNumber,
		CertificateID: newEmployee.CertificateID,
	}

	// Include tokens in the response
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"employee": employeeResponse,
			"tokens": map[string]interface{}{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}

func Login(c echo.Context) error {
	// Assuming the login request is sent as JSON in the request body
	loginRequest := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	if err := c.Bind(loginRequest); err != nil {
		data := map[string]interface{}{
			"message": "Invalid login request",
		}
		return c.JSON(http.StatusBadRequest, data)
	}

	db := db.DB()

	// Find the user by email
	var employee models.Employee
	if err := db.Where("email = ?", loginRequest.Email).First(&employee).Error; err != nil {
		data := map[string]interface{}{
			"message": "Invalid email or password",
		}
		return c.JSON(http.StatusUnauthorized, data)
	}

	// Verify the password
	if !utils.CheckPasswordHash(loginRequest.Password, employee.Password) {
		data := map[string]interface{}{
			"message": "Invalid email or password",
		}
		return c.JSON(http.StatusUnauthorized, data)
	}

	// Generate new access token
	accessToken, err := utils.GenerateJWT(employee.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Error generating access token"})
	}

	// Generate new refresh token
	refreshToken, err := utils.GenerateRefreshToken(employee.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Error generating refresh token"})
	}

	// Create response without the password field
	employeeResponse := EmployeeResponse{
		ID:            employee.ID,
		CompanyID:     employee.CompanyID,
		Name:          employee.Name,
		Email:         employee.Email,
		ContactNumber: employee.ContactNumber,
		CertificateID: employee.CertificateID,
	}

	// Include tokens in the response
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"employee": employeeResponse,
			"tokens": map[string]interface{}{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}
