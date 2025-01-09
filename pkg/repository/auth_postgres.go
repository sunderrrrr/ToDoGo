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
	return &AuthPostgres{}
}

func (ap *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNINIG id", userTable)
	row := ap.db.QueryRow(query, user.Name, user.Username, user.Password) // Направление запроса в бд
	if err := row.Scan(&id); err != nil {                                 //Записываем в id результат row
		return 0, err
	}

	return id, nil
}
