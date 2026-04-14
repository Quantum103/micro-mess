package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func NewDB() (*sql.DB, error) {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser, dbPass, dbHost, dbPort, dbName)

    var err error

    for i := 0; i < 30; i++ {
        db, err = sql.Open("mysql", dsn)
        if err != nil {
            log.Printf(" Попытка %d/30: ошибка sql.Open: %v", i+1, err)
            time.Sleep(3 * time.Second)
            continue
        }

        if err = db.Ping(); err != nil {
            log.Printf(" Попытка %d/30: MySQL не отвечает: %v", i+1, err)
            time.Sleep(3 * time.Second)
            continue
        }

        log.Println("Подключение к MySQL установлено")
        break
    }

    if err != nil {
        return nil, fmt.Errorf("не удалось подключиться к БД после 30 попыток: %w", err)
    }

    // Настройки пула соединений
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    return db, nil
}

func GetDB() *sql.DB {
	return db
}
func UpdateUsername(userID int, NewName string) error {
	query := `UPDATE users  SET username = ? WHERE id = ?`
	res, err := db.Exec(query, NewName, userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления имени: %w", err)
	}
	rowsAf, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка пол кол-ва строк: %w", err)
	}
	if rowsAf == 0 {
		return fmt.Errorf("пользователь не найден: %w", err)
	}
	return nil
}

func UpdateCity(userID int, City string) error {
	query := `UPDATE users SET location = ? WHERE id = ?`
	res, err := db.Exec(query, City, userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления города: %w", err)
	}
	rowsAf, _ := res.RowsAffected()
	if rowsAf == 0 {
		return fmt.Errorf("пользователь не найден")
	}
	return nil
}

func UpdateWork(userID int, location string) error {
	query := `UPDATE users SET work = ? WHERE id = ?`
	res, err := db.Exec(query, location, userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления места работы: %w", err)
	}
	rowsAf, _ := res.RowsAffected()
	if rowsAf == 0 {
		return fmt.Errorf("пользователь не найден")
	}
	return nil
}

func UpdatePass(userId int, OldPass, NewPass string) error {
	var OldPassOnDatabase string
	QueryPassword := `SELECT password FROM users WHERE id = ?`
	err := db.QueryRow(QueryPassword, userId).Scan(&OldPassOnDatabase)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("пользователь не найден: %w", err)
		}
		return fmt.Errorf("ошибка получения пароля из БД: %w", err)

	}
	err = bcrypt.CompareHashAndPassword([]byte(OldPassOnDatabase), []byte(OldPass))
	if err != nil {
		return fmt.Errorf("неверный старый пароль")
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(NewPass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хеширования нового пароля: %w", err)
	}
	queryUpdate := `UPDATE users SET password = ? WHERE id = ?`
	res, err := db.Exec(queryUpdate, string(newHash), userId)
	if err != nil {
		return fmt.Errorf("ошибка обновления пароля в БД: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки результата обновления: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("пользователь не найден или пароль не изменился")
	}

	return nil
}
