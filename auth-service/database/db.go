package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"os"

	_ "github.com/go-sql-driver/mysql"
)
var db *sql.DB
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDB() (*sql.DB, error) {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser, dbPass, dbHost, dbPort, dbName)

    var db *sql.DB
    var err error

    for i := 0; i < 30; i++ {
        db, err = sql.Open("mysql", dsn)
        if err != nil {
            log.Printf("Попытка %d/30: ошибка sql.Open: %v", i+1, err)
            time.Sleep(3 * time.Second)
            continue
        }

        if err = db.Ping(); err != nil {
            log.Printf(" Попытка %d/30: MySQL не отвечает: %v", i+1, err)
            time.Sleep(3 * time.Second)
            continue
        }

        log.Println(" Подключение к MySQL установлено")
        break
    }

    if err != nil {
        return nil, fmt.Errorf("не удалось подключиться к БД после 30 попыток: %w", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    return db, nil
}