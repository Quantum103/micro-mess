package handlers

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
)

func DashboardHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := r.Header.Get("X-User-ID")
        userID, err := strconv.Atoi(idStr)
        if err != nil || userID == 0 {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusBadRequest)
            json.NewEncoder(w).Encode(map[string]string{"error": "invalid user id"})
            return
        }

        email := r.Header.Get("X-User-Email")

        var username, location, birthday, work string
        
        err = db.QueryRow(`
            SELECT 
                COALESCE(username, ''),
                COALESCE(location, ''),
                COALESCE(birthday, ''),
                COALESCE(work, '')
            FROM users 
            WHERE id = ?`, userID).
            Scan(&username, &location, &birthday, &work)
        
        if err != nil {
            log.Printf("Dashboard query error: %v", err)
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusNotFound)
            json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
            return
        }

        // 3. Формируем JSON-ответ
        response := map[string]interface{}{
            "id":             userID,
            "email":          email,
            "username":       username,
            "location":       location,
            "birthday":       birthday,  
            "work":           work,
            "friendsCount":   0,
            "postsCount":     0,
            "followersCount": 0,
        }

        w.Header().Set("Content-Type", "application/json")
        if encErr := json.NewEncoder(w).Encode(response); encErr != nil {
            log.Printf(" JSON encode error: %v", encErr)
        }
    }
}