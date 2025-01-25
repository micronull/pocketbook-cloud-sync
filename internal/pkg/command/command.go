//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks -mock_names command=Command
package command

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type command interface {
	Run(args []string) error
	Description() string
	Help() string
}

type name = string

type cmnd struct {
	n name
	c command
}

type Command struct {
	cmds       []cmnd
	helpOutput io.Writer
}

func New(opts ...Option) *Command {
	c := &Command{
		cmds:       []cmnd{},
		helpOutput: os.Stdout,
	}

	for i := 0; i < len(opts); i++ {
		opts[i](c)
	}

	c.AddCommand("help", &helpCommand{
		cmds: &c.cmds,
		out:  c.helpOutput,
	})

	return c
}

func (c *Command) AddCommand(name string, cmd command) {
	c.cmds = append(c.cmds, cmnd{
		n: name,
		c: cmd,
	})
}

var errUnknownCommand = errors.New("unknown command")

func (c *Command) Run(args []string) error {
	if len(args) == 0 || isHelpFlag(args[0]) {
		if err := c.help(args); err != nil {
			return fmt.Errorf("help: %w", err)
		}

		return nil
	}

	name := args[0]

	for i := 0; i < len(c.cmds); i++ {
		cmd := c.cmds[i]

		if cmd.n == name {
			if err := cmd.c.Run(args[1:]); err != nil {
				return fmt.Errorf("%s: %w", name, err)
			}

			return nil
		}
	}

	return fmt.Errorf("%w: %s", errUnknownCommand, name)
}

func (c *Command) help(args []string) error {
	ha := []string{"help"}

	if len(args) > 1 {
		ha = append(ha, args[1:]...)
	}

	if err := c.Run(ha); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}

func isHelpFlag(name string) bool {
	name = strings.TrimSpace(name)

	switch name {
	case "-h", "--help", "-help":
		return true
	}

	return false
}

type helpCommand struct {
	cmds *[]cmnd
	out  io.Writer
}

func (h helpCommand) Run(args []string) (err error) {
	var help string

	switch {
	case len(args) > 0:
		if help, err = h.commandHelp(args[0]); err != nil {
			return fmt.Errorf("command help: %s: %w", args[0], err)
		}
	default:
		help = h.Help()
	}

	_, _ = fmt.Fprintln(h.out, help)

	return nil
}

func (h helpCommand) Help() string {
	help := "usage: <command> [<args>]"
	cmds := *h.cmds

	for i := 0; i < len(cmds); i++ {
		help += "\n\t" + cmds[i].n + " - " + cmds[i].c.Description()
	}

	return help
}

func (h helpCommand) commandHelp(cmd string) (string, error) {
	cmds := *h.cmds

	for i := 0; i < len(cmds); i++ {
		if cmds[i].n == cmd {
			return cmds[i].c.Help(), nil
		}
	}

	return "", fmt.Errorf("%w: %s", errUnknownCommand, cmd)
}

func (h helpCommand) Description() string {
	return "Print all available commands with description.\n" +
		"	       Use \"pbcsync help <command>\" for more information about a command."
}
