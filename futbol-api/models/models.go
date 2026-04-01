package models

import "time"

type Player struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	DNI       string    `json:"dni"`
	BirthDate string    `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaymentConcept struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type PlayerPayment struct {
	ID        int     `json:"id"`
	PlayerID  int     `json:"player_id"`
	ConceptID int     `json:"concept_id"`
	Paid      bool    `json:"paid"`
	Amount    float64 `json:"amount"`
	PaidDate  *string `json:"paid_date"`
	Notes     string  `json:"notes"`
	UpdatedAt string  `json:"updated_at"`
}

type PaymentMatrixRow struct {
	Player   Player                    `json:"player"`
	Payments map[int]PlayerPayment     `json:"payments"`
}

type PaymentMatrixResponse struct {
	Concepts []PaymentConcept   `json:"concepts"`
	Rows     []PaymentMatrixRow `json:"rows"`
}

type UpdatePaymentRequest struct {
	PlayerID  int     `json:"player_id"`
	ConceptID int     `json:"concept_id"`
	Paid      bool    `json:"paid"`
	Amount    float64 `json:"amount"`
	Notes     string  `json:"notes"`
}
