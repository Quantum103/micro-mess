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

	_, err := database.NewDB()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer database.GetDB().Close()
	r := mux.NewRouter()
	db := database.GetDB()
	// Get запрос - userHandler.go
	r.HandleFunc("/dashboard", handlers.DashboardHandler(db))

	// POST запрос - postHandler.go
	r.HandleFunc("/api/posts", handlers.PostHandler(db))

	// POST запросы для смены
	r.HandleFunc("/change/username", handlers.ChangeUsername(db))
	r.HandleFunc("/change/job", handlers.UpdateWork(db))
	r.HandleFunc("/change/city", handlers.UpdateGEO(db))
	// r.HandleFunc("/change/Pass", handlers.ChangePass(db))

	log.Println("User Service запущен на порту 8082")
	http.ListenAndServe(":8082", r)
}
