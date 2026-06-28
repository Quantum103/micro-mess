package chat

import (
	"chat-service/models"
	"context"
	"testing"
	"time"
)

func TestHubRegister(t *testing.T) {
	hub := NewHub(nil)

	go hub.Run()

	client := &Client{
		UserID: "user1",
		send:   make(chan models.ChatMessage, 1),
	}

	hub.register <- client

	time.Sleep(50 * time.Millisecond)

	if _, ok := hub.clients["user1"]; !ok {
		t.Fatal("client not registered")
	}
}
func TestHubUnregister(t *testing.T) {
	hub := NewHub(nil)

	client := &Client{
		UserID: "user1",
		send:   make(chan models.ChatMessage, 1),
	}

	hub.clients["user1"] = client

	go hub.Run()

	hub.unregister <- client

	time.Sleep(50 * time.Millisecond)

	if _, ok := hub.clients["user1"]; ok {
		t.Fatal("client not removed")
	}
}
func TestHubBroadcast(t *testing.T) {
	hub := NewHub(nil)

	sender := &Client{
		UserID: "user1",
		send:   make(chan models.ChatMessage, 1),
	}

	receiver := &Client{
		UserID: "user2",
		send:   make(chan models.ChatMessage, 1),
	}

	hub.clients["user1"] = sender
	hub.clients["user2"] = receiver

	go hub.Run()

	msg := models.ChatMessage{
		From:    "user1",
		To:      "user2",
		Content: "hello",
		Type:    "text",
	}

	hub.broadcast <- msg

	select {
	case received := <-receiver.send:
		if received.Content != "hello" {
			t.Fatal("wrong message")
		}
	case <-time.After(time.Second):
		t.Fatal("message not received")
	}
}

type MockDB struct {
	called bool
}

func (m *MockDB) SaveMessage(ctx context.Context, from, to, content string) error {
	m.called = true
	return nil
}

func TestHubSaveMessage(t *testing.T) {
	mockDB := &MockDB{}
	hub := NewHub(mockDB)

	go hub.Run()

	hub.broadcast <- models.ChatMessage{
		From:    "1",
		To:      "2",
		Content: "test",
		Type:    "text",
	}

	time.Sleep(100 * time.Millisecond)

	if !mockDB.called {
		t.Fatal("SaveMessage not called")
	}
}
