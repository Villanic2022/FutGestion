package main

import (
	"fmt"
	"futbol-api/db"
	"futbol-api/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	db.Init()

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// Players
	r.HandleFunc("/api/players", handlers.GetPlayers).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/players/{id}", handlers.GetPlayer).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/players", handlers.CreatePlayer).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/players/{id}", handlers.UpdatePlayer).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/players/{id}", handlers.DeletePlayer).Methods("DELETE", "OPTIONS")

	// Payment Concepts
	r.HandleFunc("/api/concepts", handlers.GetConcepts).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/concepts", handlers.CreateConcept).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/concepts/{id}", handlers.UpdateConcept).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/concepts/{id}", handlers.DeleteConcept).Methods("DELETE", "OPTIONS")

	// Payments
	r.HandleFunc("/api/payments/matrix", handlers.GetPaymentMatrix).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/payments", handlers.UpdatePayment).Methods("PUT", "OPTIONS")

	// PDF Export
	r.HandleFunc("/api/export/pdf", handlers.ExportPDF).Methods("GET", "OPTIONS")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
