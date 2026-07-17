package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MockHub struct {
	onlineUsers map[string]bool
}

func (m *MockHub) IsOnline(userID string) bool {
	return m.onlineUsers[userID]
}

func getTestDB(t *testing.T) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_HOST", "127.0.0.1"),
		getEnv("DB_PORT", "3307"),
		getEnv("DB_NAME", "chat_test_db"),
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Не удалось открыть БД: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	return db
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func TestGetFriendsHandler_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	hub := &MockHub{onlineUsers: make(map[string]bool)}
	handler := GetFriendsHandler(db, hub)

	t.Run("Успешное получение списка друзей с онлайн статусом", func(t *testing.T) {
		timestamp := time.Now().Unix()
		myID := fmt.Sprintf("user_main_%d", timestamp)
		friendID := fmt.Sprintf("user_friend_%d", timestamp)

		db.Exec("INSERT INTO users (id, username) VALUES (?, ?)", myID, "MainUser")
		db.Exec("INSERT INTO users (id, username) VALUES (?, ?)", friendID, "FriendUser")
		db.Exec("INSERT INTO friends (user_id, friend_id, status) VALUES (?, ?, 'accepted')", myID, friendID)

		hub.onlineUsers[friendID] = true

		defer func() {
			db.Exec("DELETE FROM friends WHERE user_id = ? OR friend_id = ?", myID, friendID)
			db.Exec("DELETE FROM users WHERE id = ? OR id = ?", myID, friendID)
		}()

		req := httptest.NewRequest(http.MethodGet, "/friends", nil)
		req.Header.Set("X-User-ID", myID)
		rr := httptest.NewRecorder()

		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Ожидался статус 200, получен %d", rr.Code)
		}

		var users []User
		if err := json.NewDecoder(rr.Body).Decode(&users); err != nil {
			t.Fatalf("Ошибка парсинга JSON: %v", err)
		}

		if len(users) != 1 {
			t.Fatalf("Ожидался 1 друг, получено %d", len(users))
		}

		if users[0].ID != friendID || users[0].Username != "FriendUser" {
			t.Errorf("Неверные данные друга: %+v", users[0])
		}

		if !users[0].Online {
			t.Errorf("Ожидалось, что друг будет онлайн (Online=true)")
		}
	})

	t.Run("Отсутствие заголовка X-User-ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/friends", nil)
		rr := httptest.NewRecorder()

		handler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("Ожидался статус 401, получен %d", rr.Code)
		}
	})

	t.Run("Пустой список друзей (должен быть [], а не null)", func(t *testing.T) {
		timestamp := time.Now().Unix()
		loneUserID := fmt.Sprintf("user_lone_%d", timestamp)

		db.Exec("INSERT INTO users (id, username) VALUES (?, ?)", loneUserID, "LoneUser")
		defer db.Exec("DELETE FROM users WHERE id = ?", loneUserID)

		req := httptest.NewRequest(http.MethodGet, "/friends", nil)
		req.Header.Set("X-User-ID", loneUserID)
		rr := httptest.NewRecorder()

		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Ожидался статус 200, получен %d", rr.Code)
		}

		if string(rr.Body.Bytes()) != "[]\n" {
			t.Errorf("Ожидался пустой массив '[]\\n', получено: %s", string(rr.Body.Bytes()))
		}
	})
}
