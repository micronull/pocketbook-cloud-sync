//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks -typed -exclude_interfaces Synchronizer
package factory

import (
	"context"
	"net/http"

	pc "github.com/micronull/pocketbook-cloud-client"

	"pocketbook-cloud-sync/internal/app/sync"
	"pocketbook-cloud-sync/internal/pkg/repository/books"
)

type Synchronizer interface {
	Sync(ctx context.Context) error
}

type Configurator interface {
	ClientID() string
	ClientSecret() string
	UserName() string
	Password() string
	Directory() string
}

func Factory(config Configurator) Synchronizer {
	return sync.New(
		books.New(
			pc.New(
				pc.WithHTTPClient(&http.Client{Transport: http.DefaultTransport}),
				pc.WithClientID(config.ClientID()),
				pc.WithClientSecret(config.ClientSecret()),
			),
			config.UserName(),
			config.Password(),
		),
		config.Directory(),
	)
}
