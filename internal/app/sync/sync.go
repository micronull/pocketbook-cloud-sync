package sync

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/unicode/norm"

	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/domain"
	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/download"
)

//go:generate mockgen -source $GOFILE -typed -destination mocks/$GOFILE -package mocks -typed -mock_names books=Books
type books interface {
	Books(ctx context.Context) ([]domain.Book, error)
}

type App struct {
	books      books
	dir        string
	downloader func(ctx context.Context, url, destination string) error
}

func New(books books, dir string, opts ...Option) *App {
	a := &App{
		books:      books,
		dir:        strings.TrimRight(dir, string(os.PathSeparator)),
		downloader: download.Download,
	}

	for _, o := range opts {
		o(a)
	}

	return a
}

func (a App) Sync(ctx context.Context) error {
	exist, err := readDir(a.dir)
	if err != nil {
		return fmt.Errorf("read exists files: %w", err)
	}

	slog.Info("start sync")

	bks, err := a.books.Books(ctx)
	if err != nil {
		return fmt.Errorf("get books: %w", err)
	}

	if len(bks) == 0 {
		slog.Warn("no books found")

		return nil
	}

	var skipped int

	for _, bk := range bks {
		if exist.exist(bk.FileName) {
			slog.Debug("skipped book, this is exists", "name", bk.FileName)

			skipped++

			continue
		}

		path := filepath.Join(a.dir, bk.FileName)

		slog.Debug("download", "file_name", bk.FileName, "path", path, "link", bk.Link)

		if err = a.downloader(ctx, bk.Link, path); err != nil {
			return fmt.Errorf("download %s: %w", bk.FileName, err)
		}
	}

	slog.Info("finished sync", "total", len(bks), "skipped", skipped)

	return nil
}

func readDir(dir string) (files, error) {
	f := files{}

	fls, err := os.ReadDir(dir)
	if err != nil {
		return files{}, fmt.Errorf("read dir: %w", err)
	}

	f.f = make(map[string]struct{}, len(fls))

	for _, fl := range fls {
		if !fl.IsDir() {
			f.f[norm.NFC.String(fl.Name())] = struct{}{}
		}
	}

	return f, nil
}

type files struct {
	f map[string]struct{}
}

func (e files) exist(file string) bool {
	_, ok := e.f[norm.NFC.String(file)]

	return ok
}
