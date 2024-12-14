package main

import (
	"log/slog"
	"os"

	"pocketbook-cloud-sync/internal/pkg/command"
	"pocketbook-cloud-sync/internal/pkg/command/sync"
)

func main() {
	cmd := command.New()

	cmd.AddCommand("sync", sync.New())

	if err := cmd.Run(os.Args[1:]); err != nil {
		slog.Error(err.Error())

		os.Exit(1)
	}

	os.Exit(0)
}
