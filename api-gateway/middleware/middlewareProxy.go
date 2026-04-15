package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

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
			if r.Header.Get("Accept") == "application/json" || r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "Пользователь не авторизован"})
				return
			}
			// Для обычных страниц (браузер) — редирект на логин
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
	targetURL, err := url.Parse("http://" + host)
	if err != nil {
		log.Fatalf(" Ошибка парсинга прокси %s: %v", host, err)
	}

	// Создаём прокси
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	defaultDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		defaultDirector(req)
		req.Host = targetURL.Host
		if req.Header.Get("Upgrade") != "" {
			req.Header.Set("Connection", "Upgrade")
			req.Header.Set("Upgrade", req.Header.Get("Upgrade"))
		}
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf(" Ошибка прокси для %s: %v", host, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"error": "service unavailable ыыы"})
	}

	proxy.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100,
		DisableCompression:    true,
	}

	return proxy
}
