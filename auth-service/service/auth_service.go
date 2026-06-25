package service

import (
	"auth-service/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryInterface interface {
	Create(user *models.User) (int64, error)
	FindByIdentifier(identifier string, user *models.User) error
}
type AuthService struct {
	UserRepo UserRepositoryInterface
}

func NewAuthService(repo UserRepositoryInterface) *AuthService {
	return &AuthService{UserRepo: repo}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.UserResponse, error) {
	if req.Email == "" || req.Password == "" || req.Username == "" {
		return nil, errors.New("Все поля обязательны")
	}
	if len(req.Password) < 4 {
		return nil, errors.New("Пароль меньше 4 символов невозможен")
	}
	hashPass, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     req.Username,
		Email:    req.Email,
		Password: string(hashPass),
	}

	id, err := s.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:        id,
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}, nil
}

func (s *AuthService) Login(req *models.UserLogin) (string, error) {

	var user models.User
	err := s.UserRepo.FindByIdentifier(req.Identifier, &user)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"username": user.Name,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("secret"))
}
