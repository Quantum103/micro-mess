package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"auth-service/models"

	_ "github.com/go-sql-driver/mysql"
)

func getTestDB(t *testing.T) *sql.DB {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			getEnv("DB_USER", "root"),
			getEnv("DB_PASSWORD", "password"),
			getEnv("DB_HOST", "127.0.0.1"),
			getEnv("DB_PORT", "3306"),
			getEnv("DB_NAME", "auth_test_db"),
		)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Не удалось открыть соединение с БД: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Не удалось подключиться к БД (ping failed): %v. Убедитесь, что тестовая БД запущена.", err)
	}

	return db
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func TestUserRepository_CreateAndFind_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	timestamp := time.Now().Unix()
	testEmail := fmt.Sprintf("integration_test_%d@example.com", timestamp)
	testUsername := fmt.Sprintf("testuser_%d", timestamp)
	testPassword := "hashed_password_123"

	newUser := &models.User{
		Name:     testUsername,
		Email:    testEmail,
		Password: testPassword,
	}

	defer func() {
		_, _ = db.Exec("DELETE FROM users WHERE email = ?", testEmail)
	}()

	insertedID, err := repo.Create(newUser)

	if err != nil {
		t.Fatalf("Ожидалось успешное создание пользователя, но получена ошибка: %v", err)
	}
	if insertedID == 0 {
		t.Fatal("Ожидался ID > 0 после создания")
	}

	newUser.ID = uint(insertedID)
	var foundUser models.User
	err = repo.FindByIdentifier(testEmail, &foundUser)

	if err != nil {
		if err == sql.ErrNoRows {
			t.Fatalf("Пользователь не найден в базе, хотя мы его только что создали")
		}
		t.Fatalf("Неожиданная ошибка при поиске пользователя: %v", err)
	}

	if foundUser.Email != newUser.Email {
		t.Errorf("Ожидался email %s, получен %s", newUser.Email, foundUser.Email)
	}
	if foundUser.Name != newUser.Name {
		t.Errorf("Ожидалось имя %s, получено %s", newUser.Name, foundUser.Name)
	}
	if foundUser.Password != newUser.Password {
		t.Errorf("Ожидался пароль %s, получен %s", newUser.Password, foundUser.Password)
	}
}
