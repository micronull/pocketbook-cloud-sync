package sync

import (
	"context"
	"errors"

	"pocketbook-cloud-sync/internal/pkg/domain"
)

type books interface {
	Books(ctx context.Context) ([]domain.Book, error)
}

type App struct {
	books books
	dir   string
}

func New(books books, dir string) *App {
	a := &App{
		books: books,
		dir:   dir,
	}

	return a
}

func (a App) Sync(ctx context.Context) error {
	return errors.New("implement me")
}
