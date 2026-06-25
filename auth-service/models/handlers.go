package models

import "time"

// структура ДЛЯ сервера
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// структура ИЗ сервера
type UserResponse struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// login

type UserLogin struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}
