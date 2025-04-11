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
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
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
		user.Name,
	})
	fmt.Println("EndGenToken")
	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return models.User{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return models.User{}, errors.New("token claims are not of type *tokenClaims")
	}

	returnUser := models.User{
		Id:   claims.UserId,
		Name: claims.Username,
	}

	return returnUser, nil
}

func (s *AuthService) GeneratePasswordResetToken(email, signingKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(1 * time.Hour).Unix(), // Токен действует 1 час
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ResetPassword(resetModel models.UserReset, resetToken string) error {
	token, err := jwt.ParseWithClaims(resetToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return errors.New("token claims are not of type jwt.MapClaims or token invalid")
	}
	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return errors.New("токен истёк")
	}

	// 4. Получаем email из токена
	email, ok := claims["email"].(string)
	if !ok {
		return errors.New("email не найден в токене")
	}

	return s.repo.ResetPassword(email, generatePasswordHash(resetModel.OldPass), generatePasswordHash(resetModel.NewPass))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
