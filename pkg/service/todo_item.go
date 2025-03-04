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

/*func (s *TodoItemService) GetAllItems(UserId int) ([]models.TodoItem, error) {
	return s.repo.GetAllItems(UserId)
}

func (s *TodoItemService) GetItemById(UserId int, ItemId int) (models.TodoItem, error) {
	return s.repo.GetItemById(UserId, ItemId)
}

func (s *TodoItemService) DeleteItem(UserId int, ItemId int) error {
	return s.repo.DeleteItem(UserId, ItemId)
}
*/
