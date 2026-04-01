package handlers

import (
	"encoding/json"
	"futbol-api/db"
	"futbol-api/models"
	"net/http"
	"time"
)

func GetPaymentMatrix(w http.ResponseWriter, r *http.Request) {
	// Get all concepts
	conceptRows, err := db.DB.Query(`SELECT id, name, description, created_at FROM payment_concepts ORDER BY id`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer conceptRows.Close()

	concepts := []models.PaymentConcept{}
	for conceptRows.Next() {
		var c models.PaymentConcept
		if err := conceptRows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		concepts = append(concepts, c)
	}

	// Get all players
	playerRows, err := db.DB.Query(`SELECT id, first_name, last_name, dni, birth_date::text, created_at, updated_at FROM players ORDER BY last_name, first_name`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer playerRows.Close()

	players := []models.Player{}
	for playerRows.Next() {
		var p models.Player
		if err := playerRows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.DNI, &p.BirthDate, &p.CreatedAt, &p.UpdatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		players = append(players, p)
	}

	// Get all payments
	paymentRows, err := db.DB.Query(`SELECT id, player_id, concept_id, paid, amount, paid_date::text, notes, updated_at::text FROM player_payments`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer paymentRows.Close()

	paymentMap := make(map[int]map[int]models.PlayerPayment)
	for paymentRows.Next() {
		var pp models.PlayerPayment
		var paidDate, updatedAt *string
		if err := paymentRows.Scan(&pp.ID, &pp.PlayerID, &pp.ConceptID, &pp.Paid, &pp.Amount, &paidDate, &pp.Notes, &updatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		pp.PaidDate = paidDate
		if updatedAt != nil {
			pp.UpdatedAt = *updatedAt
		}
		if _, ok := paymentMap[pp.PlayerID]; !ok {
			paymentMap[pp.PlayerID] = make(map[int]models.PlayerPayment)
		}
		paymentMap[pp.PlayerID][pp.ConceptID] = pp
	}

	// Build matrix rows
	rows := []models.PaymentMatrixRow{}
	for _, p := range players {
		row := models.PaymentMatrixRow{
			Player:   p,
			Payments: paymentMap[p.ID],
		}
		if row.Payments == nil {
			row.Payments = make(map[int]models.PlayerPayment)
		}
		rows = append(rows, row)
	}

	response := models.PaymentMatrixResponse{
		Concepts: concepts,
		Rows:     rows,
	}

	json.NewEncoder(w).Encode(response)
}

func UpdatePayment(w http.ResponseWriter, r *http.Request) {
	var req models.UpdatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	var paidDate *string
	if req.Paid {
		now := time.Now().Format("2006-01-02")
		paidDate = &now
	}

	_, err := db.DB.Exec(
		`INSERT INTO player_payments (player_id, concept_id, paid, amount, paid_date, notes, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW())
		 ON CONFLICT (player_id, concept_id)
		 DO UPDATE SET paid=$3, amount=$4, paid_date=$5, notes=$6, updated_at=NOW()`,
		req.PlayerID, req.ConceptID, req.Paid, req.Amount, paidDate, req.Notes,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
