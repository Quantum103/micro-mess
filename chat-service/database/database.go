package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"chat-service/models"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 30; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("DB connection failed: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DB{db}, nil
}

// SaveMessage сохраняет сообщение (реализует интерфейс из Hub)
func (db *DB) SaveMessage(ctx context.Context, from, to, content string) error {
	query := `INSERT INTO messages (sender_id, receiver_id, content, created_at) VALUES (?, ?, ?, NOW())`
	_, err := db.ExecContext(ctx, query, from, to, content)
	return err
}

// GetHistory загружает историю
func (db *DB) GetHistory(ctx context.Context, userID, partnerID string, limit int) ([]models.ChatMessage, error) {
	query := `SELECT sender_id, receiver_id, content, UNIX_TIMESTAMP(created_at) as timestamp 
			  FROM messages 
			  WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?) 
			  ORDER BY created_at ASC LIMIT ?`

	rows, err := db.QueryContext(ctx, query, userID, partnerID, partnerID, userID, limit)
	if err != nil {
		return []models.ChatMessage{}, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var msg models.ChatMessage
		if err := rows.Scan(&msg.From, &msg.To, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		msg.Type = "text"
		messages = append(messages, msg)
	}
	return messages, nil
}
