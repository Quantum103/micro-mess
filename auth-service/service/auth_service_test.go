package service

import (
	"auth-service/models"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct{}

func (m *MockUserRepo) Create(user *models.User) (int64, error) {
	return 1, nil
}
func (m *MockUserRepo) FindByIdentifier(identifier string, user *models.User) error {
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("12345678"),
		bcrypt.DefaultCost,
	)
	user.ID = 1
	user.Email = "test@mail.com"
	user.Name = "test"
	user.Password = string(hash)

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

/*
	Test LOGIN
*/

func TestLogin_Succses(t *testing.T) {
	service := NewAuthService(&MockUserRepo{})
	req := &models.UserLogin{
		Identifier: "anton@gmail.com",
		Password:   "12345678",
	}
	_, err := service.Login(req)
	require.NoError(t, err)
}

func TestLogin_DBerror(t *testing.T) {
	service := NewAuthService(&MockUserRepo{})
	req := &models.UserLogin{
		Identifier: "anton@gmail.com",
		Password:   "12345678",
	}
	_, err := service.Login(req)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db")
}

func TestLogin_Empty(t *testing.T) {
	service := NewAuthService(&MockUserRepo{})
	req := &models.UserLogin{
		Identifier: "",
		Password:   "",
	}
	_, err := service.Login(req)
	require.NoError(t, err)
}

func TestLogin_Short_Pass(t *testing.T) {
	service := NewAuthService(&MockUserRepo{})
	req := &models.UserLogin{
		Identifier: "anton@gmail.com",
		Password:   "12",
	}
	_, err := service.Login(req)
	require.NoError(t, err)
}
