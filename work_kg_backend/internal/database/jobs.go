package database

import (
	"fmt"
	"log"

	"work_kg_backend/internal/models"
)

func SaveJob(job *models.Job) error {
	_, err := DB.Exec(`INSERT INTO jobs (title, description, category, subcategory, city, salary, phone, company, created_by, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		job.Title, job.Description, job.Category, job.Subcategory, job.City, job.Salary, job.Phone, job.Company, job.CreatedBy, job.Source)
	return err
}

func CreateJob(job *models.Job) error {
	err := DB.QueryRow(`INSERT INTO jobs (title, description, category, subcategory, city, salary, phone, company, is_active, source)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at`,
		job.Title, job.Description, job.Category, job.Subcategory, job.City, job.Salary, job.Phone, job.Company, job.IsActive, job.Source).Scan(&job.ID, &job.CreatedAt)
	return err
}

func UpdateJob(id int64, job *models.Job) error {
	_, err := DB.Exec(`UPDATE jobs SET title=$1, description=$2, category=$3, subcategory=$4, city=$5, salary=$6, phone=$7, company=$8, is_active=$9 WHERE id=$10`,
		job.Title, job.Description, job.Category, job.Subcategory, job.City, job.Salary, job.Phone, job.Company, job.IsActive, id)
	return err
}

func DeleteJob(id int64) error {
	_, err := DB.Exec(`DELETE FROM jobs WHERE id = $1`, id)
	return err
}

func GetAllJobs() ([]models.Job, error) {
	query := `SELECT id, title, description, category, subcategory, city, salary, phone, company, is_active, COALESCE(created_by, 0), source, created_at
		FROM jobs ORDER BY created_at DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]models.Job, 0)
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.Category, &job.Subcategory, &job.City, &job.Salary, &job.Phone, &job.Company, &job.IsActive, &job.CreatedBy, &job.Source, &job.CreatedAt)
		if err != nil {
			log.Printf("Error scanning job: %v", err)
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func SearchJobs(category, subcategory, city string) ([]models.Job, error) {
	query := `SELECT id, title, description, category, subcategory, city, salary, phone, company, created_at
		FROM jobs WHERE is_active = true`
	args := []interface{}{}
	argNum := 1

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argNum)
		args = append(args, category)
		argNum++
	}
	if subcategory != "" {
		query += fmt.Sprintf(" AND subcategory = $%d", argNum)
		args = append(args, subcategory)
		argNum++
	}
	if city != "" {
		query += fmt.Sprintf(" AND city = $%d", argNum)
		args = append(args, city)
		argNum++
	}

	query += " ORDER BY created_at DESC LIMIT 10"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.Job
	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.Category, &job.Subcategory, &job.City, &job.Salary, &job.Phone, &job.Company, &job.CreatedAt)
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}
