package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tahamazari/outpatient_server/api/db"
	"github.com/tahamazari/outpatient_server/api/models"
)

func CreateDependent(c echo.Context) error {
	request_dependent := new(models.Dependent)
	db := db.DB()

	// Binding data
	if err := c.Bind(request_dependent); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	new_dependent := &models.Dependent{
		EmployeeId:   request_dependent.EmployeeId,
		Name:         request_dependent.Name,
		Relationship: request_dependent.Relationship,
	}

	if err := db.Create(&new_dependent).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": new_dependent,
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateDependent(c echo.Context) error {
	id := c.Param("id")
	request_dependent := new(models.Dependent)
	db := db.DB()

	// Binding data
	if err := c.Bind(request_dependent); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	existing_dependent := new(models.Dependent)

	if err := db.First(&existing_dependent, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	existing_dependent.Name = request_dependent.Name
	existing_dependent.Relationship = request_dependent.Relationship
	if err := db.Save(&existing_dependent).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existing_dependent,
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteDependent(c echo.Context) error {
	employeeID := c.Param("employee_id")
	dependentID := c.Param("dependent_id")
	db := db.DB()

	// Check if the dependent exists and belongs to the requesting employee
	var existingDependent models.Dependent
	if err := db.Where("id = ? AND employee_id = ?", dependentID, employeeID).First(&existingDependent).Error; err != nil {
		data := map[string]interface{}{
			"message": "Dependent not found or does not belong to the requesting employee",
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
