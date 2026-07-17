package routers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func RegisterStatic(router *mux.Router) {
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("/app/frontend/static"))),
	)

	// 2. Главная страница
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "/app/frontend/index.html")
	}).Methods("GET")

	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/ws" {
			http.NotFound(w, r)
			return
		}

		path := r.URL.Path
		if filepath.Ext(path) == "" && path != "/" {
			tryPath := filepath.Join("/app/frontend", path+".html")
			if _, err := os.Stat(tryPath); err == nil {
				http.ServeFile(w, r, tryPath)
				return
			}
		}
		http.StripPrefix("/", http.FileServer(http.Dir("/app/frontend"))).ServeHTTP(w, r)
	})
}
