//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks -mock_names app=App
package sync

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"

	"pocketbook-cloud-sync/internal/app/sync"
)

type app interface {
	Sync(ctx context.Context, params sync.Params) error
}

type Sync struct {
	flags *flag.FlagSet
	cfg   *config
	app   app
}

type config struct {
	clientID     string
	clientSecret string
	userName     string
	password     string
	debug        bool
	dir          string
}

func New(app app) *Sync {
	flags := flag.NewFlagSet("sync", flag.ContinueOnError)

	cfg := &config{}

	flags.StringVar(&cfg.clientID, "client-id", "", "Client ID of PocketBook Cloud API.\n"+
		"Read the readme to find out how to get it.")

	flags.StringVar(&cfg.clientSecret, "client-secret", "", "Client Secret of PocketBook Cloud API.\n"+
		"Read the readme to find out how to get it.")

	flags.StringVar(&cfg.userName, "username", "", "Username of PocketBook Cloud. Usually it's your email.")

	flags.StringVar(&cfg.password, "password", "", "Password from your PocketBook Cloud account.")

	flags.StringVar(&cfg.dir, "dir", "books", "Directory to sync files.")

	flags.BoolVar(&cfg.debug, "debug", false, "Enable debug output.")

	return &Sync{
		flags: flags,
		cfg:   cfg,
		app:   app,
	}
}

func (s Sync) Description() string {
	return "Uploads missing books to the directory."
}

func (s Sync) Help() string {
	buf := &bytes.Buffer{}

	s.flags.SetOutput(buf)
	s.flags.Usage()

	return buf.String()
}

func (s Sync) Run(args []string) error {
	if err := s.flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}

		return fmt.Errorf("flag parse: %v", err)
	}

	if s.cfg.debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	pms := sync.Params{
		ClientID:     s.cfg.clientID,
		ClientSecret: s.cfg.clientSecret,
		UserName:     s.cfg.userName,
		Password:     s.cfg.password,
		Dir:          s.cfg.dir,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := s.app.Sync(ctx, pms); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}
