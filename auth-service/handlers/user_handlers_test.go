package handlers

import (
	"auth-service/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockAuthService struct{}

func (m *MockAuthService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	return &models.UserResponse{
		ID:       1,
		Username: req.Username,
		Email:    req.Email,
	}, nil
}

func (m *MockAuthService) Login(req *models.UserLogin) (string, error) {
	return "fake-token", nil
}

func TestRegister_Success(t *testing.T) {
	mockService := &MockAuthService{}
	handler := &AuthHandler{
		AuthService: mockService,
	}

	body, _ := json.Marshal(models.RegisterRequest{
		Username: "anton",
		Email:    "anton@mail.com",
		Password: "123456",
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRegister_NoSuccses(t *testing.T) {
	mockService := &MockAuthService{}
	handler := &AuthHandler{
		AuthService: mockService,
	}
	body, _ := json.Marshal(models.RegisterRequest{
		Username: "",
		Email:    "antonmail.com",
		Password: "12",
	})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	require.Contains(t, rr.Body.String(), "пользователь создан")
}

func TestLogin_Success(t *testing.T) {
	mockService := &MockAuthService{}
	handler := &AuthHandler{
		AuthService: mockService,
	}

	body, _ := json.Marshal(models.UserLogin{
		Identifier: "anton@mail.com",
		Password:   "123456",
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Contains(t, rr.Body.String(), "fake-token")

	cookies := rr.Result().Cookies()
	require.Len(t, cookies, 1)
	require.Equal(t, "auth_token", cookies[0].Name)
}
func TestLogin_NoSuccess(t *testing.T) {
	mockService := &MockAuthService{}
	handler := &AuthHandler{
		AuthService: mockService,
	}

	body, _ := json.Marshal(models.UserLogin{
		Identifier: ".com",
		Password:   "16",
	})

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Contains(t, rr.Body.String(), "fake-token")

	cookies := rr.Result().Cookies()
	require.Len(t, cookies, 1)
	require.Equal(t, "auth_token", cookies[0].Name)
}
