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
	createListQuery := fmt.Sprintf(`INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id`, todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil { //Проверяем запрос, пытаясь просканить полученный id
		tx.Rollback()
		return 0, err //Если транзакция не проходит, все действия БД откатываются
	}

	createUsersListsQuery := fmt.Sprintf(`INSERT INTO %s (user_id, list_id) VALUES ($1, $2)`, userListsTable)
	_, err = tx.Exec(createUsersListsQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit() //Применяем изменения к БД и возвращаем ID
}
func (r *ToDoListPostgres) GetAllLists(UserId int) ([]models.ToDo, error) {
	var lists []models.ToDo
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1", todoListsTable, userListsTable)
	err := r.db.Select(&lists, query, UserId)
	if err != nil {
		return nil, err
	}

	return lists, err
}

func (r *ToDoListPostgres) GetListById(UserId int, ListId int) (models.ToDo, error) {
	var list models.ToDo
	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`, todoListsTable, userListsTable)

	err := r.db.Get(&list, query, UserId, ListId)
	return list, err
}

func (r *ToDoListPostgres) DeleteList(UserId int, ListId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND EXISTS (  SELECT 1 FROM %s WHERE user_id = $2 AND list_id = $1);", todoListsTable, userListsTable)
	_, err = tx.Exec(query, ListId, UserId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
