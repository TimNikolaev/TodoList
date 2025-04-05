package service

import (
	"todo-app"
	"todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(userName, password string) (string, error)
	ParseToken(token string) (int, error)
}

type ToDoList interface {
	Create(userID int, list todo.ToDoList) (int, error)
	GetAll(userID int) ([]todo.ToDoList, error)
	GetByID(userID, listID int) (todo.ToDoList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input todo.UpdateListInput) error
}

type ToDoItem interface {
	Create(userID, listID int, item todo.ToDoItem) (int, error)
	GetAll(userID, listID int) ([]todo.ToDoItem, error)
	GetByID(userID, itemID int) (todo.ToDoItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input todo.UpdateItemInput) error
}

type Service struct {
	Authorization
	ToDoList
	ToDoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		ToDoList:      NewToDoListService(repos.ToDoList),
		ToDoItem:      NewToDoItemService(repos.ToDoItem, repos.ToDoList),
	}
}
