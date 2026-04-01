package handlers

import (
	"encoding/json"
	"futbol-api/db"
	"futbol-api/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetConcepts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, name, description, created_at FROM payment_concepts ORDER BY id`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	concepts := []models.PaymentConcept{}
	for rows.Next() {
		var c models.PaymentConcept
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		concepts = append(concepts, c)
	}
	json.NewEncoder(w).Encode(concepts)
}

func CreateConcept(w http.ResponseWriter, r *http.Request) {
	var c models.PaymentConcept
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}
	if c.Name == "" {
		http.Error(w, "Name is required", 400)
		return
	}

	err := db.DB.QueryRow(
		`INSERT INTO payment_concepts (name, description) VALUES ($1, $2) RETURNING id, created_at`,
		c.Name, c.Description,
	).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Auto-create payment entries for all existing players
	db.DB.Exec(`INSERT INTO player_payments (player_id, concept_id) SELECT id, $1 FROM players ON CONFLICT DO NOTHING`, c.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func UpdateConcept(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var c models.PaymentConcept
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	idInt, _ := strconv.Atoi(id)
	_, err := db.DB.Exec(`UPDATE payment_concepts SET name=$1, description=$2 WHERE id=$3`, c.Name, c.Description, idInt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	c.ID = idInt
	json.NewEncoder(w).Encode(c)
}

func DeleteConcept(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.DB.Exec(`DELETE FROM payment_concepts WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
