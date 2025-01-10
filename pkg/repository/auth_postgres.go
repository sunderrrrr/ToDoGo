package repository

import (
	"ToDoGo/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (ap *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", userTable)
	row := ap.db.QueryRow(query, user.Name, user.Username, user.Password) // Направление запроса в бд
	if err := row.Scan(&id); err != nil {                                 //Записываем в id результат row
		return 0, err
	}

	return id, nil
}

func (ap *AuthPostgres) GetUser(username string, password string) (models.User, error) {
	var user models.User // то, что вернет функция
	query := fmt.Sprintf("SELECT FROM %s WHERE username=$1 and password_hash=$2", userTable)
	err := ap.db.Get(&user, query, username, password) // Записываем значение используя указатель на структуру user

	return user, err
}
