package main

import (
	"context"
	"database/sql"
	"time"
)

// Task — модель для сканирования результатов SELECT
type Task struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo { return &Repo{DB: db} }

// CreateTask — параметризованный INSERT с возвратом id
func (r *Repo) CreateTask(ctx context.Context, title string) (int, error) {
	var id int
	const q = `INSERT INTO tasks (title) VALUES ($1) RETURNING id;`
	err := r.DB.QueryRowContext(ctx, q, title).Scan(&id)
	return id, err
}

// ListTasks — базовый SELECT всех задач (демо для занятия)
func (r *Repo) ListTasks(ctx context.Context) ([]Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks ORDER BY id;`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) ListDone(ctx context.Context, done bool) ([]Task, error) {
	const q = `SELECT * FROM tasks WHERE done = ($1);`
	rows, err := r.DB.QueryContext(ctx, q, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) FindByID(ctx context.Context, id int) (*Task, error) {
	const q = `SELECT * FROM tasks WHERE id = ($1);`
	rows, err := r.DB.QueryContext(ctx, q, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var t Task
	if rows.Next() {
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
	} else {
		return nil, sql.ErrNoRows // Возвращаем ошибку, если строка не найдена
	}

	return &t, rows.Err()
}

func (r *Repo) CreateMany(ctx context.Context, titles []string, dones []bool) error {
	// Начало транзакции
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if rerr := tx.Rollback(); rerr != nil && err == nil {
			err = rerr
		}
	}()

	const q = `INSERT INTO tasks (title, done) VALUES ($1, $2);`

	for i := 0; i < len(titles); i++ {
		_, err := tx.ExecContext(ctx, q, titles[i], dones[i])
		if err != nil {
			return err
		}
	}

	// Фиксация транзакции
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
