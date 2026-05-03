package routers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func RegisterStatic(router *mux.Router) {
	// 1. Статика (картинки, стили, скрипты)
	// Запрос /static/css/style.css → файл /app/frontend/static/css/style.css
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

	// 3. Fallback (на всякий случай)
	// Если маршрут не найден выше, пробуем найти файл с .html
	// Это подстраховка, основные маршруты уже прописаны в auth.go и user.go
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
