//go:generate mockgen -source $GOFILE -typed -destination mocks/$GOFILE -package mocks -typed -mock_names client=Client
package books

import (
	"context"

	pbclient "github.com/micronull/pocketbook-cloud-client"

	"pocketbook-cloud-sync/internal/pkg/domain"
)

type client interface {
	Providers(ctx context.Context, userName string) ([]pbclient.Provider, error)
	Login(ctx context.Context, req pbclient.LoginRequest) (pbclient.Token, error)
	Books(ctx context.Context, token string, limit, offset int) (pbclient.Books, error)
}

type Repository struct {
	client client
	login  string
	pswd   string
}

func New(client client, login, password string) *Repository {
	return &Repository{
		client: client,
		login:  login,
		pswd:   password,
	}
}

func (r Repository) Books(ctx context.Context) ([]domain.Book, error) {
	return nil, nil
}
