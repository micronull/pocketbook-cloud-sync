package books_test

import (
	"errors"
	"math/rand/v2"
	"testing"

	pbclient "github.com/micronull/pocketbook-cloud-client"
	"github.com/stretchr/testify/assert"
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

	gomock.InOrder(
		clientMock.EXPECT().
			Books(gomock.Any(), "token-1", 0, 0).
			Return(pbclient.Books{Total: nr}, nil),
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
			}, nil),
	)

	nr = rand.N(100)

	gomock.InOrder(
		clientMock.EXPECT().
			Books(gomock.Any(), "token-2", 0, 0).
			Return(pbclient.Books{Total: nr}, nil),
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
			}, nil),
	)

	got, err := repo.Books(t.Context())
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

var errStub = errors.New("stub error")

func TestRepository_Books_Error_Provider(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	clientMock := mocks.NewClient(mockCtrl)
	repo := books.New(clientMock, "", "")

	clientMock.EXPECT().
		Providers(gomock.Any(), gomock.Any()).
		Return(nil, errStub)

	_, err := repo.Books(t.Context())
	require.ErrorIs(t, err, errStub)
}

func TestRepository_Books_Providers_Empty(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name  string
		items []pbclient.Provider
	}{
		{
			name:  "empty",
			items: []pbclient.Provider{},
		},
		{
			name:  "nil",
			items: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			clientMock := mocks.NewClient(mockCtrl)
			repo := books.New(clientMock, "", "")

			clientMock.EXPECT().
				Providers(gomock.Any(), gomock.Any()).
				Return(tt.items, nil)

			got, err := repo.Books(t.Context())
			require.NoError(t, err)

			assert.Empty(t, got)
		})
	}
}

func TestRepository_Books_Error_Login(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	clientMock := mocks.NewClient(mockCtrl)
	repo := books.New(clientMock, "", "")

	clientMock.EXPECT().
		Providers(gomock.Any(), gomock.Any()).
		Return([]pbclient.Provider{{}}, nil)

	clientMock.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(pbclient.Token{}, errStub)

	_, err := repo.Books(t.Context())
	require.ErrorIs(t, err, errStub)
}

func TestRepository_Books_Error_BooksCount(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	clientMock := mocks.NewClient(mockCtrl)
	repo := books.New(clientMock, "", "")

	clientMock.EXPECT().
		Providers(gomock.Any(), gomock.Any()).
		Return([]pbclient.Provider{{}}, nil)

	clientMock.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(pbclient.Token{}, nil)

	clientMock.EXPECT().
		Books(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(pbclient.Books{}, errStub)

	_, err := repo.Books(t.Context())
	require.ErrorIs(t, err, errStub)
}

func TestRepository_Books_BooksCount_Zero(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	clientMock := mocks.NewClient(mockCtrl)
	repo := books.New(clientMock, "", "")

	clientMock.EXPECT().
		Providers(gomock.Any(), gomock.Any()).
		Return([]pbclient.Provider{{}}, nil)

	clientMock.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(pbclient.Token{}, nil)

	clientMock.EXPECT().
		Books(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(pbclient.Books{}, nil)

	got, err := repo.Books(t.Context())
	require.NoError(t, err)

	assert.Empty(t, got)
}

func TestRepository_Books_Error_Books(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	clientMock := mocks.NewClient(mockCtrl)
	repo := books.New(clientMock, "", "")

	clientMock.EXPECT().
		Providers(gomock.Any(), gomock.Any()).
		Return([]pbclient.Provider{{}}, nil)

	clientMock.EXPECT().
		Login(gomock.Any(), gomock.Any()).
		Return(pbclient.Token{}, nil)

	clientMock.EXPECT().
		Books(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(pbclient.Books{
			Total: 1,
		}, nil)

	clientMock.EXPECT().
		Books(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(pbclient.Books{}, errStub)

	_, err := repo.Books(t.Context())
	require.ErrorIs(t, err, errStub)
}
