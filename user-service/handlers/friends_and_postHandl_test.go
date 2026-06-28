package handlers

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-service/database"

	"github.com/DATA-DOG/go-sqlmock"
)

type MockUserRepo struct{}

func TestPostHandler_GetPostsSuccess(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "user_id", "text", "updated_at"}).
		AddRow(1, 10, "hello", time.Now())

	mock.ExpectQuery("SELECT id, user_id, text, updated_at").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	handler := PostHandler(db)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}
}

func TestPostHandler_GetPostsEmpty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "user_id", "text", "updated_at"})

	mock.ExpectQuery("SELECT id, user_id, text, updated_at").
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	handler := PostHandler(db)
	handler(w, req)

	var posts []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &posts)

	if len(posts) != 0 {
		t.Fatal("expected empty posts")
	}
}

func TestPostHandler_GetPostsDBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery("SELECT id, user_id, text, updated_at").
		WillReturnError(assertAnError())

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	w := httptest.NewRecorder()

	handler := PostHandler(db)
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", w.Code)
	}
}

func TestCreatedPostSuccess(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("INSERT INTO posts").
		WithArgs(
			int64(1),
			"test post",
			sqlmock.AnyArg(),
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	body := []byte(`{"text":"test post"}`)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	handler := PostHandler(db)
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", w.Code)
	}

	var resp CreatePostResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if !resp.Success {
		t.Fatal("expected success true")
	}
}

func TestCreatedPostInvalidUserID(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	body := []byte(`{"text":"test post"}`)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "abc")

	w := httptest.NewRecorder()

	handler := PostHandler(db)
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestCreatedPostInvalidJSON(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer([]byte(`bad json`)))
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	handler := PostHandler(db)
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestCreatedPostDBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("INSERT INTO posts").
		WillReturnError(assertAnError())

	body := []byte(`{"text":"test post"}`)

	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")

	w := httptest.NewRecorder()

	handler := PostHandler(db)
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func assertAnError() error {
	return driver.ErrBadConn
}

func (m *MockUserRepo) ListAllUsers(ctx context.Context, userID int, limit int, offset int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"id": 2, "username": "alice"},
	}, nil
}

func (m *MockUserRepo) GetFrineds(ctx context.Context, userID int, limit int, offset int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{"id": 3, "username": "bob"},
	}, nil
}

func (m *MockUserRepo) AddFriend(ctx context.Context, userID int, friendID int) error {
	return nil
}

func (m *MockUserRepo) AcceptFriend(ctx context.Context, userID int, friendID int) error {
	return nil
}

func TestListUserHandler_Unauthorized(t *testing.T) {
	repo := &MockUserRepo{}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler := ListUserHandler((*database.UserRepository)(nil))
	handler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", w.Code)
	}

	_ = repo
}

func TestAddFriendHandler_BadJSON(t *testing.T) {
	handler := AddFriendHandler((*database.UserRepository)(nil))

	req := httptest.NewRequest(http.MethodPost, "/add-friend", bytes.NewBuffer([]byte("bad json")))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestAcceptFriendHandler_BadJSON(t *testing.T) {
	handler := AcceptFriendHandler((*database.UserRepository)(nil))

	req := httptest.NewRequest(http.MethodPost, "/accept-friend", bytes.NewBuffer([]byte("bad json")))
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", w.Code)
	}
}

func TestMockRepo_AddFriend(t *testing.T) {
	repo := &MockUserRepo{}

	err := repo.AddFriend(context.Background(), 1, 2)

	if err != nil {
		t.Fatal("expected nil")
	}
}

func TestMockRepo_AcceptFriend(t *testing.T) {
	repo := &MockUserRepo{}

	err := repo.AcceptFriend(context.Background(), 1, 2)

	if err != nil {
		t.Fatal("expected nil")
	}
}

func TestMockRepo_ErrorCase(t *testing.T) {
	err := errors.New("some error")

	if err.Error() != "some error" {
		t.Fatal("unexpected error")
	}
}
