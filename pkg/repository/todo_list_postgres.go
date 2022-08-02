package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
	"todo"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}
func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf("select tl.id,tl.title,tl.description from %s tl inner join %s ul on tl.id = ul.list_id where ul.user_id = $1 ", todoListsTable, usersListsTable)
	if err := r.db.Select(&lists, query, userId); err != nil {
		return nil, err
	}
	return lists, nil
}
func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList
	query := fmt.Sprintf("select tl.id,tl.title,tl.description from %s tl inner join %s ul on tl.id = ul.list_id where ul.user_id = $1 and ul.list_id = $2 ", todoListsTable, usersListsTable)
	if err := r.db.Get(&list, query, userId, listId); err != nil {
		return list, err
	}
	return list, nil

}
func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("delete from %s tl using %s ul  where tl.id = ul.list_id and ul.user_id=$1 and ul.list_id=$2", todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)
	return err

}
func (r *TodoListPostgres) Update(userId, listId int, input todo.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsId := 1
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argsId))
		args = append(args, *input.Title)
		argsId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argsId))
		args = append(args, *input.Description)
		argsId++
	}
	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("update %s tl set %s from %s ul where tl.id = ul.list_id and ul.list_id=$%d and ul.user_id=$%d",
		todoListsTable, setQuery, usersListsTable, argsId, argsId+1)
	args = append(args, listId, userId)
	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)
	_, err := r.db.Exec(query, args...)
	return err
}
