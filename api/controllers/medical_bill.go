package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tahamazari/outpatient_server/api/db"
	"github.com/tahamazari/outpatient_server/api/models"
	"github.com/tahamazari/outpatient_server/utils"
)

// CreateMedicalBill creates a new medical bill
func CreateMedicalBill(c echo.Context) error {
	requestMedicalBill := new(models.MedicalBill)
	db := db.DB()

	// Binding data
	if err := c.Bind(requestMedicalBill); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}

		return c.JSON(http.StatusInternalServerError, data)
	}

	employeeID := requestMedicalBill.EmployeeID
	billingClaimID := requestMedicalBill.BillingClaimID

	fmt.Println(employeeID, billingClaimID, requestMedicalBill)

	// Check if the billing claim exists and belongs to the requesting employee
	var existingBillingClaim models.BillingClaim
	if err := utils.CheckRecordExistence(utils.ErrBillingClaimNotFound, &existingBillingClaim, "id = ? AND employee_id = ?", []interface{}{billingClaimID, employeeID}); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Create the medical bill
	newMedicalBill := &models.MedicalBill{
		EmployeeID:         requestMedicalBill.EmployeeID,
		BillingClaimID:     requestMedicalBill.BillingClaimID,
		DependentID:        requestMedicalBill.DependentID,
		Name:               requestMedicalBill.Name,
		PatientName:        requestMedicalBill.PatientName,
		PatientRelation:    requestMedicalBill.PatientRelation,
		ConsultationFee:    requestMedicalBill.ConsultationFee,
		MedicineCharges:    requestMedicalBill.MedicineCharges,
		DiagnosticTestFees: requestMedicalBill.DiagnosticTestFees,
		OtherFees:          requestMedicalBill.OtherFees,
		TotalBill:          requestMedicalBill.TotalBill,
		ImageData:          requestMedicalBill.ImageData,
	}

	if err := db.Create(&newMedicalBill).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": newMedicalBill,
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateMedicalBill updates an existing medical bill
func UpdateMedicalBill(c echo.Context) error {
	medicalBillID := c.Param("id")
	requestMedicalBill := new(models.MedicalBill)
	db := db.DB()

	// Binding data
	if err := c.Bind(requestMedicalBill); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	employeeID := requestMedicalBill.EmployeeID
	billingClaimID := requestMedicalBill.BillingClaimID

	var existingBillingClaim models.BillingClaim
	if err := utils.CheckRecordExistence(utils.ErrBillingClaimNotFound, &existingBillingClaim, "id = ? AND employee_id = ?", []interface{}{billingClaimID, employeeID}); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Check if the medical bill exists and belongs to the specified billing claim
	var existingMedicalBill models.MedicalBill
	if err := utils.CheckRecordExistence(utils.ErrMedicalBillNotFound, &existingMedicalBill, "id = ? AND billing_claim_id = ?", []interface{}{medicalBillID, billingClaimID}); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Update the medical bill
	existingMedicalBill.Name = requestMedicalBill.Name
	existingMedicalBill.PatientName = requestMedicalBill.PatientName
	existingMedicalBill.PatientRelation = requestMedicalBill.PatientRelation
	existingMedicalBill.ConsultationFee = requestMedicalBill.ConsultationFee
	existingMedicalBill.MedicineCharges = requestMedicalBill.MedicineCharges
	existingMedicalBill.DiagnosticTestFees = requestMedicalBill.DiagnosticTestFees
	existingMedicalBill.OtherFees = requestMedicalBill.OtherFees
	existingMedicalBill.TotalBill = requestMedicalBill.TotalBill
	existingMedicalBill.ImageData = requestMedicalBill.ImageData

	if err := db.Save(&existingMedicalBill).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"data": existingMedicalBill,
	}

	return c.JSON(http.StatusOK, response)
}

// DeleteMedicalBill deletes a medical bill
func DeleteMedicalBill(c echo.Context) error {
	employeeID := c.Param("employee_id")
	billingClaimID := c.Param("billing_claim_id")
	medicalBillID := c.Param("medical_bill_id")
	db := db.DB()

	var existingBillingClaim models.BillingClaim
	if err := utils.CheckRecordExistence(utils.ErrBillingClaimNotFound, &existingBillingClaim, "id = ? AND employee_id = ?", []interface{}{billingClaimID, employeeID}); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Check if the medical bill exists and belongs to the specified billing claim
	var existingMedicalBill models.MedicalBill
	if err := utils.CheckRecordExistence(utils.ErrMedicalBillNotFound, &existingMedicalBill, "id = ? AND billing_claim_id = ?", []interface{}{medicalBillID, billingClaimID}); err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusNotFound, data)
	}

	// Delete the medical bill
	if err := db.Delete(&existingMedicalBill).Error; err != nil {
		data := map[string]interface{}{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	response := map[string]interface{}{
		"message": "Medical Bill has been deleted",
	}
	return c.JSON(http.StatusOK, response)
}
