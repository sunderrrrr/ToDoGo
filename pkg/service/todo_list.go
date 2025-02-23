package service

import (
	"ToDoGo/models"
	"ToDoGo/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list models.ToDo) (int, error) {
	return s.repo.Create(userId, list)
}
func (s *TodoListService) GetAllLists(UserId int) ([]models.ToDo, error) {
	return s.repo.GetAllLists(UserId)
}

func (s *TodoListService) GetListById(UserId int, ListId int) (models.ToDo, error) {
	return s.repo.GetListById(UserId, ListId)
}

func (s *TodoListService) DeleteList(UserId int, ListId int) error {
	return s.repo.DeleteList(UserId, ListId)
}
func (s *TodoListService) UpdateList(UserId int, ListId int, NewList models.ToDo) error {
	return s.repo.UpdateList(UserId, ListId, NewList)
}
