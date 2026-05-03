package routers

import (
	"api-gateway/middleware"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

func RegisterChat(mux *mux.Router, proxy *httputil.ReverseProxy) {
	mux.HandleFunc("/chat", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/chat.html")
	})).Methods("GET")

	mux.HandleFunc("/api/users", middleware.AuthMiddleware(proxy.ServeHTTP)).Methods("GET")
	mux.HandleFunc("/api/history", middleware.AuthMiddleware(proxy.ServeHTTP)).Methods("GET")

	mux.HandleFunc("/ws", middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}))
}
