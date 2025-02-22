package factory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/micronull/pocketbook-cloud-sync/internal/app/sync"
	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command/sync/factory"
	"github.com/micronull/pocketbook-cloud-sync/internal/pkg/command/sync/factory/mocks"
)

func TestFactory(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	cfgMock := mocks.NewMockConfigurator(ctrl)

	cfgMock.EXPECT().ClientID().Return("some client id")
	cfgMock.EXPECT().ClientSecret().Return("some client secret")
	cfgMock.EXPECT().UserName().Return("some user name")
	cfgMock.EXPECT().Password().Return("some password")
	cfgMock.EXPECT().Directory().Return("some directory")

	got := factory.Factory(cfgMock)

	assert.IsType(t, (*sync.App)(nil), got)
}
