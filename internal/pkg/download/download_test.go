package download_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/download"
)

func TestDownload(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/test.txt" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("test"))
		}
	}))

	t.Cleanup(srv.Close)

	err := download.Download(t.Context(), srv.URL+"/test.txt", "test_dest.txt")
	require.NoError(t, err)

	t.Cleanup(func() { _ = os.Remove("test_dest.txt") })

	_, err = os.Stat("test_dest.txt")
	require.NoError(t, err)
}

func TestDownload_StatusNotOk(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	t.Cleanup(srv.Close)

	err := download.Download(t.Context(), srv.URL+"/test.txt", "test_dest.txt")
	require.ErrorContains(t, err, "418 I'm a teapot")
}

func TestDownload_ContextCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := download.Download(ctx, "http://foo", "bar")
	require.ErrorIs(t, err, context.Canceled)
}
