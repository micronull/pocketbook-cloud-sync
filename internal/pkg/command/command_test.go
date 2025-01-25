package command_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"pocketbook-cloud-sync/internal/pkg/command"
	"pocketbook-cloud-sync/internal/pkg/command/mocks"
)

func TestCommand_Run(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	commandMock := mocks.NewCommand(mockCtrl)
	buf := &bytes.Buffer{}

	cmd := command.New(command.WithHelpOutput(buf))
	cmd.AddCommand("test", commandMock)

	commandMock.EXPECT().
		Run([]string{"some"}).
		Return(nil)

	err := cmd.Run([]string{"test", "some"})
	require.NoError(t, err)
}

func TestCommand_Run_Error(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	commandMock := mocks.NewCommand(mockCtrl)
	buf := &bytes.Buffer{}

	cmd := command.New(command.WithHelpOutput(buf))
	cmd.AddCommand("test", commandMock)

	errExpected := errors.New("some error")

	commandMock.EXPECT().
		Run(gomock.Any()).
		Return(errExpected)

	err := cmd.Run([]string{"test"})
	require.ErrorIs(t, err, errExpected)
	require.ErrorContains(t, err, "test: some error")
}

func TestCommand_Run_Help(t *testing.T) {
	t.Parallel()

	tests := [...]struct {
		name string
		args []string
	}{
		{"default", []string{}},
		{"help", []string{"help"}},
		{"help help", []string{"help", "help"}},
		{"--help", []string{"--help"}},
		{"-h", []string{"-h"}},
		{"-help", []string{"-help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCtrl := gomock.NewController(t)
			commandMock := mocks.NewCommand(mockCtrl)
			buf := &bytes.Buffer{}

			cmd := command.New(command.WithHelpOutput(buf))
			cmd.AddCommand("test", commandMock)

			commandMock.EXPECT().
				Description().
				Return("Test command mock.")

			err := cmd.Run(tt.args)
			require.NoError(t, err)

			const expected = `usage: <command> [<args>]
	help - Print all available commands with description.
	       Use "pbcsync help <command>" for more information about a command.
	test - Test command mock.
`

			assert.Equal(t, expected, buf.String())
		})
	}
}

func TestCommand_Run_Help_Subcommand(t *testing.T) {
	t.Parallel()

	mockCtrl := gomock.NewController(t)
	commandMock := mocks.NewCommand(mockCtrl)
	buf := &bytes.Buffer{}

	cmd := command.New(command.WithHelpOutput(buf))
	cmd.AddCommand("test", commandMock)

	commandMock.EXPECT().
		Help().
		Return("Test help.")

	err := cmd.Run([]string{"help", "test"})
	require.NoError(t, err)

	const expected = `Test help.
`

	assert.Equal(t, expected, buf.String())
}
