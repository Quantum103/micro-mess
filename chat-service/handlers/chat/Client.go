package chat

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"chat-service/models"

	"github.com/gorilla/websocket"
)

const (
	WriteWait      = 10 * time.Second
	PongWait       = 60 * time.Second
	PingPeriod     = (PongWait * 9) / 10
	MaxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // В продакшене ограничьте домены!
	},
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan models.ChatMessage
	UserID string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(PongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg models.ChatMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("JSON parse error:", err)
			continue
		}

		// Валидация
		if msg.To == "" {
			log.Println("Missing recipient (To)")
			continue
		}
		if msg.Content == "" {
			continue
		}

		// Очистка контента
		msg.From = c.UserID
		msg.Timestamp = time.Now().Unix()
		msg.Type = "text"
		msg.Content = string(bytes.TrimSpace(bytes.Replace([]byte(msg.Content), []byte{'\n'}, []byte{' '}, -1)))

		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Println("JSON marshal error:", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("🔥 WS HANDLER HIT")
	// Gateway прокидывает user_id из JWT в заголовке
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "Missing X-User-ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan models.ChatMessage, 256),
		UserID: userID,
	}

	client.hub.register <- client

	// Запускаем горутины
	go client.writePump()
	go client.readPump()
}
