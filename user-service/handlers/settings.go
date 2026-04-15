package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
	"user-service/database"
)

func decodeJSON(w http.ResponseWriter, r *http.Request, dest interface{}) bool {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(dest); err != nil {
		log.Printf("decodeJSON error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Неверный формат JSON",
		})
		return false
	}
	return true
}

func GetUserID(r *http.Request) int {
	userIDstr := r.Header.Get("X-User-ID")
	if userIDstr == "" {
		return 0
	}
	var userID int
	_, err := fmt.Sscanf(userIDstr, "%d", &userID)
	if err != nil {
		return 0
	}
	return userID
}

type Useranme struct {
	NewName string `json:"newName"`
}

func ChangeUsername(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := GetUserID(r)
		if userID == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Пользователь не авторизован",
			})
			return
		}

		var req Useranme
		if !decodeJSON(w, r, &req) {
			return
		}

		req.NewName = strings.TrimSpace(req.NewName)
		if req.NewName == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "введите имя",
			})
			return
		}
		if utf8.RuneCountInString(req.NewName) < 2 || utf8.RuneCountInString(req.NewName) > 50 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Имя должно быть от 2 до 50 символов",
			})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := userRepo.UpdateUsername(ctx, userID, req.NewName)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Ошибка базы данных",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"message": "Имя обновлено",
			"name":    req.NewName,
		})
	}
}

type City struct {
	City string `json:"city"`
}

func UpdateGEO(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		var city City
		if !decodeJSON(w, r, &city) {
			return
		}
		if city.City == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "введите город"})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := userRepo.UpdateCity(ctx, userID, city.City)
		if err != nil {
			if strings.Contains(err.Error(), "не найден") {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "город сменен",
		})
	}
}

type Work struct {
	WorkLocation string `json:"work_location"`
}

func UpdateWork(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := GetUserID(r)
		var work Work

		if !decodeJSON(w, r, &work) {
			return
		}
		if work.WorkLocation == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "введите место работы!!!"})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := userRepo.UpdateWork(ctx, userID, work.WorkLocation)
		if err != nil {
			if strings.Contains(err.Error(), "не найден") {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "место работы сменено",
		})
	}
}

type Password struct {
	OldPass string `json:"OldPass"`
	NewPass string `json:"NewPass"`
}

func UpdatePassword(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		var pass Password
		if !decodeJSON(w, r, &pass) {
			return
		}
		if pass.OldPass == "" || pass.NewPass == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Введите и старый и новый пароль",
			})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := userRepo.UpdatePass(ctx, userID, pass.OldPass, pass.NewPass)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка базы данных"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Пароль успешно обновлен",
		})

	}
}
