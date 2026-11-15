package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todo_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userID string) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
			INSERT INTO todos (title, completed, user_id)
			VALUES ($1, $2, $3)
			RETURNING id, title, completed, created_at, updated_at, user_id
	`

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, title, completed, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func GetAll(pool *pgxpool.Pool) ([]models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `SELECT * FROM todos ORDER BY created_at DESC`

	var todos []models.Todo = []models.Todo{}

	var rows, err = pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo models.Todo

		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserID,
		)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil

}

func GetTodoByID(pool *pgxpool.Pool, id int) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `SELECT * FROM todos WHERE id = $1`

	var todo models.Todo

	var err = pool.QueryRow(ctx, query, id).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt, &todo.UserID)
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed *bool) (*models.Todo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string
	var args []any

	switch {
	case title != "" && completed != nil:
		query = `UPDATE todos SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING *`
		args = []any{title, *completed, id}
	case title != "" && completed == nil:
		query = `UPDATE todos SET title = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING *`
		args = []any{title, id}
	case completed != nil && title == "":
		query = `UPDATE todos SET completed = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING *`
		args = []any{*completed, id}
	default:
		return nil, errors.New("title or completed must be provided")
	}

	var todo models.Todo
	var err error = pool.QueryRow(ctx, query, args...).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)
	if err != nil {
		return nil, err
	}

	return &todo, nil

}

func DeleteTodo(pool *pgxpool.Pool, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `DELETE FROM todos WHERE id = $1`
	var commandTag, err = pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Todo id '%d' not found", id)
	}

	return nil

}