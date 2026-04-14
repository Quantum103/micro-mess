package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) UpdateUsername(ctx context.Context, userID int, NewName string) error {
	query := `UPDATE users  SET username = ? WHERE id = ?`
	res, err := r.db.Exec(query, NewName, userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления имени: ")
	}
	rowsAf, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка пол кол-ва строк: ")
	}
	if rowsAf == 0 {
		return fmt.Errorf("пользователь не найден: ")
	}
	return nil
}

func (r *UserRepository) UpdateCity(ctx context.Context, userID int, City string) error {
	query := `UPDATE users SET location = ? WHERE id = ?`
	res, err := r.db.Exec(query, City, userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления города: %w", err)
	}
	rowsAf, _ := res.RowsAffected()
	if rowsAf == 0 {
		return fmt.Errorf("пользователь не найден")
	}
	return nil
}

func (r *UserRepository) UpdateWork(ctx context.Context, userID int, location string) error {
	query := `UPDATE users SET work = ? WHERE id = ?`
	res, err := r.db.Exec(query, location, userID)
	if err != nil {
		return fmt.Errorf("ошибка обновления места работы: %w", err)
	}
	rowsAf, _ := res.RowsAffected()
	if rowsAf == 0 {
		return fmt.Errorf("пользователь не найден")
	}
	return nil
}

func (r *UserRepository) UpdatePass(ctx context.Context, userId int, OldPass, NewPass string) error {
	var OldPassOnDatabase string
	QueryPassword := `SELECT password FROM users WHERE id = ?`
	err := r.db.QueryRow(QueryPassword, userId).Scan(&OldPassOnDatabase)
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
	res, err := r.db.Exec(queryUpdate, string(newHash), userId)
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
