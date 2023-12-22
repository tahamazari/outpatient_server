package models

type MedicalBill struct {
	ID                 int    `json:"id"`
	EmployeeID         int    `json:"employee_id"`
	BillingClaimID     int    `json:"billing_claim_id"`
	DependentID        int    `json:"dependent_id"`
	Name               string `json:"name"`
	PatientName        string `json:"patient_name"`
	PatientRelation    string `json:"patient_relation"`
	ConsultationFee    int    `json:"consultation_fee"`
	MedicineCharges    int    `json:"medicine_charges"`
	DiagnosticTestFees int    `json:"diagnostic_test_fees"`
	OtherFees          int    `json:"other_fees"`
	TotalBill          int    `json:"total_bill"`
	ImageData          []byte `json:"image_data,omitempty"`
}
