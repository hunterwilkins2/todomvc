package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Query int8

const (
	All Query = iota
	Active
	Completed
)

func (q Query) String() string {
	switch q {
	case Active:
		return "active"
	case Completed:
		return "completed"
	default:
		return "all"
	}
}

type Todo struct {
	ID        int
	Name      string
	Completed bool
}

type TodoComponentData struct {
	Todos        []Todo
	Total        int
	Completed    int
	NotCompleted int
	Query        string
}

type TodoModelInterface interface {
	Insert(name string) error
	GetAll(q Query) (*TodoComponentData, error)
	GetTodo(id int) (Todo, error)
	UpdateStatus(completed bool) error
	UpdateTodo(todo *Todo) error
	DeleteCompleted() error
	Delete(id int) error
}

type TodoModel struct {
	db *sql.DB
}

func (m TodoModel) Insert(name string) error {
	query := `
		INSERT INTO todo (name, completed)
		VALUES (?, ?)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, name, false)
	return err
}

func (m TodoModel) GetAll(q Query) (*TodoComponentData, error) {
	query := `
		SELECT id, name, completed
		FROM todo
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &TodoComponentData{
				Total:        0,
				NotCompleted: 0,
				Completed:    0,
				Query:        q.String(),
			}, nil
		default:
			return nil, err
		}
	}

	todoData := TodoComponentData{
		Query: q.String(),
	}
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Name, &todo.Completed)
		if err != nil {
			return nil, err
		}

		switch q {
		case Active:
			if !todo.Completed {
				todoData.Todos = append(todoData.Todos, todo)
			}
		case Completed:
			if todo.Completed {
				todoData.Todos = append(todoData.Todos, todo)
			}
		default:
			todoData.Todos = append(todoData.Todos, todo)
		}
		todoData.Total++
		if todo.Completed {
			todoData.Completed++
		} else {
			todoData.NotCompleted++
		}
	}

	return &todoData, nil
}

func (m TodoModel) GetTodo(id int) (Todo, error) {
	query := `
		SELECT id, name, completed
		FROM todo
		WHERE id = ?
	`

	var todo Todo
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(&todo.ID, &todo.Name, &todo.Completed)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Todo{}, nil
		default:
			return Todo{}, err
		}
	}
	return todo, nil
}

func (m TodoModel) UpdateStatus(completed bool) error {
	query := `
		UPDATE todo
		SET completed = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, completed)
	return err
}

func (m TodoModel) UpdateTodo(todo *Todo) error {
	query := `
		UPDATE todo
		SET name = ?, completed = ?
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, todo.Name, todo.Completed, todo.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		default:
			return err
		}
	}
	return nil
}

func (m TodoModel) DeleteCompleted() error {
	query := `
		DELETE FROM todo
		WHERE completed = true
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		default:
			return err
		}
	}
	return nil
}

func (m TodoModel) Delete(id int) error {
	query := `
		DELETE FROM todo
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil
		default:
			return err
		}
	}
	return nil
}
