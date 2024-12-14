package sync_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	syncApp "pocketbook-cloud-sync/internal/app/sync"
	"pocketbook-cloud-sync/internal/pkg/command/sync"
	"pocketbook-cloud-sync/internal/pkg/command/sync/mocks"
)

func TestSync_Description(t *testing.T) {
	t.Parallel()

	cmd := sync.New(nil)

	assert.Equal(t, "Uploads missing books to the directory.", cmd.Description())
}

func TestSync_Help(t *testing.T) {
	t.Parallel()

	cmd := sync.New(nil)

	const expected = `Usage of sync:
  -client-id string
    	Client ID of PocketBook Cloud API.
    	Read the readme to find out how to get it.
  -client-secret string
    	Client Secret of PocketBook Cloud API.
    	Read the readme to find out how to get it.
  -debug
    	Enable debug output.
  -dir string
    	Directory to sync files. (default "books")
  -password string
    	Password from your PocketBook Cloud account.
  -username string
    	Username of PocketBook Cloud. Usually it's your email.
`

	assert.Equal(t, expected, cmd.Help())
}

func TestSync_Run(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	appMock := mocks.NewApp(mockCtrl)
	cmd := sync.New(appMock)

	args := []string{
		"-client-id", "some-id",
		"-client-secret", "some-secret",
		"-dir", "some-dir",
		"-password", "some-password",
		"-username", "some-username",
	}

	expected := syncApp.Params{
		ClientID:     "some-id",
		ClientSecret: "some-secret",
		Dir:          "some-dir",
		Password:     "some-password",
		UserName:     "some-username",
	}

	appMock.EXPECT().
		Sync(gomock.Any(), expected).
		Return(nil)

	err := cmd.Run(args)
	require.NoError(t, err)
}

func TestSync_Run_Error(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	appMock := mocks.NewApp(mockCtrl)
	cmd := sync.New(appMock)
	errExpected := errors.New("some error")

	appMock.EXPECT().
		Sync(gomock.Any(), gomock.Any()).
		Return(errExpected)

	err := cmd.Run(nil)
	require.ErrorIs(t, err, errExpected)
}
