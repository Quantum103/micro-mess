// middleware and proxy находится в папке middleware
package main

import (
	"api-gateway/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	authProxy := middleware.CreateProxy("auth-service:8081")
	userProxy := middleware.CreateProxy("user-service:8082")
<<<<<<< HEAD
	chatProxy := middleware.CreateProxy("chat-service:8083")
=======
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56

	r.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		authProxy.ServeHTTP(w, r)
	}).Methods("POST")

	r.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		authProxy.ServeHTTP(w, r)
	}).Methods("POST")

	r.HandleFunc("/dashboard", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("GET")

	r.HandleFunc("/api/posts", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	}))
	r.HandleFunc("/settings", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	}))
	r.HandleFunc("/change/username", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("POST")
	r.HandleFunc("/change/city", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("POST")
	r.HandleFunc("/change/work", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("POST")
	r.HandleFunc("/change/Pass", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("POST")

<<<<<<< HEAD
	r.HandleFunc("/people", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("GET")
	r.HandleFunc("/friends", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("GET")
	r.HandleFunc("/friends", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("POST")
	r.HandleFunc("/friends/accept", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	}))

	r.HandleFunc("/chat", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/chat.html")
	})).Methods("GET")

	r.HandleFunc("/api/users", middleware.AuthMiddleware(chatProxy.ServeHTTP)).Methods("GET")
	r.HandleFunc("/api/history", middleware.AuthMiddleware(chatProxy.ServeHTTP)).Methods("GET")

	r.HandleFunc("/ws", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		chatProxy.ServeHTTP(w, r)
	}))

=======
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("/app/frontend"))))

	log.Fatal(http.ListenAndServe(":8080", r))
}
