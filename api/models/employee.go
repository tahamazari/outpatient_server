package models

type Employee struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	ContactNumber string `json:"contact_number"`
	CertificateID string `json:"certificate_id"`
}
