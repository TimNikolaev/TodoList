package service

import (
	"todo-app"
	"todo-app/pkg/repository"
)

type ToDoListService struct {
	repo repository.ToDoList
}

func NewToDoListService(repo repository.ToDoList) *ToDoListService {
	return &ToDoListService{repo: repo}
}

func (s *ToDoListService) Create(userID int, list todo.ToDoList) (int, error) {
	return s.repo.Create(userID, list)
}

func (s *ToDoListService) GetAll(userID int) ([]todo.ToDoList, error) {
	return s.repo.GetAll(userID)
}

func (s *ToDoListService) GetByID(userID, listID int) (todo.ToDoList, error) {
	return s.repo.GetByID(userID, listID)
}

func (s *ToDoListService) Delete(userID, listID int) error {
	return s.repo.Delete(userID, listID)
}

func (s *ToDoListService) Update(userID, listID int, input todo.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userID, listID, input)
}
