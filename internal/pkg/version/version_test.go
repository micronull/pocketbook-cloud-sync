package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/version"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "undefined", version.Version())
}
