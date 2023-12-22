package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tahamazari/outpatient_server/api/db"
	"github.com/tahamazari/outpatient_server/api/models"
	"github.com/tahamazari/outpatient_server/utils"
)

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

	// Clear the password before sending the response
	newEmployee.Password = ""

	// Include tokens in the response
	response := map[string]interface{}{
		"data": map[string]interface{}{
			"employee": newEmployee,
			"tokens": map[string]interface{}{
				"access_token":  accessToken,
				"refresh_token": refreshToken,
			},
		},
	}

	return c.JSON(http.StatusOK, response)
}
