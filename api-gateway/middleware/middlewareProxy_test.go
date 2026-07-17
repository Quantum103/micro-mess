package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

func TestAuthMiddleware(t *testing.T) {
	var capturedReq *http.Request

	nextHandler := func(w http.ResponseWriter, r *http.Request) {
		capturedReq = r
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	handler := AuthMiddleware(nextHandler)

	t.Run("no token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rr.Code)
		}
		if rr.Body.String() != "не авториз\n" {
			t.Errorf("expected body 'не авториз\\n', got '%s'", rr.Body.String())
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?token=invalid.token.here", nil)
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rr.Code)
		}
		if rr.Body.String() != "invalid token\n" {
			t.Errorf("expected body 'invalid token\\n', got '%s'", rr.Body.String())
		}
	})

	t.Run("valid token in query", func(t *testing.T) {
		token := generateTestToken(jwt.MapClaims{
			"user_id":  123,
			"email":    "test@example.com",
			"username": "testuser",
		})
		req := httptest.NewRequest(http.MethodGet, "/?token="+token, nil)
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
		if capturedReq.Header.Get("X-User-ID") != "123" {
			t.Errorf("expected X-User-ID '123', got '%s'", capturedReq.Header.Get("X-User-ID"))
		}
		if capturedReq.Header.Get("X-User-Email") != "test@example.com" {
			t.Errorf("expected X-User-Email 'test@example.com', got '%s'", capturedReq.Header.Get("X-User-Email"))
		}
		if capturedReq.Header.Get("X-User-Username") != "testuser" {
			t.Errorf("expected X-User-Username 'testuser', got '%s'", capturedReq.Header.Get("X-User-Username"))
		}
	})

	t.Run("valid token in authorization header with bearer", func(t *testing.T) {
		token := generateTestToken(jwt.MapClaims{
			"user_id": 456,
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
		if capturedReq.Header.Get("X-User-ID") != "456" {
			t.Errorf("expected X-User-ID '456', got '%s'", capturedReq.Header.Get("X-User-ID"))
		}
	})

	t.Run("valid token in cookie", func(t *testing.T) {
		token := generateTestToken(jwt.MapClaims{
			"username": "cookieuser",
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})
		rr := httptest.NewRecorder()
		handler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}
		if capturedReq.Header.Get("X-User-Username") != "cookieuser" {
			t.Errorf("expected X-User-Username 'cookieuser', got '%s'", capturedReq.Header.Get("X-User-Username"))
		}
	})
}

func TestCreateProxy(t *testing.T) {
	t.Run("proxy creation and director", func(t *testing.T) {
		proxy := CreateProxy("backend.local:8080")
		if proxy == nil {
			t.Fatal("expected proxy to be created, got nil")
		}

		req := httptest.NewRequest(http.MethodGet, "http://frontend.local/path", nil)
		req.Header.Set("X-User-ID", "999")
		req.Header.Set("X-User-Email", "proxy@example.com")
		req.Header.Set("X-Other-Header", "keep-me")

		proxy.Director(req)

		if req.Host != "backend.local:8080" {
			t.Errorf("expected Host 'backend.local:8080', got '%s'", req.Host)
		}
		if req.Header.Get("X-User-ID") != "999" {
			t.Errorf("expected X-User-ID '999', got '%s'", req.Header.Get("X-User-ID"))
		}
		if req.Header.Get("X-User-Email") != "proxy@example.com" {
			t.Errorf("expected X-User-Email 'proxy@example.com', got '%s'", req.Header.Get("X-User-Email"))
		}
	})

	t.Run("proxy error handler", func(t *testing.T) {
		proxy := CreateProxy("backend.local:8080")

		proxy.Transport = &errorTransport{}

		req := httptest.NewRequest(http.MethodGet, "http://frontend.local/path", nil)
		rr := httptest.NewRecorder()

		proxy.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", rr.Code)
		}
		if rr.Body.String() != "gateway error\n" {
			t.Errorf("expected body 'gateway error\\n', got '%s'", rr.Body.String())
		}
	})
}

type errorTransport struct{}

func (t *errorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("simulated proxy transport error")
}
