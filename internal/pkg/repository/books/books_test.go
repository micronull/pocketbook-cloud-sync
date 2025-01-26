package books_test

import (
	"context"
	"math/rand/v2"
	"testing"

	pbclient "github.com/micronull/pocketbook-cloud-client"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"pocketbook-cloud-sync/internal/pkg/domain"
	"pocketbook-cloud-sync/internal/pkg/repository/books"
	"pocketbook-cloud-sync/internal/pkg/repository/books/mocks"
)

func TestRepository_Books(t *testing.T) {
	t.Parallel()

	const (
		login = "some login"
		pwd   = "some password"
	)

	mockCtrl := gomock.NewController(t)
	clientMock := mocks.NewClient(mockCtrl)
	repo := books.New(clientMock, login, pwd)

	clientMock.EXPECT().
		Providers(gomock.Any(), login).
		Return([]pbclient.Provider{
			{
				Alias:  "provider-1",
				ShopID: "1",
			},
			{
				Alias:  "provider-2",
				ShopID: "2",
			},
		}, nil)

	clientMock.EXPECT().
		Login(gomock.Any(), pbclient.LoginRequest{
			ShopID:   "1",
			UserName: login,
			Password: pwd,
			Provider: "provider-1"},
		).
		Return(pbclient.Token{AccessToken: "token-1"}, nil)

	clientMock.EXPECT().
		Login(gomock.Any(), pbclient.LoginRequest{
			ShopID:   "2",
			UserName: login,
			Password: pwd,
			Provider: "provider-2"},
		).
		Return(pbclient.Token{AccessToken: "token-2"}, nil)

	nr := rand.N(100)

	clientMock.EXPECT().
		Books(gomock.Any(), "token-1", 0, 0).
		Return(pbclient.Books{Total: nr}, nil)

	clientMock.EXPECT().
		Books(gomock.Any(), "token-1", nr, 0).
		Return(pbclient.Books{
			Total: 1,
			Books: []pbclient.Book{
				{
					Link: "https://example.com/first.txt",
					Name: "first.txt",
				},
			},
		}, nil)

	nr = rand.N(100)

	clientMock.EXPECT().
		Books(gomock.Any(), "token-2", 0, 0).
		Return(pbclient.Books{Total: nr}, nil)

	clientMock.EXPECT().
		Books(gomock.Any(), "token-2", nr, 0).
		Return(pbclient.Books{
			Total: 1,
			Books: []pbclient.Book{
				{
					Link: "https://example.com/second.txt",
					Name: "second.txt",
				},
			},
		}, nil)

	got, err := repo.Books(context.Background())
	require.NoError(t, err)

	expected := []domain.Book{
		{
			FileName: "first.txt",
			Link:     "https://example.com/first.txt",
		},
		{
			FileName: "second.txt",
			Link:     "https://example.com/second.txt",
		},
	}

	require.Equal(t, expected, got)
}
