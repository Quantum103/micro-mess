package chat

import (
	"context"
	"log"
	"time"

	"chat-service/models"
)

type Database interface {
	SaveMessage(ctx context.Context, from, to, content string) error
}

type Hub struct {
	clients    map[string]*Client
	broadcast  chan models.ChatMessage
	register   chan *Client
	unregister chan *Client
	db         Database
}

func NewHub(db Database) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan models.ChatMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		db:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.UserID] = client
			log.Printf("User %s connected", client.UserID)

		case client := <-h.unregister:
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)
				log.Printf("User %s disconnected", client.UserID)
			}

		case msg := <-h.broadcast:
			// Сохраняем в БД асинхронно
			if h.db != nil && msg.Type == "text" && msg.To != "" {
				go func(m models.ChatMessage) {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					defer cancel()
					if err := h.db.SaveMessage(ctx, m.From, m.To, m.Content); err != nil {
						log.Println("DB save error:", err)
					}
				}(msg)
			}

			// Отправляем получателю (если онлайн)
			if recipient, ok := h.clients[msg.To]; ok {
				select {
				case recipient.send <- msg:
				default:
					close(recipient.send)
					delete(h.clients, msg.To)
				}
			}

			// Отправляем отправителю
			if sender, ok := h.clients[msg.From]; ok {
				select {
				case sender.send <- msg:
				default:
					close(sender.send)
					delete(h.clients, msg.From)
				}
			}
		}
	}
}

func (h *Hub) IsOnline(userID string) bool {
	_, ok := h.clients[userID]
	return ok
}
