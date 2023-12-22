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

	var existingEmployee models.Employee
	if err := db.Where("email = ?", requestSignupEmployee.Email).First(&existingEmployee).Error; err == nil {
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
		ContactNumber: requestSignupEmployee.ContactNumber,
		CertificateID: requestSignupEmployee.CertificateID,
	}

	if err := db.Create(&newEmployee).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	// Clear the password before sending the response
	newEmployee.Password = ""
	response := map[string]interface{}{
		"data": newEmployee,
	}

	return c.JSON(http.StatusOK, response)
}
