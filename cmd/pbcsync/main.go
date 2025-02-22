package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command"
	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command/sync"
	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command/sync/factory"
	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command/version"
)

func main() {
	log.SetOutput(os.Stdout)

	cmd := command.New()
	cmd.AddCommand("sync", sync.New(factory.Factory))
	cmd.AddCommand("version", version.New())

	if err := cmd.Run(os.Args[1:]); err != nil {
		slog.Error(err.Error())

		os.Exit(1)
	}

	os.Exit(0)
}
