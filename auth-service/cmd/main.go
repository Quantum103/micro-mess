package main

import (
	"log"
	"net/http"

	"auth-service/database"
	"auth-service/handlers"
	"auth-service/repository"
	"auth-service/service"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	// repository layer
	userRepo := repository.NewUserRepository(db)

	// service layer
	authService := service.NewAuthService(userRepo)

	// handler layer
	authHandler := handlers.NewAuthHandler(authService)
	r.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	log.Println(" User Service запущен на порту 8081")
	log.Fatal(http.ListenAndServe(":8081", r))

}
