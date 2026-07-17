package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Online   bool   `json:"online"`
}
type OnlineChecker interface {
	IsOnline(userID string) bool
}

func GetFriendsHandler(db *sql.DB, hub OnlineChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("X-User-ID")
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		query := `
		SELECT u.id, u.username
		FROM users u
		JOIN friends f 
			ON (f.friend_id = u.id OR f.user_id = u.id)
		WHERE 
			(f.user_id = ? OR f.friend_id = ?)
			AND f.status = 'accepted'
			AND u.id != ?
		`

		rows, err := db.Query(query, userID, userID, userID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		var users []User

		for rows.Next() {
			var u User

			if err := rows.Scan(&u.ID, &u.Username); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			u.Online = hub.IsOnline(u.ID)
			users = append(users, u)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		if users == nil {
			users = []User{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
