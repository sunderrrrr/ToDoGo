package repository

import (
	"ToDoGo/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ToDoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *ToDoListPostgres {
	return &ToDoListPostgres{db: db}
}

func (r *ToDoListPostgres) Create(userId int, list models.ToDo) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %S (title, descrtiption) VALUES $1, $2 RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil { //Проверяем запрос, пытаясь просканить полученный id
		tx.Rollback()
		return 0, err //Если транзакция не проходит, все действия БД откатываются
	}

	createUsersListsQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", userListsTable)
	_, err = tx.Exec(createUsersListsQuery, userId, list.Id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit() //Применяем изменения к БД и возвращаем ID
}
