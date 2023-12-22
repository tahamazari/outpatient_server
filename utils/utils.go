package utils

import (
	"fmt"

	"github.com/tahamazari/outpatient_server/api/db"
)

type ErrorMessageArgs struct {
	ErrorMessage string
	Model        interface{}
	Query        string
	Args         []interface{}
}

func CheckRecordExistence(errorMessage string, model interface{}, query string, args []interface{}) error {
	db := db.DB()

	if err := db.Where(query, args...).First(model).Error; err != nil {
		return fmt.Errorf(errorMessage)
	}

	return nil
}

const ErrDependentNotFound = "Dependent not found or does not belong to the requesting employee"
const ErrBillingClaimNotFound = "Billing Claim not found or does not belong to the requesting employee"
const ErrMedicalBillNotFound = "Medical Bill not found or does not belong to the specified Billing Claim"
