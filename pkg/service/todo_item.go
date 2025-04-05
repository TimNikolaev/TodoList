package service

import (
	"todo-app"
	"todo-app/pkg/repository"
)

type ToDoItemService struct {
	repo     repository.ToDoItem
	listRepo repository.ToDoList
}

func NewToDoItemService(repo repository.ToDoItem, listRepo repository.ToDoList) *ToDoItemService {
	return &ToDoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *ToDoItemService) Create(userID, listID int, item todo.ToDoItem) (int, error) {
	_, err := s.listRepo.GetByID(userID, listID)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listID, item)
}

func (s *ToDoItemService) GetAll(userID, listID int) ([]todo.ToDoItem, error) {
	return s.repo.GetAll(userID, listID)
}

func (s *ToDoItemService) GetByID(userID, itemID int) (todo.ToDoItem, error) {
	return s.repo.GetByID(userID, itemID)
}

func (s *ToDoItemService) Delete(userID, itemID int) error {
	return s.repo.Delete(userID, itemID)
}

func (s *ToDoItemService) Update(userID, itemID int, input todo.UpdateItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userID, itemID, input)
}
