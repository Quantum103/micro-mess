package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("my-super-secret-key-12345")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.URL.Query().Get("token")

		if tokenString == "" {
			tokenString = r.Header.Get("Authorization")
		}

		if tokenString == "" {
			if cookie, err := r.Cookie("auth_token"); err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			http.Error(w, "не авториз", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		if userID, ok := claims["user_id"]; ok {
			r.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))
		}

		if email, ok := claims["email"].(string); ok {
			r.Header.Set("X-User-Email", email)
		}

		if username, ok := claims["username"].(string); ok {
			r.Header.Set("X-User-Username", username)
		}

		next(w, r)
	}
}

func CreateProxy(host string) *httputil.ReverseProxy {
	targetURL, err := url.Parse("http://" + host)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		req.Host = targetURL.Host

		if uid := req.Header.Get("X-User-ID"); uid != "" {
			req.Header.Set("X-User-ID", uid)
		}
		if email := req.Header.Get("X-User-Email"); email != "" {
			req.Header.Set("X-User-Email", email)
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Println("proxy error:", err)
		http.Error(w, "gateway error", 502)
	}

	return proxy
}
