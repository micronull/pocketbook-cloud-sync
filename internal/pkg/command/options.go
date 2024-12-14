package command

import "io"

type Option func(*Command)

func WithHelpOutput(output io.Writer) Option {
	return func(c *Command) {
		c.helpOutput = output
	}
}
