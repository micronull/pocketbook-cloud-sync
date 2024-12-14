package sync

import (
	"context"
	"errors"
	"fmt"
)

type Params struct {
	ClientID     string
	ClientSecret string
	UserName     string
	Password     string
	Dir          string
}

type App struct{}

func New() *App {
	return &App{}
}

func (a App) Sync(ctx context.Context, params Params) error {
	if err := validation(params); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return errors.New("implement me")
}

var errRequired = errors.New("required param")

func validation(pms Params) error {
	switch {
	case pms.ClientID == "":
		return fmt.Errorf("%w: client-id", errRequired)
	case pms.ClientSecret == "":
		return fmt.Errorf("%w: client-secret", errRequired)
	case pms.UserName == "":
		return fmt.Errorf("%w: username", errRequired)
	case pms.Password == "":
		return fmt.Errorf("%w: password", errRequired)
	case pms.Dir == "":
		return fmt.Errorf("%w: dir", errRequired)
	}

	return nil
}
