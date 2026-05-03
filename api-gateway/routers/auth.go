package routers

import (
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/mux"
)

func RegisterAuth(router *mux.Router, proxy *httputil.ReverseProxy) {
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/login.html")
	}).Methods("GET")

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/register.html")
	}).Methods("GET")

	//  API

	router.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}).Methods("POST")

	router.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}).Methods("POST")
}
