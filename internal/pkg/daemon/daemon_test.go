package daemon_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thejerf/slogassert"

	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/daemon"
)

type syncFunc func(ctx context.Context) error

func (s syncFunc) Sync(ctx context.Context) error {
	return s(ctx)
}

func TestDaemon_Sync(t *testing.T) {
	t.Parallel()

	timeout := 100 * time.Millisecond

	var c int

	counter := syncFunc(func(context.Context) error { c++; return nil })

	dn := daemon.New(timeout, counter)

	ctx, cancel := context.WithTimeout(t.Context(), timeout*10)
	t.Cleanup(cancel)

	err := dn.Sync(ctx)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	assert.Greater(t, c, 5)
}

type serverErrorMock struct {
	code int
}

func (s serverErrorMock) Error() string {
	return fmt.Sprintf("server error %d", s.code)
}

func (s serverErrorMock) Code() int {
	return s.code
}

func TestDaemon_Sync_HTTPError_5xx(t *testing.T) {
	orig := slog.Default()

	logassert := slogassert.New(t, slog.LevelWarn, orig.Handler())
	slog.SetDefault(slog.New(logassert))

	t.Cleanup(func() { slog.SetDefault(orig) })

	timeout := 100 * time.Millisecond
	syncMock := syncFunc(func(context.Context) error {
		return serverErrorMock{
			code: rand.N(100) + http.StatusInternalServerError,
		}
	})

	dn := daemon.New(timeout, syncMock)

	ctx, cancel := context.WithTimeout(t.Context(), timeout*10)
	t.Cleanup(cancel)

	err := dn.Sync(ctx)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	logassert.AssertSomeMessage("server error")
}

func TestDaemon_Sync_HTTPError_4xx(t *testing.T) {
	t.Parallel()

	code := rand.N(100) + http.StatusBadRequest

	timeout := 100 * time.Millisecond
	syncMock := syncFunc(func(context.Context) error {
		return serverErrorMock{
			code: code,
		}
	})

	dn := daemon.New(timeout, syncMock)

	err := dn.Sync(t.Context())
	require.ErrorContains(t, err, "call sync: server error "+strconv.Itoa(code))
}

func TestDaemon_Sync_Error(t *testing.T) {
	t.Parallel()

	errExpected := errors.New("expected error")

	counter := syncFunc(func(context.Context) error {
		return errExpected
	})

	dn := daemon.New(time.Millisecond, counter)

	err := dn.Sync(t.Context())
	require.ErrorIs(t, err, errExpected)
}
