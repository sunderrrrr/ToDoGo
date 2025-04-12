package service

import (
	"ToDoGo/models"
	"ToDoGo/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"net/smtp"
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

func (s *AuthService) ResetPassword(resetModel models.UserReset) error {
	token, err := jwt.ParseWithClaims(resetModel.Token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
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

func (s *AuthService) ResetPasswordRequest(email models.ResetRequest) error {
	token, err := s.GeneratePasswordResetToken(email.Login, signingKey)
	if err != nil {
		return err
	}
	resetLink := fmt.Sprintf("%s/reset-confirm/?token=%s", viper.GetString("frontendUrl"), token)
	from := viper.GetString("smtp.senderMail")
	password := viper.GetString("smtp.senderPass")

	// Информация о получателе
	to := []string{
		email.Login,
	}

	// smtp сервер конфигурация
	smtpHost := viper.GetString("smtp.host")
	smtpPort := viper.GetString("smtp.port")

	// Сообщение.
	message := []byte("<h1>Сброс пароля</h1>\n" +
		"<p>Перейдите по ссылке, чтобы сбросить пароль</p>\n" +
		"<a href=\"" + resetLink + "\">Сброс</a>\n" +
		"<p>Если вы не запрашивали сброс, не переходите. Время действия ссылки один час</p>")

	// Авторизация.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Отправка почты.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err

	}
	fmt.Println("Почта отправлена!")
	return nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
