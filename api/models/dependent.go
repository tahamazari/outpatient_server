package models

type Dependent struct {
	ID           int    `json:"id"`
	EmployeeId   int    `json:"employee_id"`
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
}
