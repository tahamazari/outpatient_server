package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tahamazari/outpatient_server/api/db"
	"github.com/tahamazari/outpatient_server/api/models"
)

func CreateEmployee(c echo.Context) error {
	employee_model := new(models.Employee)
	db := db.DB()

	// Binding data
	if err := c.Bind(employee_model); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	new_employee := &models.Employee{
		Name:  employee_model.Name,
		Email: employee_model.Email,
	}

	if err := db.Create(&new_employee).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": new_employee,
	}

	return c.JSON(http.StatusOK, response)
}

// GetEmployee retrieves an employee by ID
func GetEmployee(c echo.Context) error {
	id := c.Param("id")
	db := db.DB()

	var employees []*models.Employee

	if res := db.Find(&employees, id); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusOK, data)
	}

	response := map[string]interface{}{
		"data": employees[0],
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateEmployee(c echo.Context) error {
	id := c.Param("id")
	employee_model := new(models.Employee)
	db := db.DB()

	// Binding data
	if err := c.Bind(employee_model); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	existing_employee := new(models.Employee)

	if err := db.First(&existing_employee, id).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	existing_employee.Name = employee_model.Name
	existing_employee.Email = employee_model.Email
	if err := db.Save(&existing_employee).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existing_employee,
	}

	return c.JSON(http.StatusOK, response)
}

func GetEmployeeDependents(c echo.Context) error {
	employeeID := c.Param("id")
	db := db.DB()

	var employeeDependents []*models.Dependent

	// Assuming there is a foreign key relationship between Employee and Dependent
	if res := db.Where("employee_id", employeeID).Find(&employeeDependents); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": employeeDependents,
	}

	return c.JSON(http.StatusOK, response)
}

func GetEmployeeBillingClaims(c echo.Context) error {
	employeeID := c.Param("id")
	db := db.DB()

	var employeeBillingClaims []*models.BillingClaim

	// Assuming there is a foreign key relationship between Employee and Dependent
	if res := db.Where("employee_id", employeeID).Find(&employeeBillingClaims); res.Error != nil {
		data := map[string]interface{}{
			"message": res.Error.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": employeeBillingClaims,
	}

	return c.JSON(http.StatusOK, response)
}
