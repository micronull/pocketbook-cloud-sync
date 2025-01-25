package sync

import (
	"errors"
	"fmt"
)

type requiredError struct {
	param string
}

func (e requiredError) Error() string {
	return fmt.Sprintf("%s is required", e.param)
}

var errIsNotDirectory = errors.New("is not a directory")

type directoryError struct {
	dir string
	err error
}

func (e directoryError) Error() string {
	return fmt.Sprintf("incorrect directory \"%s\": %s", e.dir, e.err)
}

func (e directoryError) Unwrap() error {
	return e.err
}
