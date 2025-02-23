package sync_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command/sync"
	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command/sync/factory"
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
  -daemon
    	Enable daemon mode. Use the daemon-timeout flag for setting sync interval.
  -daemon-timeout duration
    	Timeout for sync operation. 
    	Used only daemon mode. (default 24h0m0s)
  -debug
    	Enable debug output.
  -dir string
    	Directory to sync files. (default "books")
  -env
    	Enable environment variables mode.
    	Ignores all command-line flags and loads values from environment variables:
    	PBC_CLIENT_ID as -client-id
    	PBC_CLIENT_SECRET as -client-secret
    	PBC_USERNAME as -username
    	PBC_PASSWORD as -password
    	DEBUG as -debug
    	DIR as -dir
    	DAEMON as -daemon
    	DAEMON_TIMEOUT as -daemon-timeout
  -password string
    	Password from your PocketBook Cloud account.
  -username string
    	Username of PocketBook Cloud. Usually it's your email.
`

	assert.Equal(t, expected, cmd.Help())
}

func defaultArgs() []string {
	return []string{
		"-client-id", "some-id",
		"-client-secret", "some-secret",
		"-dir", "testdata",
		"-password", "some-password",
		"-username", "some-username",
	}
}

type mockSync struct {
	mock.Mock
}

func (m *mockSync) Sync(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func TestSync_Run(t *testing.T) {
	t.Parallel()
	_ = os.Mkdir("testdata", 0777)

	appMock := &mockSync{}
	cmd := sync.New(func(config factory.Configurator) factory.Synchronizer {
		assert.Equal(t, "some-id", config.ClientID())
		assert.Equal(t, "some-secret", config.ClientSecret())
		assert.Equal(t, "testdata", config.Directory())
		assert.Equal(t, "some-password", config.Password())
		assert.Equal(t, "some-username", config.UserName())

		return appMock
	})

	args := defaultArgs()

	appMock.On("Sync", mock.Anything).Return(nil)

	err := cmd.Run(args)
	require.NoError(t, err)

	appMock.AssertExpectations(t)
}

func TestSync_Run_Env(t *testing.T) {
	_ = os.Mkdir("testdata", 0777)

	appMock := &mockSync{}
	cmd := sync.New(func(config factory.Configurator) factory.Synchronizer {
		assert.Equal(t, "some-id from env", config.ClientID())
		assert.Equal(t, "some-secret from env", config.ClientSecret())
		assert.Equal(t, "testdata", config.Directory())
		assert.Equal(t, "some-password from env", config.Password())
		assert.Equal(t, "some-username from env", config.UserName())

		return appMock
	})

	args := []string{
		"-client-id", "some-id",
		"-client-secret", "some-secret",
		"-dir", "foobar",
		"-password", "some-password",
		"-username", "some-username",
		"-env",
	}

	t.Setenv("PBC_CLIENT_ID", "some-id from env")
	t.Setenv("PBC_CLIENT_SECRET", "some-secret from env")
	t.Setenv("PBC_USERNAME", "some-username from env")
	t.Setenv("PBC_PASSWORD", "some-password from env")
	t.Setenv("DIR", "testdata")

	appMock.On("Sync", mock.Anything).Return(nil)

	err := cmd.Run(args)
	require.NoError(t, err)

	appMock.AssertExpectations(t)
}

func TestSync_Run_Error(t *testing.T) {
	t.Parallel()
	_ = os.Mkdir("testdata", 0777)

	appMock := &mockSync{}
	cmd := sync.New(func(factory.Configurator) factory.Synchronizer {
		return appMock
	})
	errExpected := errors.New("some error")

	appMock.On("Sync", mock.Anything).Return(errExpected)

	args := defaultArgs()

	err := cmd.Run(args)
	require.ErrorIs(t, err, errExpected)
}

func TestSync_Run_Error_Validation(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name   string
		args   []string
		expect string
	}{
		{
			name:   "empty args",
			args:   []string{},
			expect: "validate: client-id is required",
		},
		{
			name:   "empty client-secret",
			args:   []string{"-client-id", "some-id"},
			expect: "validate: client-secret is required",
		},
		{
			name: "empty username",
			args: []string{
				"-client-id", "some-id",
				"-client-secret", "some-secret",
			},
			expect: "validate: username is required",
		},
		{
			name: "empty password",
			args: []string{
				"-client-id", "some-id",
				"-client-secret", "some-secret",
				"-username", "some-username",
			},
			expect: "validate: password is required",
		},
		{
			name: "empty dir",
			args: []string{
				"-client-id", "some-id",
				"-client-secret", "some-secret",
				"-username", "some-username",
				"-dir", "",
				"-password", "some-password",
			},
			expect: "validate: dir is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := sync.New(nil)

			err := app.Run(tt.args)
			require.EqualError(t, err, tt.expect)
		})
	}
}

func TestSync_Run_Error_NonExistingDirectory(t *testing.T) {
	t.Parallel()

	args := []string{
		"-client-id", "some-id",
		"-client-secret", "some-secret",
		"-dir", "some-dir",
		"-password", "some-password",
		"-username", "some-username",
	}

	app := sync.New(nil)

	err := app.Run(args)

	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestSync_Run_Error_ErrPermissionDirectory(t *testing.T) {
	t.Parallel()

	args := []string{
		"-client-id", "some-id",
		"-client-secret", "some-secret",
		"-dir", "test",
		"-password", "some-password",
		"-username", "some-username",
	}

	err := os.Mkdir("test", 0o666)
	require.NoError(t, err)

	t.Cleanup(func() { _ = os.Remove("test") })

	app := sync.New(nil)
	err = app.Run(args)

	require.ErrorIs(t, err, os.ErrPermission)
}
