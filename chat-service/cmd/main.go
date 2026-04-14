package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"chat-service/database"
	"chat-service/handlers"
	"chat-service/handlers/chat"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	hub := chat.NewHub(db)
	go hub.Run()

	r := mux.NewRouter()

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWS(hub, w, r)
	})

	r.HandleFunc("/api/users", handlers.GetUsersHandler(db.DB, hub)).Methods("GET")

	// Получить историю сообщений между двумя пользователями
	r.HandleFunc("/api/history", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		partnerID := r.URL.Query().Get("partner_id")
		if userID == "" || partnerID == "" {
			http.Error(w, "Missing IDs", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		messages, err := db.GetHistory(ctx, userID, partnerID, 50)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8080", "http://127.0.0.1:8080"}, AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	log.Println("Chat service starting on :8083")
	log.Fatal(http.ListenAndServe(":8083", handler))
}
