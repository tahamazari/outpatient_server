package models

type BillingClaim struct {
	ID         int    `json:"id"`
	EmployeeID int    `json:"employee_id"`
	Name       string `json:"name"`
	Status     bool   `json:"status"`
}
