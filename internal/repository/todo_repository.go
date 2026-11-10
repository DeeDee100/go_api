package repository

import (
	"context"
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
