package models

import "time"

type User struct {
	ID         int64     `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Phone      string    `json:"phone"`
	City       string    `json:"city"`
	Specialty  string    `json:"specialty"`
	Experience string    `json:"experience"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"created_at"`
}

type AdminUser struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Job struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Subcategory string    `json:"subcategory"`
	City        string    `json:"city"`
	Salary      string    `json:"salary"`
	Phone       string    `json:"phone"`
	Company     string    `json:"company"`
	IsActive    bool      `json:"is_active"`
	CreatedBy   int64     `json:"created_by"`
	Source      string    `json:"source"`
	CreatedAt   time.Time `json:"created_at"`
}

type Resume struct {
	ID         int64     `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	City       string    `json:"city"`
	Specialty  string    `json:"specialty"`
	Experience string    `json:"experience"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserState struct {
	State       string
	Category    string
	Subcategory string
	City        string
	SearchType  string
	TempJob     *Job
	// Form data
	FormName       string
	FormPhone      string
	FormCity       string
	FormSpecialty  string
	FormExperience string
	// Message IDs for deletion (collect all, delete at end)
	FormMessageIDs []int
}

type Stats struct {
	TotalJobs    int `json:"total_jobs"`
	ActiveJobs   int `json:"active_jobs"`
	TotalUsers   int `json:"total_users"`
	TotalResumes int `json:"total_resumes"`
	TodayJobs    int `json:"today_jobs"`
	TodayUsers   int `json:"today_users"`
	TodayResumes int `json:"today_resumes"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string                 `json:"token"`
	User  map[string]interface{} `json:"user"`
}
