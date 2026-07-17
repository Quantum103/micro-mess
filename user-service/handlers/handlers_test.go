package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"user-service/database"

	"github.com/DATA-DOG/go-sqlmock"
)

func toJSON(v interface{}) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

type MockUserRepo struct{}

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func setupWeatherMock(resp *http.Response, err error) func() {
	originalTransport := http.DefaultTransport
	http.DefaultTransport = &mockTransport{
		roundTripFunc: func(req *http.Request) (*http.Response, error) {
			return resp, err
		},
	}
	return func() {
		http.DefaultTransport = originalTransport
	}
}

func setupMockRepo(t *testing.T) (*database.UserRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating mock: %s", err)
	}

	repo := database.NewUserRepository(db)

	cleanup := func() {
		db.Close()
	}
	return repo, mock, cleanup
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected int
	}{
		{"пустой заголовок", "", 0},
		{"не число", "abc", 0},
		{"валидный ID", "123", 123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				req.Header.Set("X-User-ID", tt.header)
			}
			if got := GetUserID(req); got != tt.expected {
				t.Errorf("GetUserID() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestChangeUsername(t *testing.T) {
	repo, mock, cleanup := setupMockRepo(t)
	defer cleanup()
	handler := ChangeUsername(repo)

	t.Run("не авторизован", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", toJSON(map[string]string{"newName": "test"}))
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", rr.Code)
		}
	})

	t.Run("неверный JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
		req.Header.Set("X-User-ID", "1")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("слишком короткое имя", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", toJSON(map[string]string{"newName": "a"}))
		req.Header.Set("X-User-ID", "1")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("успешное обновление", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET username = \\? WHERE id = \\?").
			WithArgs("ValidName", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req := httptest.NewRequest(http.MethodPost, "/", toJSON(map[string]string{"newName": "ValidName"}))
		req.Header.Set("X-User-ID", "1")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestUpdateGEO(t *testing.T) {
	repo, mock, cleanup := setupMockRepo(t)
	defer cleanup()
	handler := UpdateGEO(repo)

	t.Run("пустой город", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", toJSON(map[string]string{"city": ""}))
		req.Header.Set("X-User-ID", "1")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("успешное обновление", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET location = \\? WHERE id = \\?").
			WithArgs("Москва", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req := httptest.NewRequest(http.MethodPost, "/", toJSON(map[string]string{"city": "Москва"}))
		req.Header.Set("X-User-ID", "1")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET location = \\? WHERE id = \\?").
			WithArgs("Москва", 1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		req := httptest.NewRequest(http.MethodPost, "/", toJSON(map[string]string{"city": "Москва"}))
		req.Header.Set("X-User-ID", "1")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
	})
}

func TestUpdateWork(t *testing.T) {
	repo, mock, cleanup := setupMockRepo(t)
	defer cleanup()
	handler := UpdateWork(repo)

	t.Run("успешное обновление", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET work = \\? WHERE id = \\?").
			WithArgs("Google", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req := httptest.NewRequest(http.MethodPost, "/", toJSON(map[string]string{"work_location": "Google"}))
		req.Header.Set("X-User-ID", "1")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestUpdatePassword(t *testing.T) {
	repo, mock, cleanup := setupMockRepo(t)
	defer cleanup()
	handler := UpdatePassword(repo)

	t.Run("пустые поля", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", toJSON(map[string]string{"OldPass": "", "NewPass": "new"}))
		req.Header.Set("X-User-ID", "1")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("успешное обновление", func(t *testing.T) {
		mock.ExpectQuery("SELECT password FROM users WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("$2a$10$...")) // фейковый хеш

		mock.ExpectExec("UPDATE users SET password = \\? WHERE id = \\?").
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		req := httptest.NewRequest(http.MethodPost, "/", toJSON(map[string]string{"OldPass": "oldpass", "NewPass": "newpass"}))
		req.Header.Set("X-User-ID", "1")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestDashboardHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating mock: %s", err)
	}
	defer db.Close()

	handler := DashboardHandler(db)

	origKey := os.Getenv("WEATHER_API_KEY")
	defer os.Setenv("WEATHER_API_KEY", origKey)
	os.Setenv("WEATHER_API_KEY", "test-key")

	t.Run("невалидный ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "abc")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("пользователь не найден", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "1")
		req.Header.Set("X-User-Email", "test@test.com")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("успешный запрос с погодой", func(t *testing.T) {
		mock.ExpectQuery("SELECT COALESCE").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"username", "location", "birthday", "work"}).
				AddRow("ivan", "Moscow", "1990-01-01", "Developer"))

		mockWeatherResp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"main": {"temp": 20}, "weather": [{"main": "Clear"}]}`)),
		}
		mockWeatherResp.Header.Set("Content-Type", "application/json")
		defer setupWeatherMock(mockWeatherResp, nil)()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-User-ID", "1")
		req.Header.Set("X-User-Email", "test@test.com")
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)

		if response["username"] != "ivan" {
			t.Errorf("expected username 'ivan', got %v", response["username"])
		}
		if response["weather"] == nil {
			t.Errorf("expected weather data to be present")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

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
