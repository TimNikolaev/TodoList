package postgres

import (
	"fmt"
	"strings"
	"todo-app"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ToDoListPostgres struct {
	db *sqlx.DB
}

func NewToDoListPostgres(db *sqlx.DB) *ToDoListPostgres {
	return &ToDoListPostgres{db: db}
}

func (r *ToDoListPostgres) Create(userID int, list todo.ToDoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var listID int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, discription) VALUES ($1, $2) RETURNING id", todoListsTable)
	if err := tx.QueryRow(createListQuery, list.Title, list.Discription).Scan(&listID); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListsQuery := fmt.Sprintf("INSERT INTO %s (user_id, lists_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListsQuery, userID, listID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return listID, tx.Commit()

}

func (r ToDoListPostgres) GetAll(userID int) ([]todo.ToDoList, error) {
	var lists []todo.ToDoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.discription FROM %s tl INNER JOIN %s ul on tl.id = ul.lists_id WHERE ul.user_id = $1", todoListsTable, usersListsTable)
	err := r.db.Select(&lists, query, userID)

	return lists, err
}

func (r ToDoListPostgres) GetByID(userID, listID int) (todo.ToDoList, error) {
	var list todo.ToDoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.discription FROM %s tl INNER JOIN %s ul on tl.id = ul.lists_id WHERE ul.user_id = $1 AND ul.lists_id = $2", todoListsTable, usersListsTable)
	err := r.db.Get(&list, query, userID, listID)

	return list, err
}

func (r ToDoListPostgres) Delete(userID, listID int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.lists_id AND ul.user_id=$1 AND ul.lists_id=$2", todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userID, listID)

	return err

}

func (r ToDoListPostgres) Update(userID, listID int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Discription != nil {
		setValues = append(setValues, fmt.Sprintf("discription=$%d", argID))
		args = append(args, *input.Discription)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		"UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.lists_id AND ul.lists_id = $%d AND ul.user_id = $%d",
		todoListsTable,
		setQuery,
		usersListsTable,
		argID,
		argID+1,
	)

	args = append(args, listID, userID)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
