package handlers

import (
	"chat-service/handlers/chat"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func assertAnError() error {
	return errors.New("db error")
}

func TestGetFriendsHandler_Unauthorized(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	hub := chat.NewHub(nil)

	req := httptest.NewRequest("GET", "/friends", nil)
	w := httptest.NewRecorder()

	handler := GetFriendsHandler(db, hub)
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGetFriendsHandler_EmptyFriends(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	hub := chat.NewHub(nil)

	req := httptest.NewRequest("GET", "/friends", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	handler := GetFriendsHandler(db, hub)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var users []User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 0 {
		t.Fatalf("expected empty array")
	}
}

func TestGetFriendsHandler_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow("2", "alice").
		AddRow("3", "bob")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	hub := chat.NewHub(nil)

	req := httptest.NewRequest("GET", "/friends", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	handler := GetFriendsHandler(db, hub)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var users []User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if users[0].ID != "2" {
		t.Fatal("wrong first user id")
	}

	if users[1].Username != "bob" {
		t.Fatal("wrong second username")
	}
}

func TestGetFriendsHandler_OnlineStatus(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username"}).
		AddRow("2", "alice").
		AddRow("3", "bob")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	hub := chat.NewHub(nil)

	hub.AddClient(&chat.Client{
		UserID: "2",
	})

	req := httptest.NewRequest("GET", "/friends", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	handler := GetFriendsHandler(db, hub)
	handler(w, req)

	var users []User
	json.Unmarshal(w.Body.Bytes(), &users)

	if !users[0].Online {
		t.Fatal("alice should be online")
	}

	if users[1].Online {
		t.Fatal("bob should be offline")
	}
}

func TestGetFriendsHandler_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT").
		WillReturnError(assertAnError())

	hub := chat.NewHub(nil)

	req := httptest.NewRequest("GET", "/friends", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	handler := GetFriendsHandler(db, hub)
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetFriendsHandler_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"wrong_column"}).
		AddRow("oops")

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	hub := chat.NewHub(nil)

	req := httptest.NewRequest("GET", "/friends", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	handler := GetFriendsHandler(db, hub)
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetFriendsHandler_ContentType(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	hub := chat.NewHub(nil)

	req := httptest.NewRequest("GET", "/friends", nil)
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	handler := GetFriendsHandler(db, hub)
	handler(w, req)

	contentType := w.Header().Get("Content-Type")

	if contentType != "application/json" {
		t.Fatalf("expected application/json, got %s", contentType)
	}
}
