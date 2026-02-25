package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"user-service/database"
)

func decodeJSON(w http.ResponseWriter, r *http.Request, dest interface{}) bool {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(dest); err != nil {
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

func ChangeUsername(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		var req Useranme
		if !decodeJSON(w, r, &req) {
			return
		}
		if req.NewName == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Имя не может быть пустым или его длина слишком большая"})
			return
		}

		err := database.UpdateUsername(userID, req.NewName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ошибка базы данных"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Имя успешно обновлено",
		})
	}
}

type City struct {
	City string `json:"city"`
}

func UpdateGEO(db *sql.DB) http.HandlerFunc {
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
		err := database.UpdateCity(userID, city.City)
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
	WorkLocaion string `json:"work_location"`
}

func UpdateWork(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		var work Work
		if !decodeJSON(w, r, &work) {
			return
		}
		if work.WorkLocaion == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "введите место работы!!!"})
			return
		}

		err := database.UpdateWork(userID, work.WorkLocaion)
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
