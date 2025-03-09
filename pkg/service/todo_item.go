package service

import (
	"ToDoGo/models"
	"ToDoGo/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) CreateItem(UserId int, ListId int, Item models.TodoItem) (int, error) {
	_, err := s.listRepo.GetListById(UserId, ListId)
	if err != nil { //список не существует или не принадлежит пользователю
		return 0, err
	}
	return s.repo.CreateItem(UserId, ListId, Item)
}

func (s *TodoItemService) GetAllItemsOfList(UserId int, ListId int) ([]models.TodoItem, error) {
	return s.repo.GetAllItemsOfList(UserId, ListId)
}

func (s *TodoItemService) GetItemById(UserId int, ItemId int) (models.TodoItem, error) {
	return s.repo.GetItemById(UserId, ItemId)
}

func (s *TodoItemService) DeleteItem(UserId int, ItemId int) error {
	return s.repo.DeleteItem(UserId, ItemId)
}

func (s *TodoItemService) UpdateItem(UserId int, ListId int, ItemId int, UpdatedItem models.TodoItem) error {
	return s.repo.UpdateItem(UserId, ListId, ItemId, UpdatedItem)
}
