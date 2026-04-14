package handlers

import (
<<<<<<< HEAD
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"
	"user-service/database"
=======
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"user-service/database"
	"unicode/utf8"
	"log"
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
)

func decodeJSON(w http.ResponseWriter, r *http.Request, dest interface{}) bool {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(dest); err != nil {
<<<<<<< HEAD
		log.Printf("decodeJSON error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Неверный формат JSON",
		})
		return false
=======
		log.Printf("decodeJSON error: %v", err) 
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
		"error": "Неверный формат JSON",
		})
		return false	
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
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

<<<<<<< HEAD
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
=======
func ChangeUsername(w http.ResponseWriter, r *http.Request) {
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
    if req.NewName == ""{
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

    err := database.UpdateUsername(userID, req.NewName)
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
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
}

type City struct {
	City string `json:"city"`
}

<<<<<<< HEAD
func UpdateGEO(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
=======
func UpdateGEO(w http.ResponseWriter, r *http.Request) {
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
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
<<<<<<< HEAD
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := userRepo.UpdateCity(ctx, userID, city.City)
=======
		err := database.UpdateCity(userID, city.City)
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
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
<<<<<<< HEAD
}
=======

>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56

type Work struct {
	WorkLocation string `json:"work_location"`
}

<<<<<<< HEAD
func UpdateWork(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := GetUserID(r)
		var work Work

=======
func UpdateWork(w http.ResponseWriter, r *http.Request) {
	    log.Printf("UpdateWork: ЗАПРОС ДОШЁЛ! Method=%s, Path=%s", r.Method, r.URL.Path)
	
	userID := GetUserID(r)
		var work Work

		 log.Printf(" UpdateWork: userID=%d, raw header='%s'", 
        userID, r.Header.Get("X-User-ID"))

>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
		if !decodeJSON(w, r, &work) {
			return
		}
		if work.WorkLocation == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "введите место работы!!!"})
			return
		}
<<<<<<< HEAD
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := userRepo.UpdateWork(ctx, userID, work.WorkLocation)
=======

		err := database.UpdateWork(userID, work.WorkLocation)
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
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
<<<<<<< HEAD
		})
	}
}

=======
	})
}


>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
type Password struct {
	OldPass string `json:"OldPass"`
	NewPass string `json:"NewPass"`
}

<<<<<<< HEAD
func UpdatePassword(userRepo *database.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
=======
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
		userID := GetUserID(r)
		var pass Password
		if !decodeJSON(w, r, &pass) {
			return
		}
		if pass.OldPass == "" || pass.NewPass == "" {
			w.Header().Set("Content-Type", "application/json")
<<<<<<< HEAD
			w.WriteHeader(http.StatusBadRequest)
=======
			w.WriteHeader(http.StatusBadRequest)  
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Введите и старый и новый пароль",
			})
			return
		}
<<<<<<< HEAD
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		err := userRepo.UpdatePass(ctx, userID, pass.OldPass, pass.NewPass)
=======
		err := database.UpdatePass(userID, pass.OldPass, pass.NewPass)
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
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
<<<<<<< HEAD

	}
=======
	
>>>>>>> 504fd3e5a511bf68e5f35cecc92e257c0bb17d56
}
