package database

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestUpdateUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating mock: %s", err)
	}
	defer db.Close()

	repo := &UserRepository{db: db}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET username = \\? WHERE id = \\?").
			WithArgs("newname", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateUsername(context.Background(), 1, "newname")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET username = \\? WHERE id = \\?").
			WithArgs("newname", 1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateUsername(context.Background(), 1, "newname")
		if err == nil || err.Error() != "пользователь не найден: " {
			t.Errorf("expected 'пользователь не найден: ', got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("exec error", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET username = \\? WHERE id = \\?").
			WithArgs("newname", 1).
			WillReturnError(errors.New("db error"))

		err := repo.UpdateUsername(context.Background(), 1, "newname")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestUpdateCity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating mock: %s", err)
	}
	defer db.Close()

	repo := &UserRepository{db: db}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET location = \\? WHERE id = \\?").
			WithArgs("Moscow", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateCity(context.Background(), 1, "Moscow")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET location = \\? WHERE id = \\?").
			WithArgs("Moscow", 1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateCity(context.Background(), 1, "Moscow")
		if err == nil || err.Error() != "пользователь не найден" {
			t.Errorf("expected 'пользователь не найден', got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestUpdateWork(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating mock: %s", err)
	}
	defer db.Close()

	repo := &UserRepository{db: db}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET work = \\? WHERE id = \\?").
			WithArgs("Developer", 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdateWork(context.Background(), 1, "Developer")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE users SET work = \\? WHERE id = \\?").
			WithArgs("Developer", 1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdateWork(context.Background(), 1, "Developer")
		if err == nil || err.Error() != "пользователь не найден" {
			t.Errorf("expected 'пользователь не найден', got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}

func TestUpdatePass(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error creating mock: %s", err)
	}
	defer db.Close()

	repo := &UserRepository{db: db}

	oldPass := "oldpassword"
	newPass := "newpassword"
	hashedOldPass, _ := bcrypt.GenerateFromPassword([]byte(oldPass), bcrypt.DefaultCost)

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery("SELECT password FROM users WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(string(hashedOldPass)))

		mock.ExpectExec("UPDATE users SET password = \\? WHERE id = \\?").
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.UpdatePass(context.Background(), 1, oldPass, newPass)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("user not found on select", func(t *testing.T) {
		mock.ExpectQuery("SELECT password FROM users WHERE id = \\?").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		err := repo.UpdatePass(context.Background(), 1, oldPass, newPass)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("wrong old password", func(t *testing.T) {
		mock.ExpectQuery("SELECT password FROM users WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(string(hashedOldPass)))

		err := repo.UpdatePass(context.Background(), 1, "wrongpassword", newPass)
		if err == nil || err.Error() != "неверный старый пароль" {
			t.Errorf("expected 'неверный старый пароль', got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})

	t.Run("update not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT password FROM users WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow(string(hashedOldPass)))

		mock.ExpectExec("UPDATE users SET password = \\? WHERE id = \\?").
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.UpdatePass(context.Background(), 1, oldPass, newPass)
		if err == nil || err.Error() != "пользователь не найден или пароль не изменился" {
			t.Errorf("expected 'пользователь не найден или пароль не изменился', got %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %s", err)
		}
	})
}
