package todo

import "errors"

type ToDoList struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binbing:"required"`
	Discription string `json:"discription" db:"discription"`
}

type ToDoItem struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Discription string `json:"discription" db:"discription"`
	Done        bool   `json:"done" db:"done"`
}

type UpdateListInput struct {
	Title       *string `json:"title"`
	Discription *string `json:"discription"`
}

func (i *UpdateListInput) Validate() error {
	if i.Title == nil && i.Discription == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

type UpdateItemInput struct {
	Title       *string `json:"title"`
	Discription *string `json:"discription"`
	Done        *bool   `json:"done"`
}

func (i *UpdateItemInput) Validate() error {
	if i.Title == nil && i.Discription == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
