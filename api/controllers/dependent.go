package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tahamazari/outpatient_server/api/db"
	"github.com/tahamazari/outpatient_server/api/models"
	"github.com/tahamazari/outpatient_server/utils"
)

func CreateDependent(c echo.Context) error {
	employeeID, err := utils.ExtractEmployeeIDFromToken(c)
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusUnauthorized, data)
	}

	requestDependent := new(models.Dependent)
	db := db.DB()

	// Binding data
	if err := c.Bind(requestDependent); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	newDependent := &models.Dependent{
		EmployeeID:   employeeID,
		Name:         requestDependent.Name,
		Relationship: requestDependent.Relationship,
	}

	if err := db.Create(&newDependent).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": newDependent,
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateDependent(c echo.Context) error {
	employeeID, err := utils.ExtractEmployeeIDFromToken(c)
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusUnauthorized, data)
	}

	dependentID := c.Param("id")
	requestDependent := new(models.Dependent)
	db := db.DB()

	// Binding data
	if err := c.Bind(requestDependent); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	var existingDependent models.Dependent
	// Check if the dependent exists and belongs to the specified employee
	if err := utils.CheckRecordExistence(utils.ErrDependentNotFound, &existingDependent, "id = ? AND employee_id = ?", []interface{}{dependentID, employeeID}); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusNotFound, data)
	}

	existingDependent.Name = requestDependent.Name
	existingDependent.Relationship = requestDependent.Relationship
	if err := db.Save(&requestDependent).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existingDependent,
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteDependent(c echo.Context) error {
	employeeID, err := utils.ExtractEmployeeIDFromToken(c)
	if err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusUnauthorized, data)
	}

	dependentID := c.Param("dependent_id")
	db := db.DB()

	var existingDependent models.Dependent
	// Check if the dependent exists and belongs to the specified employee
	if err := utils.CheckRecordExistence(utils.ErrDependentNotFound, &existingDependent, "id = ? AND employee_id = ?", []interface{}{dependentID, employeeID}); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Delete the dependent
	if err := db.Delete(&existingDependent).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"message": "Dependent has been deleted",
	}
	return c.JSON(http.StatusOK, response)
}
