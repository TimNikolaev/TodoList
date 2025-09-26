package postgres

import (
	"fmt"
	"strings"
	"todo-app"

	"github.com/jmoiron/sqlx"
)

type ToDoItemPostgres struct {
	db *sqlx.DB
}

func NewToDoItemPostgres(db *sqlx.DB) *ToDoItemPostgres {
	return &ToDoItemPostgres{db: db}
}

func (r *ToDoItemPostgres) Create(listID int, item todo.ToDoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemID int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) values ($1, $2) RETURNING id", todoItemsTable)

	if err := tx.QueryRow(createItemQuery, item.Title, item.Description).Scan(&itemID); err != nil {
		return 0, err
	}

	createListsItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListsItemsQuery, listID, itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemID, tx.Commit()

}

func (r *ToDoItemPostgres) GetAll(userID, listID int) ([]todo.ToDoItem, error) {
	var items []todo.ToDoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti
	 INNER JOIN %s li on li.item_id = ti.id
		 INNER JOIN %s ul on ul.lists_id = li.list_id
		  WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable,
	)
	if err := r.db.Select(&items, query, listID, userID); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ToDoItemPostgres) GetByID(userID, itemID int) (todo.ToDoItem, error) {
	var item todo.ToDoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti
	 INNER JOIN %s li on li.item_id = ti.id
		 INNER JOIN %s ul on ul.lists_id = li.list_id
		  WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable,
	)
	if err := r.db.Get(&item, query, itemID, userID); err != nil {
		return item, err
	}

	return item, nil
}

func (r *ToDoItemPostgres) Delete(userID, itemID int) error {
	query := fmt.Sprintf("DELETE FROM %s ti USING %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.lists_id AND ul.user_id=$1 AND ti.id=$2", todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userID, itemID)

	return err

}

func (r *ToDoItemPostgres) Update(userID, itemID int, input todo.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *input.Description)
		argID++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argID))
		args = append(args, *input.Done)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(
		`UPDATE %s ti 
     SET %s 
     FROM %s li, %s ul 
     WHERE ti.id = li.item_id 
     AND li.list_id = ul.lists_id 
     AND ul.user_id = $%d 
     AND ti.id = $%d`,
		todoItemsTable,
		setQuery,
		listsItemsTable,
		usersListsTable,
		argID,
		argID+1,
	)

	args = append(args, userID, itemID)

	_, err := r.db.Exec(query, args...)
	return err
}
