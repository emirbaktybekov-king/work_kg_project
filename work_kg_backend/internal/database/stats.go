package database

import "work_kg_backend/internal/models"

func GetStats() models.Stats {
	var stats models.Stats

	DB.QueryRow(`SELECT COUNT(*) FROM jobs`).Scan(&stats.TotalJobs)
	DB.QueryRow(`SELECT COUNT(*) FROM jobs WHERE is_active = true`).Scan(&stats.ActiveJobs)
	DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	DB.QueryRow(`SELECT COUNT(*) FROM resumes`).Scan(&stats.TotalResumes)
	DB.QueryRow(`SELECT COUNT(*) FROM jobs WHERE DATE(created_at) = CURRENT_DATE`).Scan(&stats.TodayJobs)
	DB.QueryRow(`SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURRENT_DATE`).Scan(&stats.TodayUsers)
	DB.QueryRow(`SELECT COUNT(*) FROM resumes WHERE DATE(created_at) = CURRENT_DATE`).Scan(&stats.TodayResumes)

	return stats
}
