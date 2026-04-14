package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"user-service/database"
)

func ListUserHandler(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		if userID == 0 {
			http.Error(w, `{"error":"не знаем такого"}`, http.StatusUnauthorized)
			return
		}
		limit, offset := 20, 0
		if l := r.URL.Query().Get("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 {
				limit = v
			}
		}
		if o := r.URL.Query().Get("offset"); o != "" {
			if v, err := strconv.Atoi(o); err == nil && v >= 0 {
				offset = v
			}
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		users, err := userRepo.ListAllUsers(ctx, userID, limit, offset)
		if err != nil {
			http.Error(w, `{"ошибка":"ошибка базы данных"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users":  users,
			"total":  len(users),
			"limit":  limit,
			"offset": offset,
		})

	}
}

func GetFrineds(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		UserID := GetUserID(r)
		if UserID == 0 {
			http.Error(w, `{"error":"не знаем такого"}`, http.StatusUnauthorized)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		friends, err := userRepo.GetFrineds(ctx, UserID, 50, 0)
		if err != nil {
			http.Error(w, `{"ошибка":"ошибка базы данных"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"friends": friends,
			"count":   len(friends),
		})
	}
}

func AddFriendHandler(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		if userID == 0 {
			http.Error(w, `{"error":"не знаем такого"}`, http.StatusUnauthorized)
			return
		}
		var req struct {
			FriendID int `json:"friend_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"ошибка JSON"}`, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		err := userRepo.AddFriend(ctx, userID, req.FriendID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Заявка отправлена"})
	}
}

func AcceptFriendHandler(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)

		var req struct {
			FriendID int `json:"friend_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"bad json"}`, http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		err := userRepo.AcceptFriend(ctx, userID, req.FriendID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "accepted",
		})
	}
}
