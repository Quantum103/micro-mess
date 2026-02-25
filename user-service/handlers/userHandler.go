// функции для GET запросов - показать
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

var db *sql.DB

func DashboardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-User-ID")
		email := r.Header.Get("X-User-Email")

		var (
			username string
			location string
			birthday sql.NullString
			work     string
		)

		err := db.QueryRow(`
                SELECT COALESCE(username, '') as username, COALESCE(location, '') as location, 
                   COALESCE(birthday, '') as birthday, COALESCE(work, '') as work
            FROM users WHERE id = ?`, id).
			Scan(&username, &location, &birthday, &work)
		if err != nil {
			http.Error(w, "пользователь не найден", http.StatusNotFound)
			return
		}
		response := map[string]interface{}{
			"id":             id,
			"email":          email,
			"username":       username,
			"location":       location,
			"birthday":       birthday.String,
			"work":           work,
			"friendsCount":   0,
			"postsCount":     0,
			"followersCount": 0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
