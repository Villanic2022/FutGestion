package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5433")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "futbol")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	migrate()
}

func migrate() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS players (
			id          SERIAL PRIMARY KEY,
			first_name  VARCHAR(100) NOT NULL,
			last_name   VARCHAR(100) NOT NULL,
			dni         VARCHAR(20) UNIQUE NOT NULL,
			birth_date  DATE NOT NULL,
			created_at  TIMESTAMP DEFAULT NOW(),
			updated_at  TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS payment_concepts (
			id          SERIAL PRIMARY KEY,
			name        VARCHAR(200) NOT NULL,
			description TEXT DEFAULT '',
			created_at  TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS player_payments (
			id            SERIAL PRIMARY KEY,
			player_id     INT REFERENCES players(id) ON DELETE CASCADE,
			concept_id    INT REFERENCES payment_concepts(id) ON DELETE CASCADE,
			paid          BOOLEAN DEFAULT FALSE,
			amount        DECIMAL(10,2) DEFAULT 0,
			paid_date     DATE,
			notes         TEXT DEFAULT '',
			updated_at    TIMESTAMP DEFAULT NOW(),
			UNIQUE(player_id, concept_id)
		)`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			log.Fatal("Migration error:", err)
		}
	}
	log.Println("✅ Database migrated")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
