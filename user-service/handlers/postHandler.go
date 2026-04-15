package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func PostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getPost(w, r, db)
		} else if r.Method == http.MethodPost {
			createdPost(w, r, db)
		}
	}
}

func getPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	row, err := db.Query(`SELECT id, user_id, text, updated_at 
	FROM posts
	ORDER BY created_at DESC
	LIMIT 50`)
	if err != nil {
		log.Printf("ошибка в БД: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database error"})
		return
	}

	defer row.Close()
	var posts []map[string]interface{}
	for row.Next() {
		var id uint64
		var userID uint64
		var text string
		var createdAt time.Time

		err := row.Scan(&id, &userID, &text, &createdAt)
		if err != nil {
			continue
		}
		posts = append(posts, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"text":       text,
			"created_at": createdAt,
		})
	}
	if posts == nil {
		posts = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)

}

type CreatePostRequest struct {
	Text string `json:"text"`
}

type CreatePostResponse struct {
	Success bool  `json:"success"`
	PostID  int64 `json:"post_id"`
}

func createdPost(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	UserIDstr := r.Header.Get("X-User-ID")
	UserID, err := strconv.ParseInt(UserIDstr, 10, 64)
	if err != nil {
		http.Error(w, "UserID не число", http.StatusBadRequest)
		return
	}

	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "ошибка JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`INSERT INTO posts (user_id, text, created_at) VALUES (?,?,?)`, UserID, req.Text, time.Now())
	if err != nil {
		http.Error(w, "ошибка сохранение в БД", http.StatusBadRequest)
		return
	}

	postID, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreatePostResponse{
		Success: true,
		PostID:  postID,
	})
}
