package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"petproject2/cmd/internal/repository"
	"petproject2/cmd/models"
	"time"
)

var jwtSecret = "secret"

type UserService interface {
	RegisterUser(ctx context.Context, user *models.Users) (int, error)
	LoginUser(ctx context.Context, username, password string) (string, error)
}

type ServiceUser struct {
	repo *repository.RepositoryUser
}

func NewServiceUser(repo *repository.RepositoryUser) *ServiceUser {
	return &ServiceUser{repo: repo}
}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err

}

func (s *ServiceUser) RegisterUser(ctx context.Context, user *models.Users) (int, error) {
	hashedPassword, err := GeneratehashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("Ошибка с Генерацией пароля")
	}

	user.Password = hashedPassword

	return s.repo.CreateUser(ctx, user)
}

func (s *ServiceUser) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByName(ctx, username)
	if err != nil {
		return "", fmt.Errorf("Ошибка с Получением Имени: %w", err)
	}

	fmt.Println("Stored password hash:", user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Ошибка при сравнении паролей:", err)
		return "", fmt.Errorf("Ошибка, неверный пароль")
	}

	token, err := generateJwt(user.ID, user.Name)
	if err != nil {
		return "", err
	}

	return token, nil
}

func generateJwt(userID int, name string) (string, error) {
	exptime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"exp":     exptime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))

}
