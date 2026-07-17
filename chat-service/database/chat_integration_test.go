package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func getChatTestDB(t *testing.T) *DB {
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
	return &DB{db}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func TestChatDatabase_Integration(t *testing.T) {
	db := getChatTestDB(t)
	defer db.Close()

	ctx := context.Background()

	t.Run("Сохранение и получение истории сообщений", func(t *testing.T) {
		timestamp := time.Now().Unix()
		user1 := fmt.Sprintf("chat_u1_%d", timestamp)
		user2 := fmt.Sprintf("chat_u2_%d", timestamp)

		defer func() {
			db.Exec("DELETE FROM messages WHERE sender_id = ? OR receiver_id = ?", user1, user2)
		}()

		err := db.SaveMessage(ctx, user1, user2, "Привет от 1 к 2")
		if err != nil {
			t.Fatalf("Ошибка SaveMessage (1->2): %v", err)
		}

		time.Sleep(1 * time.Second)

		err = db.SaveMessage(ctx, user2, user1, "Привет от 2 к 1")
		if err != nil {
			t.Fatalf("Ошибка SaveMessage (2->1): %v", err)
		}

		history, err := db.GetHistory(ctx, user1, user2, 10)
		if err != nil {
			t.Fatalf("Ошибка GetHistory: %v", err)
		}

		// 3. Проверки
		if len(history) != 2 {
			t.Fatalf("Ожидалось 2 сообщения, получено %d", len(history))
		}

		if history[0].Content != "Привет от 1 к 2" {
			t.Errorf("Первое сообщение должно быть от 1 к 2, получено: %s", history[0].Content)
		}
		if history[1].Content != "Привет от 2 к 1" {
			t.Errorf("Второе сообщение должно быть от 2 к 1, получено: %s", history[1].Content)
		}

		if history[0].Type != "text" || history[1].Type != "text" {
			t.Errorf("Ожидался Type='text', получено: %s, %s", history[0].Type, history[1].Type)
		}

		if history[0].Timestamp == 0 {
			t.Errorf("Timestamp не должен быть равен 0")
		}
	})

	t.Run("Получение истории с ограничением (Limit)", func(t *testing.T) {
		timestamp := time.Now().Unix()
		u1 := fmt.Sprintf("limit_u1_%d", timestamp)
		u2 := fmt.Sprintf("limit_u2_%d", timestamp)

		defer func() {
			db.Exec("DELETE FROM messages WHERE sender_id = ? OR receiver_id = ?", u1, u2)
		}()

		// Вставляем 3 сообщения
		db.SaveMessage(ctx, u1, u2, "msg 1")
		time.Sleep(500 * time.Millisecond)
		db.SaveMessage(ctx, u1, u2, "msg 2")
		time.Sleep(500 * time.Millisecond)
		db.SaveMessage(ctx, u1, u2, "msg 3")

		history, err := db.GetHistory(ctx, u1, u2, 2)
		if err != nil {
			t.Fatalf("Ошибка GetHistory с limit: %v", err)
		}

		if len(history) != 2 {
			t.Errorf("Ожидалось 2 сообщения из-за limit, получено %d", len(history))
		}

		if history[0].Content != "msg 1" || history[1].Content != "msg 2" {
			t.Errorf("Неверные сообщения при limit. Получено: %+v", history)
		}
	})

	t.Run("Пустая история (пользователи не переписывались)", func(t *testing.T) {
		history, err := db.GetHistory(ctx, "non_existent_1", "non_existent_2", 10)

		if err != nil {
			t.Errorf("Ожидалась nil ошибка для пустой истории, получено: %v", err)
		}
		if history == nil {
			t.Errorf("Ожидался пустой срез [], а не nil")
		}
		if len(history) != 0 {
			t.Errorf("Ожидалась длина 0, получено %d", len(history))
		}
	})
}
