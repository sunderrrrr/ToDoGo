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
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err //Если транзакция не проходит, все действия БД откатываются
	}

	createUsersListsQuery := fmt.Sprintf(`INSERT INTO %s (user_id, list_id) VALUES ($1, $2)`, userListsTable)
	_, err = tx.Exec(createUsersListsQuery, userId, id)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
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
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func (r *ToDoListPostgres) UpdateList(UserId int, ListId int, NewList models.ToDo) error {
	var OldList models.ToDo
	var ResList models.ToDo
	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`, todoListsTable, userListsTable)
	err := r.db.Get(&OldList, query, UserId, ListId)
	if err != nil {
		return err
	}

	if OldList.Title != NewList.Title && NewList.Title != "" {
		ResList.Title = NewList.Title
	} else {
		ResList.Title = OldList.Title
	}

	if OldList.Description != NewList.Description && NewList.Description != "" {
		ResList.Description = NewList.Description
	} else {
		ResList.Description = OldList.Description
	}
	ResList.Id = OldList.Id
	tx, err := r.db.Begin()
	fmt.Printf("todo_list_postgres.go:\n Old list: %v\n  New list %v\n  Result list: %v\n ", OldList, NewList, ResList)
	//query = fmt.Sprintf("UPDATE %s SET title = $1, description = $2 WHERE id = $3 AND EXISTS (SELECT 1 FROM lists_items li JOIN users_lists ul ON li.list_id = ul.list_id  WHERE ul.user_id = $4)", todoListsTable)
	query = "UPDATE todo_lists SET title = $1, description = $2 WHERE id = $3 AND EXISTS (SELECT 1 FROM users_lists WHERE users_lists.list_id = todo_lists.id AND users_lists.user_id = $4);"
	_, err = tx.Exec(query, ResList.Title, ResList.Description, ResList.Id, UserId)
	if err != nil {
		tx.Rollback()
		return err
	} else {
		err := tx.Commit()
		if err != nil {
			return err
		}
	}

	return err
}
