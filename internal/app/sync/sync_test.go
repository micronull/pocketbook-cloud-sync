package sync_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"pocketbook-cloud-sync/internal/app/sync"
)

func TestApp_Sync_Error_Validation(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name   string
		params sync.Params
		expect string
	}{
		{
			name:   "empty params",
			params: sync.Params{},
			expect: "validate: client-id is required",
		},
		{
			name: "empty client-secret",
			params: sync.Params{
				ClientID:     "some",
				ClientSecret: "",
			},
			expect: "validate: client-secret is required",
		},
		{
			name: "empty username",
			params: sync.Params{
				ClientID:     "some",
				ClientSecret: "some",
			},
			expect: "validate: username is required",
		},
		{
			name: "empty password",
			params: sync.Params{
				ClientID:     "some",
				ClientSecret: "some",
				UserName:     "some",
			},
			expect: "validate: password is required",
		},
		{
			name: "empty dir",
			params: sync.Params{
				ClientID:     "some",
				ClientSecret: "some",
				UserName:     "some",
				Password:     "some",
			},
			expect: "validate: dir is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			app := sync.New()
			err := app.Sync(context.Background(), tt.params)
			require.EqualError(t, err, tt.expect)
		})
	}
}

func TestApp_Sync_Error_NonExistingDirectory(t *testing.T) {
	t.Parallel()

	params := sync.Params{
		ClientID:     "some",
		ClientSecret: "some",
		UserName:     "some",
		Password:     "some",
		Dir:          "some",
	}

	app := sync.New()
	err := app.Sync(context.Background(), params)

	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestApp_Sync_Error_ErrPermissionDirectory(t *testing.T) {
	t.Parallel()

	params := sync.Params{
		ClientID:     "some",
		ClientSecret: "some",
		UserName:     "some",
		Password:     "some",
		Dir:          "test",
	}

	err := os.Mkdir("test", 0o666)
	require.NoError(t, err)

	t.Cleanup(func() { _ = os.Remove("test") })

	app := sync.New()
	err = app.Sync(context.Background(), params)

	require.ErrorIs(t, err, os.ErrPermission)
}
