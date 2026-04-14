package main

import (
	"log"
	"net/http"

	"auth-service/database"
	"auth-service/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/api/register", handlers.HandleRegister(db)).Methods("POST")
	r.HandleFunc("/api/login", handlers.HandlerLogin(db)).Methods("POST")

	log.Println(" User Service запущен на порту 8081")
	log.Fatal(http.ListenAndServe(":8081", r))

}
