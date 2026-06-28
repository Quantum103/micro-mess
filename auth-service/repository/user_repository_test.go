package repository

import (
	"auth-service/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewUserRepository(db)

	mock.ExpectExec("INSERT INTO users").
		WithArgs("anton", "anton@mail.com", "123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	user := &models.User{
		Name:     "anton",
		Email:    "anton@mail.com",
		Password: "123",
	}

	id, err := repo.Create(user)

	require.NoError(t, err)
	require.Equal(t, int64(1), id)
}
