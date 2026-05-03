// middleware and proxy находится в папке middleware
package main

import (
	"api-gateway/middleware"
	"api-gateway/routers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	authProxy := middleware.CreateProxy("auth-service:8081")
	userProxy := middleware.CreateProxy("user-service:8082")
	chatProxy := middleware.CreateProxy("chat-service:8083")

	routers.RegisterAuth(r, authProxy)
	routers.RegisterUser(r, userProxy)
	routers.RegisterChat(r, chatProxy)
	routers.RegisterStatic(r)

	log.Println("🚀 Gateway запущен на порту :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
