package handlers

import (
	"fmt"
	"futbol-api/db"
	"futbol-api/models"
	"net/http"
	"strconv"

	"github.com/go-pdf/fpdf"
)

func ExportPDF(w http.ResponseWriter, r *http.Request) {
	conceptFilter := r.URL.Query().Get("concept_id")

	// Get concepts
	conceptQuery := `SELECT id, name, description, created_at FROM payment_concepts ORDER BY id`
	conceptRows, err := db.DB.Query(conceptQuery)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer conceptRows.Close()

	allConcepts := []models.PaymentConcept{}
	for conceptRows.Next() {
		var c models.PaymentConcept
		if err := conceptRows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		allConcepts = append(allConcepts, c)
	}

	// Filter concepts if needed
	concepts := allConcepts
	if conceptFilter != "" {
		cid, _ := strconv.Atoi(conceptFilter)
		filtered := []models.PaymentConcept{}
		for _, c := range allConcepts {
			if c.ID == cid {
				filtered = append(filtered, c)
			}
		}
		concepts = filtered
	}

	if len(concepts) == 0 {
		http.Error(w, "No concepts found", 404)
		return
	}

	// Get players
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

	// Get payments
	paymentRows, err := db.DB.Query(`SELECT player_id, concept_id, paid, amount FROM player_payments`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer paymentRows.Close()

	paymentMap := make(map[int]map[int]models.PlayerPayment)
	for paymentRows.Next() {
		var pp models.PlayerPayment
		if err := paymentRows.Scan(&pp.PlayerID, &pp.ConceptID, &pp.Paid, &pp.Amount); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if _, ok := paymentMap[pp.PlayerID]; !ok {
			paymentMap[pp.PlayerID] = make(map[int]models.PlayerPayment)
		}
		paymentMap[pp.PlayerID][pp.ConceptID] = pp
	}

	// Generate PDF
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(30, 30, 30)
	pdf.CellFormat(0, 12, "Reporte de Pagos - Seniors", "", 1, "C", false, 0, "")
	
	// Add Logo (Top Right)
	// (x, y, w, h). 277 total width in landscape A4. 277 - 35 = 242 (x pos). 10 = y pos.
	pdf.Image("logo.png", 242, 10, 30, 0, false, "", 0, "")

	pdf.Ln(4)

	// Subtitle
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 100, 100)
	if conceptFilter != "" && len(concepts) > 0 {
		pdf.CellFormat(0, 6, fmt.Sprintf("Filtrado por: %s", concepts[0].Name), "", 1, "C", false, 0, "")
	} else {
		pdf.CellFormat(0, 6, "Todos los conceptos de pago", "", 1, "C", false, 0, "")
	}
	pdf.Ln(6)

	// Table header
	nameColW := 50.0
	dniColW := 30.0
	conceptColW := 35.0
	availableWidth := 277.0 - nameColW - dniColW
	if float64(len(concepts))*conceptColW > availableWidth {
		conceptColW = availableWidth / float64(len(concepts))
	}

	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(44, 62, 80)
	pdf.SetTextColor(255, 255, 255)

	pdf.CellFormat(nameColW, 8, "Jugador", "1", 0, "C", true, 0, "")
	pdf.CellFormat(dniColW, 8, "DNI", "1", 0, "C", true, 0, "")
	for _, c := range concepts {
		name := c.Name
		if len(name) > 15 {
			name = name[:15] + ".."
		}
		pdf.CellFormat(conceptColW, 8, toUTF8(name), "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 9)
	for i, p := range players {
		if i%2 == 0 {
			pdf.SetFillColor(241, 245, 249)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		pdf.SetTextColor(30, 30, 30)
		playerName := fmt.Sprintf("%s %s", p.LastName, p.FirstName)
		if len(playerName) > 25 {
			playerName = playerName[:25] + ".."
		}
		pdf.CellFormat(nameColW, 7, toUTF8(playerName), "1", 0, "L", true, 0, "")
		pdf.CellFormat(dniColW, 7, p.DNI, "1", 0, "C", true, 0, "")

		for _, c := range concepts {
			payment, exists := paymentMap[p.ID][c.ID]
			status := "NO"
			if exists && payment.Paid {
				status = "SI"
				pdf.SetTextColor(39, 174, 96)
			} else {
				pdf.SetTextColor(231, 76, 60)
			}
			pdf.CellFormat(conceptColW, 7, status, "1", 0, "C", true, 0, "")
			pdf.SetTextColor(30, 30, 30)
		}
		pdf.Ln(-1)
	}

	// Summary
	pdf.Ln(8)
	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(30, 30, 30)
	pdf.CellFormat(0, 7, "Resumen:", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	for _, c := range concepts {
		paid := 0
		total := len(players)
		for _, p := range players {
			if pp, ok := paymentMap[p.ID][c.ID]; ok && pp.Paid {
				paid++
			}
		}
		pdf.CellFormat(0, 6, fmt.Sprintf("  %s: %d/%d pagaron", toUTF8(c.Name), paid, total), "", 1, "L", false, 0, "")
	}

	// Output PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=reporte_pagos.pdf")
	if err := pdf.Output(w); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func toUTF8(s string) string {
	// fpdf handles basic ASCII, for special chars we simplify
	result := []byte{}
	for _, b := range []byte(s) {
		if b < 128 {
			result = append(result, b)
		} else {
			result = append(result, '?')
		}
	}
	return string(result)
}
