package database

import (
	"database/sql"
	"log"

	"work_kg_backend/internal/models"
)

func SaveUser(telegramID int64, username, firstName, lastName, city string) error {
	_, err := DB.Exec(`INSERT INTO users (telegram_id, username, first_name, last_name, city)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (telegram_id) DO UPDATE SET
		username = EXCLUDED.username,
		first_name = EXCLUDED.first_name,
		last_name = EXCLUDED.last_name`,
		telegramID, username, firstName, lastName, city)
	if err != nil {
		log.Printf("Error saving user: %v", err)
	}
	return err
}

func GetUserByTelegramID(telegramID int64) (*models.User, error) {
	var user models.User
	var specialty, experience sql.NullString

	err := DB.QueryRow(`SELECT id, telegram_id, COALESCE(username, ''), COALESCE(first_name, ''),
		COALESCE(last_name, ''), COALESCE(phone, ''), COALESCE(city, ''), specialty, experience, role, created_at
		FROM users WHERE telegram_id = $1`, telegramID).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName,
		&user.Phone, &user.City, &specialty, &experience, &user.Role, &user.CreatedAt)

	if specialty.Valid {
		user.Specialty = specialty.String
	}
	if experience.Valid {
		user.Experience = experience.String
	}

	return &user, err
}

func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query(`SELECT id, telegram_id, COALESCE(username, ''), COALESCE(first_name, ''),
		COALESCE(last_name, ''), COALESCE(phone, ''), COALESCE(city, ''), COALESCE(specialty, ''),
		COALESCE(experience, ''), role, created_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName,
			&user.Phone, &user.City, &user.Specialty, &user.Experience, &user.Role, &user.CreatedAt)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUsernameByTelegramID(telegramID int64) string {
	var username string
	DB.QueryRow(`SELECT COALESCE(username, '') FROM users WHERE telegram_id = $1`, telegramID).Scan(&username)
	return username
}
