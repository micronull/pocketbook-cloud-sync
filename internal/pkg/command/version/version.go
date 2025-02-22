package version

import (
	"fmt"

	"pocketbook-cloud-sync/internal/pkg/version"
)

type Version struct{}

func New() *Version {
	return &Version{}
}

func (v Version) Run([]string) error {
	fmt.Println("version: " + version.Version())

	return nil
}

func (v Version) Description() string {
	return "Print current version."
}

func (v Version) Help() string {
	return "Use subcommand the version."
}
