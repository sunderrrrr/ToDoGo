package service

import (
	"ToDoGo/models"
	"ToDoGo/pkg/repository"
	"crypto/sha1"
	"fmt"
)

const salt = "gdfgdf789fsd798ghdfh9d8f79d8fs" //абфускатор пароля "соль"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}
func (s *AuthService) CreateUser(user models.User) (int, error) { //Сначала создаем хэш пароля и передаем его на уровень репозитория
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
