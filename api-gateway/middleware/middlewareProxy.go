package middleware

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("my-super-secret-key-12345")

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			if cookie, err := r.Cookie("auth_token"); err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			http.Redirect(w, r, "/login.html", http.StatusSeeOther)
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Неверный или просроченный токен", http.StatusUnauthorized)
			return
		}

		if userID, ok := claims["user_id"].(float64); ok {
			r.Header.Set("X-User-ID", fmt.Sprintf("%.0f", userID))
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
	target, _ := url.Parse("http://" + host)
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Ошибка прокси: %v", err)
		http.Error(w, "Сервис временно недоступен", http.StatusServiceUnavailable)
	}

	return proxy
}
