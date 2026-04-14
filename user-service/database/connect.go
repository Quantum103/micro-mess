package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func NewDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	var newDB *sql.DB
	for i := 0; i < 30; i++ {
		newDB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf(" Попытка %d/30: ошибка sql.Open: %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}

		if err = newDB.Ping(); err != nil {
			newDB.Close()
			log.Printf(" Попытка %d/30: MySQL не отвечает: %v", i+1, err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД после 30 попыток: %w", err)
	}

	// Настройки пула соединений
	newDB.SetMaxOpenConns(25)
	newDB.SetMaxIdleConns(5)
	newDB.SetConnMaxLifetime(5 * time.Minute)

	return newDB, nil
}
