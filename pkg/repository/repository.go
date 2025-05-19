package repository

import (
	"todo-app"
	"todo-app/pkg/repository/postgres"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(userName, password string) (todo.User, error)
}

type ToDoList interface {
	Create(userID int, list todo.ToDoList) (int, error)
	GetAll(userID int) ([]todo.ToDoList, error)
	GetByID(userID, listID int) (todo.ToDoList, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input todo.UpdateListInput) error
}

type ToDoItem interface {
	Create(listID int, item todo.ToDoItem) (int, error)
	GetAll(userID, listID int) ([]todo.ToDoItem, error)
	GetByID(userID, itemID int) (todo.ToDoItem, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input todo.UpdateItemInput) error
}

type Repository struct {
	Authorization
	ToDoList
	ToDoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(db),
		ToDoList:      postgres.NewToDoListPostgres(db),
		ToDoItem:      postgres.NewToDoItemPostgres(db),
	}
}
