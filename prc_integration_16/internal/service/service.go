package service

import (
	"context"
	"errors"

	"example.com/prc_integr/internal/models"
	"example.com/prc_integr/internal/repo"
)

type Service struct {
	Notes repo.NoteRepo
}

func (s Service) Create(ctx context.Context, n *models.Note) error {
	if n.Title == "" || n.Content == "" {
		return errors.New("title and content cannot be empty")
	}
	return s.Notes.Create(ctx, n)
}

func (s Service) Get(ctx context.Context, id int64) (models.Note, error) {
	if id <= 0 {
		return models.Note{}, errors.New("invalid id")
	}
	return s.Notes.Get(ctx, id)
}

func (s Service) Update(ctx context.Context, n *models.Note) error {
	if n.ID <= 0 || n.Title == "" || n.Content == "" {
		return errors.New("invalid note data")
	}
	return s.Notes.Update(ctx, n)
}

func (s Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	return s.Notes.Delete(ctx, id)
}

func (s Service) List(ctx context.Context, limit, offset int) ([]models.Note, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.Notes.List(ctx, limit, offset)
}
