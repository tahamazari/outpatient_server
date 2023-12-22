package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tahamazari/outpatient_server/api/db"
	"github.com/tahamazari/outpatient_server/api/models"
)

func CreateEmployeeBillingClaim(c echo.Context) error {
	request_billing_claim := new(models.BillingClaim)
	db := db.DB()

	// Binding data
	if err := c.Bind(request_billing_claim); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	new_billing_claim := &models.BillingClaim{
		EmployeeID: request_billing_claim.EmployeeID,
		Name:       request_billing_claim.Name,
		Status:     false,
	}

	if err := db.Create(&new_billing_claim).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": new_billing_claim,
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateEmployeeBillingClaimStatus(c echo.Context) error {
	employeeID := c.Param("employee_id")
	billingClaimID := c.Param("billing_claim_id")
	request_billing_claim := new(models.BillingClaim)
	db := db.DB()

	// Binding data
	if err := c.Bind(request_billing_claim); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	// Check if the billing claim exists and belongs to the requesting employee
	var existingBillingClaim models.BillingClaim
	if err := db.Where("id = ? AND employee_id = ?", billingClaimID, employeeID).First(&existingBillingClaim).Error; err != nil {
		data := map[string]interface{}{
			"message": "Billing Claim not found or does not belong to the requesting employee",
		}
		return c.JSON(http.StatusNotFound, data)
	}

	existing_billing_claim := new(models.BillingClaim)

	if err := db.First(&existing_billing_claim, billingClaimID).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusNotFound, data)
	}

	existing_billing_claim.Status = request_billing_claim.Status
	if err := db.Save(&existing_billing_claim).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existing_billing_claim,
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteEmployeeBillingClaim(c echo.Context) error {
	employeeID := c.Param("employee_id")
	billingClaimID := c.Param("billing_claim_id")
	db := db.DB()

	// Check if the dependent exists and belongs to the requesting employee
	var existingBillingClaim models.BillingClaim
	if err := db.Where("id = ? AND employee_id = ?", billingClaimID, employeeID).First(&existingBillingClaim).Error; err != nil {
		data := map[string]interface{}{
			"message": "Dependent not found or does not belong to the requesting employee",
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Delete the dependent
	if err := db.Delete(&existingBillingClaim).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"message": "Billing Claim has been deleted",
	}
	return c.JSON(http.StatusOK, response)
}
