package database

import (
	"work_kg_backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func GetAdminByEmail(email string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := DB.QueryRow(`SELECT id, email, password, name, role FROM admin_users WHERE email = $1`, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Name, &user.Role)
	return &user, err
}

func GetAdminByEmailWithoutPassword(email string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := DB.QueryRow(`SELECT id, email, name, role FROM admin_users WHERE email = $1`, email).Scan(
		&user.ID, &user.Email, &user.Name, &user.Role)
	return &user, err
}

func ValidateAdminPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
