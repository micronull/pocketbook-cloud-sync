package sync_test

import (
	"context"
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
			expect: "validate: required param: client-id",
		},
		{
			name: "empty client-secret",
			params: sync.Params{
				ClientID:     "some",
				ClientSecret: "",
			},
			expect: "validate: required param: client-secret",
		},
		{
			name: "empty username",
			params: sync.Params{
				ClientID:     "some",
				ClientSecret: "some",
			},
			expect: "validate: required param: username",
		},
		{
			name: "empty password",
			params: sync.Params{
				ClientID:     "some",
				ClientSecret: "some",
				UserName:     "some",
			},
			expect: "validate: required param: password",
		},
		{
			name: "empty dir",
			params: sync.Params{
				ClientID:     "some",
				ClientSecret: "some",
				UserName:     "some",
				Password:     "some",
			},
			expect: "validate: required param: dir",
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

//func TestApp_Sync(t *testing.T) {
//	t.Parallel()
//
//	fstest.MapFS{}
//}
