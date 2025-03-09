package repository

import (
	"ToDoGo/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type TodoList interface {
	Create(UserId int, list models.ToDo) (int, error)
	GetAllLists(UserId int) ([]models.ToDo, error)
	GetListById(UserId int, ListId int) (models.ToDo, error)
	DeleteList(UserId int, ListId int) error
	UpdateList(UserId int, ListId int, NewList models.ToDo) error
}

type TodoItem interface {
	CreateItem(UserId int, ListId int, Item models.TodoItem) (int, error)
	GetAllItemsOfList(UserId int, ListId int) ([]models.TodoItem, error)
	DeleteItem(UserId int, ItemId int) error
	GetItemById(UserId int, ItemId int) (models.TodoItem, error)
	UpdateItem(UserId int, ListId int, ItemId int, UpdatedItem models.TodoItem) error
}
type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
