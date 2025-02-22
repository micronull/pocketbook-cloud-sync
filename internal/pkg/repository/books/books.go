//go:generate mockgen -source $GOFILE -typed -destination mocks/$GOFILE -package mocks -typed -mock_names client=Client
package books

import (
	"context"
	"fmt"
	"log/slog"

	pbclient "github.com/micronull/pocketbook-cloud-client"

	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/domain"
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
	providers, err := r.client.Providers(ctx, r.login)
	if err != nil {
		return nil, fmt.Errorf("get providers: %w", err)
	}

	books := make([]domain.Book, 0)

	for i := 0; i < len(providers); i++ {
		provider := providers[i]

		token, err := r.client.Login(ctx, pbclient.LoginRequest{
			ShopID:   provider.ShopID,
			UserName: r.login,
			Password: r.pswd,
			Provider: provider.Alias,
		})
		if err != nil {
			return nil, fmt.Errorf("login: %w", err)
		}

		pbooks, err := r.client.Books(ctx, token.AccessToken, 0, 0)
		if err != nil {
			return nil, fmt.Errorf("get books count: %w", err)
		}

		slog.Debug("books",
			"total", pbooks.Total,
			"provider_shop_id", provider.ShopID,
			"provider_name", provider.Name,
			"provider_alias", provider.Alias,
		)

		if pbooks.Total == 0 {
			continue
		}

		pbooks, err = r.client.Books(ctx, token.AccessToken, pbooks.Total, 0)
		if err != nil {
			return nil, fmt.Errorf("get books: %w", err)
		}

		for n := 0; n < len(pbooks.Books); n++ {
			pbook := pbooks.Books[n]

			if pbook.Link == "" {
				slog.Warn("book link is empty", "book_id", pbook.ID, "book_name", pbook.Name)

				continue
			}

			books = append(books, domain.Book{
				FileName: pbook.Name,
				Link:     pbook.Link,
			})
		}
	}

	return books, nil
}
