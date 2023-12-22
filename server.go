package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tahamazari/outpatient_server/api/controllers"
	"github.com/tahamazari/outpatient_server/api/db"
)

func main() {
	// Create an instance of Echo
	e := echo.New()

	// Define a route
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, this is your API!")
	})

	db.DatabaseInit()
	gorm := db.DB()

	dbGorm, err := gorm.DB()
	if err != nil {
		panic(err)
	}

	dbGorm.Ping()

	e.POST("/employees", controllers.CreateEmployee)
	e.GET("/employees/:id", controllers.GetEmployee)
	e.PUT("/employees/:id", controllers.UpdateEmployee)
	e.GET("/employees/:id/dependents", controllers.GetEmployeeDependents)
	e.GET("/employees/:id/billing_claims", controllers.GetEmployeeBillingClaims)

	e.POST("/dependents", controllers.CreateDependent)
	e.PUT("/dependents/:id", controllers.UpdateDependent)
	e.DELETE("/dependents/:employee_id/:dependent_id", controllers.DeleteDependent)

	e.POST("/billing_claims", controllers.CreateEmployeeBillingClaim)
	e.PUT("/billing_claims/:employee_id/:billing_claim_id", controllers.UpdateEmployeeBillingClaimStatus)
	e.DELETE("/billing_claims/:employee_id/:billing_claim_id", controllers.DeleteEmployeeBillingClaim)

	e.POST("/medical_bills", controllers.CreateMedicalBill)
	e.PUT("/medical_bills/:id", controllers.UpdateMedicalBill)
	e.DELETE("/medical_bills/:employee_id/:billing_claim_id/:medical_bill_id", controllers.DeleteMedicalBill)

	e.POST("/auth/sign_up", controllers.SignUp)
	e.POST("/auth/login", controllers.Login)

	// Start the server
	e.Start(":8080")
}
