package routers

import (
	"api-gateway/middleware"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

func RegisterUser(router *mux.Router, proxy *httputil.ReverseProxy) {

	router.HandleFunc("/dashboard", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/dashboard.html")
	})).Methods("GET")

	router.HandleFunc("/settings", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/settings.html")
	})).Methods("GET")

	router.HandleFunc("/people", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/people.html")
	})).Methods("GET")

	router.HandleFunc("/friends", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/friends.html")
	})).Methods("GET")

	// API

	router.HandleFunc("/api/dashboard", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/posts", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}))

	router.HandleFunc("/api/people", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/change/username", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})).Methods("POST")
	router.HandleFunc("/change/city", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})).Methods("POST")
	router.HandleFunc("/change/work", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})).Methods("POST")
	router.HandleFunc("/change/Pass", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})).Methods("POST")
	router.HandleFunc("/api/friends", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})).Methods("GET", "POST")
	router.HandleFunc("/api/friends/accept", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}))
}
