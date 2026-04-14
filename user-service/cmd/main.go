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
<<<<<<< HEAD
	defer db.Close()

	userRepo := database.NewUserRepository(db)

=======
	defer database.GetDB().Close()
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
	r := mux.NewRouter()
	// Get запрос - userHandler.go
	r.HandleFunc("/dashboard", handlers.DashboardHandler(db))

	// POST запрос - postHandler.go
	r.HandleFunc("/api/posts", handlers.PostHandler(db))

	// POST запросы для смены настроек
<<<<<<< HEAD
	r.HandleFunc("/change/username", handlers.ChangeUsername(userRepo))
	r.HandleFunc("/change/work", handlers.UpdateWork(userRepo))
	r.HandleFunc("/change/city", handlers.UpdateGEO(userRepo))
	r.HandleFunc("/change/Pass", handlers.UpdatePassword(userRepo))

	// маршруты для друзей
	r.HandleFunc("/people", handlers.ListUserHandler(userRepo)).Methods("GET")
	r.HandleFunc("/friends", handlers.GetFrineds(userRepo)).Methods("GET")
	r.HandleFunc("/friends", handlers.AddFriendHandler(userRepo)).Methods("POST")
	r.HandleFunc("/friends/accept", handlers.AcceptFriendHandler(userRepo)).Methods("POST")
=======
	r.HandleFunc("/change/username", handlers.ChangeUsername)
	r.HandleFunc("/change/work", handlers.UpdateWork)
	r.HandleFunc("/change/city", handlers.UpdateGEO)
	r.HandleFunc("/change/Pass", handlers.UpdatePassword)
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56

	log.Println("User Service запущен на порту 8082")
	http.ListenAndServe(":8082", r)
}
