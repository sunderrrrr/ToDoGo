package service

import (
	"ToDoGo/models"
	"ToDoGo/pkg/repository"
)

// Authorization Todo Перенести методы в User
type Authorization interface { //Методы авторизации
	CreateUser(User models.User) (int, error)
	GenerateToken(Username, Password string) (string, error)
	ParseToken(Token string) (models.User, error)
}

type TodoList interface {
	Create(UserId int, list models.ToDo) (int, error)
	GetAllLists(UserId int) ([]models.ToDo, error)
	GetListById(UserId int, ListId int) (models.ToDo, error)
	DeleteList(UserId int, ListId int) error
	UpdateList(UserId int, ListId int, NewList models.ToDo) error
}

type TodoItem interface {
	CreateItem(UserId int, ListId int, ItemText models.TodoItem) (int, error)
	GetAllItemsOfList(UserId int, ListId int) ([]models.TodoItem, error)
	DeleteItem(UserId int, ItemId int) error
	GetItemById(UserId int, ItemId int) (models.TodoItem, error)
	UpdateItem(UserId int, ListId int, ItemId int, UpdatedItem models.TodoItem) error
	//DeleteItem(UserId int, ItemId int) error
}
type User interface {
	ResetPassword(resetModel models.UserReset) error
	ResetPasswordRequest(email models.ResetRequest) error
	GeneratePasswordResetToken(email, signingKey string) (string, error)
}
type Service struct {
	Authorization
	TodoList
	TodoItem
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
		User:          NewUserService(repos.User),
	}
}
