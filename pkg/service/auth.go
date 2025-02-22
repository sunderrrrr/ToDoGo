package service

import (
	"ToDoGo/models"
	"ToDoGo/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	salt       = "gdfgdf789fsd798ghdfh9d8f79d8fs"                            //абфускатор пароля "соль"
	signingKey = "js786b87^*bn98v79&(*jhkjhKj6kiu6iU^^u6iU^uk6tiuufv6biu^u6" //ключ подписи
	tokenTTL   = time.Hour * 12                                              //время действия токена
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}
func (s *AuthService) CreateUser(user models.User) (int, error) { //Сначала создаем хэш пароля и передаем его на уровень репозитория
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}
func (s *AuthService) GenerateToken(username, password string) (string, error) {
	//get user from db
	fmt.Println("StartGenToken")
	user, err := s.repo.GetUser(username, generatePasswordHash(password))

	if err != nil {

		fmt.Println("Failed Get User")
		return "fail", err
	}
	fmt.Println("auth.go: ", user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})
	fmt.Println("EndGenToken")
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
