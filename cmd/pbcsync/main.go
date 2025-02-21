package main

import (
	"log/slog"
	"os"

	"pocketbook-cloud-sync/internal/pkg/command"
	"pocketbook-cloud-sync/internal/pkg/command/sync"
	"pocketbook-cloud-sync/internal/pkg/command/sync/factory"
)

func main() {
	cmd := command.New()
	cmd.AddCommand("sync", sync.New(factory.Factory))

	if err := cmd.Run(os.Args[1:]); err != nil {
		slog.Error(err.Error())

		os.Exit(1)
	}

	os.Exit(0)
}
