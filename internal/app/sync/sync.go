package sync

import (
	"context"
	"errors"
)

type Params struct {
	ClientID     string
	ClientSecret string
	UserName     string
	Password     string
	Dir          string
}

type App struct{}

func (a App) Sync(ctx context.Context, params Params) error {
	return errors.New("implement me")
}

func New() *App {
	return &App{}
}
