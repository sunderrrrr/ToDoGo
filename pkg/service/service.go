package service

import (
	"ToDoGo/models"
	"ToDoGo/pkg/repository"
)

type Authorization interface { //Методы авторизации
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(UserId int, list models.ToDo) (int, error)
	GetAllLists(UserId int) ([]models.ToDo, error)
	GetListById(UserId int, ListId int) (models.ToDo, error)
	DeleteList(UserId int, ListId int) error
	UpdateList(UserId int, ListId int) error
}

type TodoItem interface {
}
type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
	}
}
