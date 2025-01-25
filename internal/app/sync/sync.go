package sync

import (
	"context"
	"errors"
	"fmt"
	"os"
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

func validation(pms Params) error {
	switch {
	case pms.ClientID == "":
		return requiredError{param: "client-id"}
	case pms.ClientSecret == "":
		return requiredError{param: "client-secret"}
	case pms.UserName == "":
		return requiredError{param: "username"}
	case pms.Password == "":
		return requiredError{param: "password"}
	case pms.Dir == "":
		return requiredError{param: "dir"}
	}

	if err := dirCheck(pms.Dir); err != nil {
		return fmt.Errorf("check direcrory: %w", err)
	}

	return nil
}

func dirCheck(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("stat: %w", directoryError{dir: dir, err: err})
	}

	if !info.IsDir() {
		return fmt.Errorf("is dir: %w", directoryError{dir: dir, err: errIsNotDirectory})
	}

	const tmpFile = "test_write_file.txt"

	err = os.WriteFile(dir+"/"+tmpFile, []byte("test"), 0666)
	if err != nil {
		return fmt.Errorf("check write: %w", directoryError{dir: dir, err: err})
	}

	_ = os.Remove(dir + "/" + tmpFile)

	return nil
}
