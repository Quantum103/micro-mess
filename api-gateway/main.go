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

	authProxy := middleware.CreateProxy("localhost:8081")
	userProxy := middleware.CreateProxy("localhost:8082")

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
	r.HandleFunc("/changesettings", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	}))
	r.HandleFunc("/change/username", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("POST")
	r.HandleFunc("/change/city", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("POST")
	r.HandleFunc("/change/job", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userProxy.ServeHTTP(w, r)
	})).Methods("POST")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("../frontend/"))))

	log.Fatal(http.ListenAndServe(":8080", r))
}
