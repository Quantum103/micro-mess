package handlers

import (
	"auth-service/models"
	"auth-service/repository"
	"auth-service/service"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
)

func getTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			getEnv("DB_USER", "root"),
			getEnv("DB_PASSWORD", "password"),
			getEnv("DB_HOST", "127.0.0.1"),
			getEnv("DB_PORT", "3307"),
			getEnv("DB_NAME", "auth_test_db"),
		)
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Не удалось открыть БД: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("Не удалось подключиться к БД: %v. Запустите docker-compose.test.yml", err)
	}
	return db
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func TestAuth_RegisterAndLogin_E2E(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := NewAuthHandler(authService)

	timestamp := time.Now().Unix()
	testEmail := fmt.Sprintf("e2e_user_%d@test.com", timestamp)
	testUsername := fmt.Sprintf("e2euser%d", timestamp)
	testPassword := "SecurePass123!"

	defer func() {
		_, _ = db.Exec("DELETE FROM users WHERE email = ?", testEmail)
	}()

	t.Run("Регистрация нового пользователя", func(t *testing.T) {
		// Формируем и отправляем HTTP-запрос
		regReq := models.RegisterRequest{
			Username: testUsername,
			Email:    testEmail,
			Password: testPassword,
		}
		body, _ := json.Marshal(regReq)

		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		authHandler.Register(rr, req)

		if rr.Code != http.StatusCreated {
			t.Fatalf("Ожидался статус 201 Created, получен %d. Тело: %s", rr.Code, rr.Body.String())
		}

		var resp map[string]interface{}
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("Не удалось распарсить ответ регистрации: %v", err)
		}
		if resp["message"] != "пользователь создан" {
			t.Errorf("Неверное сообщение в ответе: %v", resp["message"])
		}

		//  Проверяем, что данные записались в БД правильно
		var dbPassword string
		err := db.QueryRow("SELECT password FROM users WHERE email = ?", testEmail).Scan(&dbPassword)
		if err != nil {
			t.Fatalf("Пользователь не найден в базе данных после регистрации: %v", err)
		}
		if dbPassword == testPassword {
			t.Errorf("ОШИБКА БЕЗОПАСНОСТИ: Пароль сохранен в открытом виде, а не захеширован!")
		}
		if !strings.HasPrefix(dbPassword, "$2a$") && !strings.HasPrefix(dbPassword, "$2b$") {
			t.Errorf("Ожидался bcrypt-хеш пароля, получено: %s", dbPassword)
		}
	})

	t.Run("Успешный логин с полученными данными", func(t *testing.T) {
		loginReq := models.UserLogin{
			Identifier: testEmail,
			Password:   testPassword,
		}
		body, _ := json.Marshal(loginReq)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		authHandler.Login(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Ожидался статус 200 OK, получен %d. Тело: %s", rr.Code, rr.Body.String())
		}

		cookies := rr.Result().Cookies()
		var authToken string
		for _, cookie := range cookies {
			if cookie.Name == "auth_token" {
				authToken = cookie.Value
				break
			}
		}
		if authToken == "" {
			t.Error("Cookie 'auth_token' не была установлена в ответе")
		}

		var resp map[string]string
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Fatalf("Не удалось распарсить ответ логина: %v", err)
		}
		if resp["token"] == "" {
			t.Error("Токен отсутствует в теле ответа")
		}

		token, err := jwt.Parse(resp["token"], func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			t.Fatalf("Возвращенный токен невалиден или ошибка парсинга: %v", err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal("Не удалось привести claims к MapClaims")
		}

		if claims["email"] != testEmail {
			t.Errorf("Ожидался email в токене %s, получен %v", testEmail, claims["email"])
		}
		if claims["username"] != testUsername {
			t.Errorf("Ожидался username в токене %s, получен %v", testUsername, claims["username"])
		}
	})

	t.Run("Логин с неверным паролем", func(t *testing.T) {
		loginReq := models.UserLogin{
			Identifier: testEmail,
			Password:   "WrongPassword123",
		}
		body, _ := json.Marshal(loginReq)

		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		authHandler.Login(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Ожидался статус 401 Unauthorized при неверном пароле, получен %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "invalid credentials") {
			t.Errorf("Ожидалось сообщение 'invalid credentials', получено: %s", rr.Body.String())
		}
	})
}
