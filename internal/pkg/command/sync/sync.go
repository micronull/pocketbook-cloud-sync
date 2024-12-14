package sync

import "errors"

type Sync struct{}

func New() *Sync {
	return &Sync{}
}

func (s Sync) Description() string {
	return "Uploads missing books to the directory."
}

func (s Sync) Help() string {
	return "TODO"
}

func (s Sync) Run(args []string) error {
	return errors.New("not implemented")
}
