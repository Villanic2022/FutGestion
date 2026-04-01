package handlers

import (
	"database/sql"
	"encoding/json"
	"futbol-api/db"
	"futbol-api/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetPlayers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`SELECT id, first_name, last_name, dni, birth_date::text, created_at, updated_at FROM players ORDER BY last_name, first_name`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	players := []models.Player{}
	for rows.Next() {
		var p models.Player
		if err := rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.DNI, &p.BirthDate, &p.CreatedAt, &p.UpdatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		players = append(players, p)
	}
	json.NewEncoder(w).Encode(players)
}

func GetPlayer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p models.Player
	err := db.DB.QueryRow(`SELECT id, first_name, last_name, dni, birth_date::text, created_at, updated_at FROM players WHERE id = $1`, id).
		Scan(&p.ID, &p.FirstName, &p.LastName, &p.DNI, &p.BirthDate, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Player not found", 404)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var p models.Player
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}
	if p.FirstName == "" || p.LastName == "" || p.DNI == "" || p.BirthDate == "" {
		http.Error(w, "All fields are required: first_name, last_name, dni, birth_date", 400)
		return
	}

	err := db.DB.QueryRow(
		`INSERT INTO players (first_name, last_name, dni, birth_date) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		p.FirstName, p.LastName, p.DNI, p.BirthDate,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Auto-create payment entries for all existing concepts
	db.DB.Exec(`INSERT INTO player_payments (player_id, concept_id) SELECT $1, id FROM payment_concepts ON CONFLICT DO NOTHING`, p.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p models.Player
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	idInt, _ := strconv.Atoi(id)
	_, err := db.DB.Exec(
		`UPDATE players SET first_name=$1, last_name=$2, dni=$3, birth_date=$4, updated_at=NOW() WHERE id=$5`,
		p.FirstName, p.LastName, p.DNI, p.BirthDate, idInt,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	p.ID = idInt
	json.NewEncoder(w).Encode(p)
}

func DeletePlayer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	_, err := db.DB.Exec(`DELETE FROM players WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
