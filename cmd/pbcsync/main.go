package main

import (
	"log/slog"
	"os"

	syncApp "pocketbook-cloud-sync/internal/app/sync"
	"pocketbook-cloud-sync/internal/pkg/command"
	"pocketbook-cloud-sync/internal/pkg/command/sync"
)

func main() {
	cmd := command.New()
	appSync := syncApp.New()

	cmd.AddCommand("sync", sync.New(appSync))

	if err := cmd.Run(os.Args[1:]); err != nil {
		slog.Error(err.Error())

		os.Exit(1)
	}

	os.Exit(0)
}
