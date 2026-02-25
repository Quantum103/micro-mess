package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func NewDB() (*sql.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/micro?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия соединения: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("MySQL не отвечает: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Подключение к MySQL установлено")
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
