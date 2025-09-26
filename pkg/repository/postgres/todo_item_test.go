package postgres

import (
	"errors"
	"log"
	"testing"
	"todo-app"

	"github.com/stretchr/testify/assert"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	type args struct {
		listID int
		item   todo.ToDoItem
	}

	type mockBehavior func(m sqlxmock.Sqlmock, args args, id int)

	testTable := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		itemID       int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				listID: 1,
				item: todo.ToDoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			mockBehavior: func(m sqlxmock.Sqlmock, args args, itemID int) {
				m.ExpectBegin()

				rows := sqlxmock.NewRows([]string{"id"}).AddRow(itemID)
				m.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				m.ExpectExec("INSERT INTO lists_items").
					WithArgs(args.listID, itemID).
					WillReturnResult(sqlxmock.NewResult(1, 1))

				m.ExpectCommit()
			},
			itemID: 2,
		},
		{
			name: "Empty Fields",
			args: args{
				listID: 1,
				item: todo.ToDoItem{
					Title:       "",
					Description: "test description",
				},
			},
			mockBehavior: func(m sqlxmock.Sqlmock, args args, itemID int) {
				m.ExpectBegin()

				rows := sqlxmock.NewRows([]string{"id"}).AddRow(itemID).RowError(1, errors.New("some error"))
				m.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				m.ExpectRollback()
			},
			wantErr: true,
		},

		{
			name: "2nd Insert Error",
			args: args{
				listID: 1,
				item: todo.ToDoItem{
					Title:       "",
					Description: "test description",
				},
			},
			mockBehavior: func(m sqlxmock.Sqlmock, args args, itemID int) {
				m.ExpectBegin()

				rows := sqlxmock.NewRows([]string{"id"}).AddRow(itemID).RowError(1, errors.New("some error"))
				m.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				m.ExpectExec("INSERT INTO lists_items").
					WithArgs(args.listID, itemID).
					WillReturnError(errors.New("some error"))

				m.ExpectRollback()
			},
			wantErr: true,
		},
	}

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewToDoItemPostgres(db)

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(mock, testCase.args, testCase.itemID)

			got, err := r.Create(testCase.args.listID, testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.itemID, got)
			}
		})
	}
}
