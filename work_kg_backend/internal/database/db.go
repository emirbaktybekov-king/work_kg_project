package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func Connect(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to database successfully")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func InitSchema() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			telegram_id BIGINT UNIQUE,
			username VARCHAR(255),
			first_name VARCHAR(255),
			last_name VARCHAR(255),
			phone VARCHAR(50),
			city VARCHAR(100),
			specialty VARCHAR(255),
			experience TEXT,
			role VARCHAR(50) DEFAULT 'user',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS specialty VARCHAR(255)`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS experience TEXT`,
		`CREATE TABLE IF NOT EXISTS admin_users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			name VARCHAR(255),
			role VARCHAR(50) DEFAULT 'admin',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS jobs (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			category VARCHAR(100),
			subcategory VARCHAR(100),
			city VARCHAR(100),
			salary VARCHAR(100),
			phone VARCHAR(50),
			company VARCHAR(255),
			is_active BOOLEAN DEFAULT true,
			created_by BIGINT,
			source VARCHAR(50) DEFAULT 'telegram',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS resumes (
			id SERIAL PRIMARY KEY,
			telegram_id BIGINT UNIQUE,
			username VARCHAR(255),
			name VARCHAR(255),
			phone VARCHAR(50),
			city VARCHAR(100),
			specialty VARCHAR(255),
			experience TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_resumes_telegram_id ON resumes(telegram_id)`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_category ON jobs(category)`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_city ON jobs(city)`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_is_active ON jobs(is_active)`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Printf("Error executing query: %v", err)
		}
	}

	// Create default admin user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	_, err := DB.Exec(`INSERT INTO admin_users (email, password, name, role)
		VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING`,
		"admin@workkg.com", string(hashedPassword), "Admin", "admin")
	if err != nil {
		log.Printf("Error creating admin user: %v", err)
	}

	log.Println("Database initialized successfully")
}
