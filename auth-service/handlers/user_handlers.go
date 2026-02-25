package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("my-super-secret-key-12345")

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

func HandleRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Получен запрос: метод=%s, URL=%s", r.Method, r.URL.Path)

		if r.Method != http.MethodPost {
			http.Error(w, "Вход не разрешён", http.StatusMethodNotAllowed)
			return
		}

		// ограничение по памяти 1МБ
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}
		if req.Email == "" || req.Password == "" || req.Username == "" {
			http.Error(w, "Все поля обязательны", http.StatusBadRequest)
			return
		}
		if len(req.Password) < 4 {
			http.Error(w, "пароль слишком маленький", http.StatusBadRequest)
			return
		}

		hashPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "ошибка обработки пароля", http.StatusBadRequest)
			return
		}

		// готовим вставку в SQL
		query := `
    INSERT INTO users (username, email, password, created_at, updated_at, location) 
    VALUES (?, ?, ?, NOW(), NOW(), "")
`

		result, err := db.Exec(query, req.Username, req.Email, string(hashPass))
		if err != nil {
			log.Printf("Ошибка при сохранении: %v", err)

			if strings.Contains(err.Error(), "Duplicate entry") {
				http.Error(w, "Пользователь с таким email или логином уже существует", http.StatusConflict)
				return
			}
			http.Error(w, "Ошибка сохранения пользователя", http.StatusInternalServerError)
			return
		}
		// получ последнего ID
		userID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Ошибка получения ID ", http.StatusInternalServerError)
			return
		}

		response := UserResponse{
			ID:        userID,
			Username:  req.Username,
			Email:     req.Email,
			CreatedAt: time.Now(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "пользователь создан",
			"user":    response,
		})
	}
}

type UserLogin struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func HandlerLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Проверка метода
		if r.Method != http.MethodPost {
			log.Println(" Неправильный метод запроса")
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
			return
		}

		// Ограничение размера тела
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)

		// Чтение всего тела запроса для отладки
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf(" Ошибка чтения тела запроса: %v", err)
			http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
			return
		}

		log.Printf("Тело запроса: %s", string(body))

		// Парсим JSON
		var req UserLogin
		if err := json.Unmarshal(body, &req); err != nil {
			log.Printf(" Ошибка парсинга JSON: %v", err)
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		// Валидация
		if req.Identifier == "" || req.Password == "" || len(req.Password) < 4 {
			log.Println(" Ошибка валидации: пустые поля или короткий пароль")
			http.Error(w, "Введите корректные данные", http.StatusBadRequest)
			return
		}

		var userId int
		var realEmail string
		var storPassHash string
		var realUsername string

		// Ищем пользователя в БД
		err = db.QueryRow(`
    SELECT id, email, username, password 
    FROM users 
    WHERE email = ? OR username = ?
`, req.Identifier, req.Identifier).Scan(&userId, &realEmail, &realUsername, &storPassHash)

		if err == sql.ErrNoRows {
			log.Printf(" Пользователь не найден: %s", req.Identifier)
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}

		if err != nil {
			log.Printf(" Ошибка БД: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		log.Printf(" Пользователь найден: id=%d", userId)

		// Проверяем пароль
		if err := bcrypt.CompareHashAndPassword([]byte(storPassHash), []byte(req.Password)); err != nil {
			log.Printf(" Пароли не совпадают: %v", err)
			http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
			return
		}

		log.Println("Пароль верный")

		// Генерируем токен
		claims := jwt.MapClaims{
			"user_id":  userId,
			"email":    realEmail,
			"username": realUsername,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		}
		log.Println(realUsername)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)

		http.SetCookie(w, &http.Cookie{
			Name:  "auth_token",
			Value: tokenString,
			Path:  "/",
		})

		if err != nil {
			log.Printf(" Ошибка генерации токена: %v", err)
			http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
			return
		}

		log.Println(" Токен сгенерирован успешно")

		// Возвращаем токен
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

		log.Println("Ответ отправлен клиенту")
	}
}
