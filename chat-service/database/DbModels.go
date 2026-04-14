package database

import "time"

type Message struct {
	ID         uint      `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	isRead     bool      `json:"is_read"`
}
