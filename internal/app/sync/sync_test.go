package sync_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"pocketbook-cloud-sync/internal/app/sync"
	"pocketbook-cloud-sync/internal/app/sync/mocks"
	"pocketbook-cloud-sync/internal/pkg/domain"
)

func TestApp_Sync(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name             string
		file             string
		expectedDownload assert.BoolAssertionFunc
	}{
		{
			name:             "success",
			file:             "test.txt",
			expectedDownload: assert.True,
		},
		{
			name:             "is exists",
			file:             "exist.txt",
			expectedDownload: assert.False,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			booksMock := mocks.NewBooks(mockCtrl)

			const dir = "testdata"

			var checkDownloader bool

			opts := []sync.Option{
				sync.WithDownloader(func(_ context.Context, url, destination string) error {
					checkDownloader = url == "https://test.link/foo/bar" && destination == dir+"/"+tt.file

					return nil
				}),
			}

			app := sync.New(booksMock, dir, opts...)

			booksMock.EXPECT().
				Books(gomock.Any()).
				Return([]domain.Book{
					{
						FileName: tt.file,
						Link:     "https://test.link/foo/bar",
					},
				}, nil)

			err := app.Sync(t.Context())
			assert.NoError(t, err)

			tt.expectedDownload(t, checkDownloader)
		})
	}
}

func TestApp_Sync_EmptyBooks(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	booksMock := mocks.NewBooks(mockCtrl)

	app := sync.New(booksMock, "testdata")

	booksMock.EXPECT().
		Books(gomock.Any()).
		Return([]domain.Book{}, nil)

	err := app.Sync(t.Context())
	assert.NoError(t, err)
}
