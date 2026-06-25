package repository

import (
	"auth-service/models"
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) (int64, error) {
	query := `
    INSERT INTO users (username, email, password, created_at, updated_at, location) 
    VALUES (?, ?, ?, NOW(), NOW(), "")
`

	result, err := r.DB.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *UserRepository) FindByIdentifier(identifier string, user *models.User) error {

	return r.DB.QueryRow(`
		SELECT id, email, username, password
		FROM users
		WHERE email = ? OR username = ?
	`, identifier, identifier).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
	)
}
