package repository

import (
	"ToDoGo/models"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ToDoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *ToDoItemPostgres {
	return &ToDoItemPostgres{db: db}
}

func (r *ToDoItemPostgres) CreateItem(UserId int, ListId int, Item models.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	row := tx.QueryRow(createItemQuery, Item.Title, Item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) values ($1, $2)", listItemsTable)
	_, err = tx.Exec(createListItemsQuery, ListId, itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	return itemId, tx.Commit()
}

func (r *ToDoItemPostgres) GetAllItemsOfList(UserId int, ListId int) ([]models.TodoItem, error) {
	var rowsOutput []models.TodoItem
	query := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti JOIN %s li ON ti.id = li.item_id JOIN %s ul ON li.list_id = ul.list_id WHERE ul.list_id = $1  AND ul.user_id = $2", todoItemsTable, listItemsTable, userListsTable)
	rows := r.db.Select(&rowsOutput, query, ListId, UserId)
	if rows != nil {
		return nil, rows
	}
	fmt.Printf("todo_item_postgres.go: %s", rows)

	return rowsOutput, nil
}

func (r *ToDoItemPostgres) DeleteItem(UserId int, ItemId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1 AND EXISTS (SELECT 1 FROM %s li JOIN %s ul ON li.list_id = ul.list_id WHERE li.item_id = $1 AND ul.user_id = $2);`, todoItemsTable, listItemsTable, userListsTable)

	if _, err := tx.Exec(query, ItemId, UserId); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *ToDoItemPostgres) GetItemById(UserId int, ItemId int) (models.TodoItem, error) {
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
									INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`, todoItemsTable, listItemsTable, userListsTable)
	var item models.TodoItem
	err := r.db.Get(&item, query, ItemId, UserId)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (r *ToDoItemPostgres) UpdateItem(UserId int, ListId int, ItemId int, UpdatedItem models.TodoItem) error {
	var OldItem models.TodoItem
	var ResItem models.TodoItem

	if UpdatedItem.Title == "" && UpdatedItem.Description == "" {
		return errors.New("missing fields")
	}

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li on li.item_id = ti.id
									INNER JOIN %s ul on ul.list_id = li.list_id WHERE ti.id = $1 AND ul.user_id = $2`, todoItemsTable, listItemsTable, userListsTable)

	err := r.db.Get(&OldItem, query, ItemId, UserId)

	if err != nil {
		return err
	}
	fmt.Println("first req sent")
	if UpdatedItem.Title != OldItem.Title && UpdatedItem.Title != "" {
		ResItem.Title = UpdatedItem.Title
	} else {
		ResItem.Title = OldItem.Title
	}
	if UpdatedItem.Description != OldItem.Description && UpdatedItem.Description != "" {
		ResItem.Description = UpdatedItem.Description
	} else {
		ResItem.Description = OldItem.Description
	}
	if UpdatedItem.Done != OldItem.Done {
		ResItem.Done = UpdatedItem.Done
	} else {
		ResItem.Done = OldItem.Done
	}
	fmt.Printf("second send start")
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	updateQuery := fmt.Sprintf(`UPDATE %s
			SET title = $1, description = $2, done = $3
			WHERE id = $4
			  AND EXISTS (
				  SELECT 1
				  FROM %s li
				  JOIN %s ul ON li.list_id = ul.list_id
				  WHERE li.item_id = $4
					AND ul.user_id = $5
			  );`, todoItemsTable, listItemsTable, userListsTable)
	_, err = tx.Exec(updateQuery, ResItem.Title, ResItem.Description, ResItem.Done, ItemId, UserId)

	return tx.Commit()
}
