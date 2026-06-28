package chat

import (
	"bytes"
	"chat-service/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func sanitizeContent(content string) string {
	return string(bytes.TrimSpace(
		bytes.Replace([]byte(content), []byte{'\n'}, []byte{' '}, -1),
	))
}
func TestSanitizeContent(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{" hello ", "hello"},
		{"hello\nworld", "hello world"},
		{"   hello\nworld   ", "hello world"},
		{"", ""},
	}

	for _, tt := range tests {
		got := sanitizeContent(tt.input)

		if got != tt.expected {
			t.Fatalf("expected %q, got %q", tt.expected, got)
		}
	}
}
func TestClientReadPump_ValidMessage(t *testing.T) {
	hub := NewHub(nil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan models.ChatMessage, 1),
			UserID: "user1",
		}

		go client.readPump()
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = conn.WriteJSON(models.ChatMessage{
		To:      "user2",
		Content: "hello",
	})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case msg := <-hub.broadcast:
		if msg.From != "user1" {
			t.Fatal("wrong sender")
		}
		if msg.To != "user2" {
			t.Fatal("wrong recipient")
		}
		if msg.Content != "hello" {
			t.Fatal("wrong content")
		}
		if msg.Type != "text" {
			t.Fatal("wrong type")
		}
	case <-time.After(time.Second):
		t.Fatal("message not broadcasted")
	}
}
func TestClientReadPump_InvalidJSON(t *testing.T) {
	hub := NewHub(nil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan models.ChatMessage, 1),
			UserID: "user1",
		}

		go client.readPump()
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	conn.WriteMessage(websocket.TextMessage, []byte("{invalid json"))

	select {
	case <-hub.broadcast:
		t.Fatal("invalid JSON should not broadcast")
	case <-time.After(300 * time.Millisecond):
	}
}

func TestClientWritePump(t *testing.T) {
	hub := NewHub(nil)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan models.ChatMessage, 1),
			UserID: "user1",
		}

		go client.writePump()

		client.send <- models.ChatMessage{
			From:    "user1",
			To:      "user2",
			Content: "hello",
		}
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	var msg models.ChatMessage
	err = conn.ReadJSON(&msg)
	if err != nil {
		t.Fatal(err)
	}

	if msg.Content != "hello" {
		t.Fatal("wrong content")
	}
}
