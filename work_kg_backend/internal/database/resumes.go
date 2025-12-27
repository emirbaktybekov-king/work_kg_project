package database

import (
	"log"

	"work_kg_backend/internal/models"
)

func SaveResume(telegramID int64, username, name, phone, city, specialty, experience string) error {
	_, err := DB.Exec(`INSERT INTO resumes (telegram_id, username, name, phone, city, specialty, experience)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (telegram_id) DO UPDATE SET
		username = EXCLUDED.username,
		name = EXCLUDED.name,
		phone = EXCLUDED.phone,
		city = EXCLUDED.city,
		specialty = EXCLUDED.specialty,
		experience = EXCLUDED.experience,
		updated_at = CURRENT_TIMESTAMP`,
		telegramID, username, name, phone, city, specialty, experience)
	if err != nil {
		log.Printf("Error saving resume: %v", err)
	}
	return err
}

func UpdateUserFormData(telegramID int64, name, phone, city, specialty, experience string) error {
	_, err := DB.Exec(`UPDATE users SET
		first_name = COALESCE(NULLIF($1, ''), first_name),
		phone = $2,
		city = $3,
		specialty = $4,
		experience = $5
		WHERE telegram_id = $6`,
		name, phone, city, specialty, experience, telegramID)
	if err != nil {
		log.Printf("Error saving form data to users: %v", err)
	}
	return err
}

func GetAllResumes() ([]models.Resume, error) {
	rows, err := DB.Query(`SELECT id, telegram_id, COALESCE(username, ''), COALESCE(name, ''),
		COALESCE(phone, ''), COALESCE(city, ''), COALESCE(specialty, ''), COALESCE(experience, ''),
		created_at, updated_at FROM resumes ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resumes []models.Resume
	for rows.Next() {
		var resume models.Resume
		err := rows.Scan(&resume.ID, &resume.TelegramID, &resume.Username, &resume.Name, &resume.Phone,
			&resume.City, &resume.Specialty, &resume.Experience, &resume.CreatedAt, &resume.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning resume: %v", err)
			continue
		}
		resumes = append(resumes, resume)
	}

	return resumes, nil
}
