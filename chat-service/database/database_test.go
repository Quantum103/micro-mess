package database

import (
	"context"
	"fmt"
	"os"
	"testing"

	// Убедитесь, что путь соответствует вашей структуре проекта
	"github.com/DATA-DOG/go-sqlmock"
)

func TestSaveMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("неожиданная ошибка при создании mock базы данных: %s", err)
	}
	defer db.Close()

	dbWrapper := &DB{db}

	t.Run("успешное сохранение", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO messages").
			WithArgs("user1", "user2", "hello world").
			WillReturnResult(sqlmock.NewResult(1, 1)) // last insert id = 1, rows affected = 1

		err := dbWrapper.SaveMessage(context.Background(), "user1", "user2", "hello world")
		if err != nil {
			t.Errorf("ошибка не ожидалась, но получена: %s", err)
		}

		// Проверяем, что все ожидаемые вызовы были выполнены
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("невыполненные ожидания: %s", err)
		}
	})

	t.Run("ошибка базы данных", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO messages").
			WithArgs("user1", "user2", "hello world").
			WillReturnError(fmt.Errorf("ошибка соединения с БД"))

		err := dbWrapper.SaveMessage(context.Background(), "user1", "user2", "hello world")
		if err == nil {
			t.Errorf("ожидалась ошибка, но получена nil")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("невыполненные ожидания: %s", err)
		}
	})
}

func TestGetHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("неожиданная ошибка при создании mock базы данных: %s", err)
	}
	defer db.Close()

	dbWrapper := &DB{db}

	t.Run("успешная загрузка истории", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"sender_id", "receiver_id", "content", "timestamp"}).
			AddRow("user1", "user2", "привет", 1690000000).
			AddRow("user2", "user1", "пока", 1690000010)

		mock.ExpectQuery("SELECT sender_id, receiver_id, content, UNIX_TIMESTAMP").
			WithArgs("user1", "user2", "user2", "user1", 10).
			WillReturnRows(rows)

		messages, err := dbWrapper.GetHistory(context.Background(), "user1", "user2", 10)
		if err != nil {
			t.Errorf("ошибка не ожидалась, но получена: %s", err)
		}

		if len(messages) != 2 {
			t.Errorf("ожидалось 2 сообщения, получено: %d", len(messages))
		}

		if messages[0].Type != "text" {
			t.Errorf("ожидался тип 'text', получено: %s", messages[0].Type)
		}

		if messages[0].From != "user1" || messages[0].Content != "привет" {
			t.Errorf("неверные данные сообщения: %+v", messages[0])
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("невыполненные ожидания: %s", err)
		}
	})

	t.Run("ошибка при выполнении запроса", func(t *testing.T) {
		mock.ExpectQuery("SELECT sender_id, receiver_id, content, UNIX_TIMESTAMP").
			WithArgs("user1", "user2", "user2", "user1", 10).
			WillReturnError(fmt.Errorf("ошибка выполнения запроса"))

		messages, err := dbWrapper.GetHistory(context.Background(), "user1", "user2", 10)
		if err == nil {
			t.Errorf("ожидалась ошибка, но получена nil")
		}
		if messages == nil {
			t.Errorf("ожидался пустой срез, а не nil")
		}
	})

	t.Run("пустая история", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"sender_id", "receiver_id", "content", "timestamp"})

		mock.ExpectQuery("SELECT sender_id, receiver_id, content, UNIX_TIMESTAMP").
			WithArgs("user1", "user2", "user2", "user1", 10).
			WillReturnRows(rows)

		messages, err := dbWrapper.GetHistory(context.Background(), "user1", "user2", 10)
		if err != nil {
			t.Errorf("ошибка не ожидалась, но получена: %s", err)
		}

		if len(messages) != 0 {
			t.Errorf("ожидался пустой срез, получено сообщений: %d", len(messages))
		}
	})
}

func TestNewDB(t *testing.T) {
	t.Skip("Интеграционный тест: требует реального MySQL. Раскомментируйте для проверки, но учтите возможный таймаут до 90 сек при ошибке подключения.")
	origUser := os.Getenv("DB_USER")
	origPass := os.Getenv("DB_PASSWORD")
	origHost := os.Getenv("DB_HOST")
	origPort := os.Getenv("DB_PORT")
	origName := os.Getenv("DB_NAME")

	defer func() {
		os.Setenv("DB_USER", origUser)
		os.Setenv("DB_PASSWORD", origPass)
		os.Setenv("DB_HOST", origHost)
		os.Setenv("DB_PORT", origPort)
		os.Setenv("DB_NAME", origName)
	}()

	os.Setenv("DB_USER", "baduser")
	os.Setenv("DB_PASSWORD", "badpass")
	os.Setenv("DB_HOST", "127.0.0.1")
	os.Setenv("DB_PORT", "9999")
	os.Setenv("DB_NAME", "baddb")

	_, err := NewDB()
	if err == nil {
		t.Errorf("ожидалась ошибка подключения при некорректных данных, но получена nil")
	}
}
