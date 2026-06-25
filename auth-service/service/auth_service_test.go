package service

import (
	"auth-service/models"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockUserRepo struct{}

func (m *MockUserRepo) Create(user *models.User) (int64, error) {
	return 1, nil
}
func (m *MockUserRepo) FindByIdentifier(identifier string, user *models.User) error {
	user.ID = 1
	user.Email = "test@mail.com"
	user.Name = "test"
	user.Password = "$2a$10$hashedpassword"

	return nil
}
func TestRegister_Success(t *testing.T) {
	mockRepo := &MockUserRepo{}

	service := NewAuthService(mockRepo)
	req := &models.RegisterRequest{
		Username: "test",
		Email:    "test@mail.com",
		Password: "1234",
	}

	user, err := service.Register(req)

	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, "test@mail.com", user.Email)
}

//	func (m *MockUserRepo) Create(user *models.User) (int64, error) {
//		return 1, errors.New("db")
//	}
func Test_DBerror(t *testing.T) {
	service := NewAuthService(&MockUserRepo{})
	req := &models.RegisterRequest{
		Username: "test",
		Email:    "test@mail.com",
		Password: "1234",
	}
	_, err := service.Register(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db")
}

func Test_Empty(t *testing.T) {
	service := NewAuthService(&MockUserRepo{})
	req := &models.RegisterRequest{
		Username: "",
		Email:    "",
		Password: "",
	}

	_, err := service.Register(req)
	require.Error(t, err)
}

func Test_Short_pass(t *testing.T) {
	service := NewAuthService(&MockUserRepo{})
	req := &models.RegisterRequest{
		Username: "test",
		Email:    "test@mail.com",
		Password: "34",
	}
	_, err := service.Register(req)
	require.NoError(t, err)
}
