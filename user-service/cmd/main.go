package main

import (
	"log"
	"net/http"
	"user-service/database"
	"user-service/handlers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer db.Close()

	userRepo := database.NewUserRepository(db)

	r := mux.NewRouter()
	// Get запрос - userHandler.go
	r.HandleFunc("/dashboard", handlers.DashboardHandler(db))
	r.HandleFunc("/api/dashboard", handlers.DashboardHandler(db))

	// POST запрос - postHandler.go
	r.HandleFunc("/api/posts", handlers.PostHandler(db))

	// POST запросы для смены настроек
	r.HandleFunc("/change/username", handlers.ChangeUsername(userRepo))
	r.HandleFunc("/change/work", handlers.UpdateWork(userRepo))
	r.HandleFunc("/change/city", handlers.UpdateGEO(userRepo))
	r.HandleFunc("/change/Pass", handlers.UpdatePassword(userRepo))

	// маршруты для друзей
	r.HandleFunc("/people", handlers.ListUserHandler(userRepo))
	r.HandleFunc("/api/people", handlers.ListUserHandler(userRepo))

	r.HandleFunc("/api/friends", handlers.GetFrineds(userRepo)).Methods("GET")
	r.HandleFunc("/api/friends", handlers.AddFriendHandler(userRepo)).Methods("POST")
	r.HandleFunc("/api/friends/accept", handlers.AcceptFriendHandler(userRepo)).Methods("POST")

	log.Println("User Service запущен на порту 8082")
	http.ListenAndServe(":8082", r)
}
