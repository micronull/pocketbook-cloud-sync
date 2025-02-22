//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks -mock_names app=App
package sync

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command/sync/factory"
	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/daemon"
)

const daemonTimeoutDefault = time.Hour * 24

type factorySynchronizer func(config factory.Configurator) factory.Synchronizer

type Sync struct {
	flags   *flag.FlagSet
	cfg     *config
	factory factorySynchronizer
}

func New(factory factorySynchronizer) *Sync {
	flags := flag.NewFlagSet("sync", flag.ContinueOnError)

	cfg := &config{}

	flags.BoolVar(&cfg.env, "env", false, "Enable environment variables mode.\n"+
		"Ignores all command-line flags and loads values from environment variables:\n"+
		"PBC_CLIENT_ID as -client-id\n"+
		"PBC_CLIENT_SECRET as -client-secret\n"+
		"PBC_USERNAME as -username\n"+
		"PBC_PASSWORD as -password\n"+
		"DEBUG as -debug\n"+
		"DIR as -dir\n"+
		"DAEMON as -daemon\n"+
		"DAEMON_TIMEOUT as -daemon-timeout")

	flags.StringVar(&cfg.clientID, "client-id", "", "Client ID of PocketBook Cloud API.\n"+
		"Read the readme to find out how to get it.")

	flags.StringVar(&cfg.clientSecret, "client-secret", "", "Client Secret of PocketBook Cloud API.\n"+
		"Read the readme to find out how to get it.")

	flags.StringVar(&cfg.userName, "username", "", "Username of PocketBook Cloud. Usually it's your email.")

	flags.StringVar(&cfg.password, "password", "", "Password from your PocketBook Cloud account.")

	flags.StringVar(&cfg.dir, "dir", "books", "Directory to sync files.")

	flags.BoolVar(&cfg.debug, "debug", false, "Enable debug output.")

	flags.BoolVar(&cfg.daemon, "daemon", false, "Enable daemon mode. Use the daemon-timeout flag for setting sync interval.")

	flags.DurationVar(&cfg.daemonTimeout, "daemon-timeout", daemonTimeoutDefault, "Timeout for sync operation. \n"+
		"Used only daemon mode.")

	return &Sync{
		flags:   flags,
		cfg:     cfg,
		factory: factory,
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
	var err error

	if err = s.flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}

		return fmt.Errorf("flag parse: %v", err)
	}

	if s.cfg.env {
		if s.cfg, err = loadConfigFromEnv(); err != nil {
			return fmt.Errorf("load config from env: %v", err)
		}
	}

	if err := validation(*s.cfg); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if s.cfg.debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("debug enabled")
	}

	ctx, cancel := createContext()
	defer cancel()

	app := s.factory(s.cfg)

	slog.Info("Welcome! I will be glad to receive your star: https://github.com/micronull/pocketbook-cloud-client")

	if s.cfg.daemon {
		slog.Debug("starting daemon mode", "timeout", s.cfg.daemonTimeout)
		app = factory.Synchronizer(daemon.New(s.cfg.daemonTimeout, app))
	}

	if err := app.Sync(ctx); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}

func createContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)
}

func validation(cfg config) error {
	switch {
	case cfg.clientID == "":
		return requiredError{param: "client-id"}
	case cfg.clientSecret == "":
		return requiredError{param: "client-secret"}
	case cfg.userName == "":
		return requiredError{param: "username"}
	case cfg.password == "":
		return requiredError{param: "password"}
	case cfg.dir == "":
		return requiredError{param: "dir"}
	}

	if err := dirCheck(cfg.dir); err != nil {
		return fmt.Errorf("check directory: %w", err)
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

func loadConfigFromEnv() (*config, error) {
	cfg := &config{
		daemonTimeout: daemonTimeoutDefault,
	}

	var err error

	cfg.clientID = os.Getenv("PBC_CLIENT_ID")
	cfg.clientSecret = os.Getenv("PBC_CLIENT_SECRET")
	cfg.userName = os.Getenv("PBC_USERNAME")
	cfg.password = os.Getenv("PBC_PASSWORD")
	cfg.debug = os.Getenv("DEBUG") == "true"
	cfg.daemon = os.Getenv("DAEMON") == "true"

	if dt := os.Getenv("DAEMON_TIMEOUT"); dt != "" {
		if cfg.daemonTimeout, err = time.ParseDuration(dt); err != nil {
			return nil, fmt.Errorf("set daemon timeout: %w", err)
		}
	}

	cfg.dir = os.Getenv("DIR")

	return cfg, err
}
