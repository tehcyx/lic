package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tehcyx/lic/pkg/lic/core"
)

//Version contains the kyma-cli binary version injected by the build system
var Version string

//VersionOptions defines available options for the command
type VersionOptions struct {
	*core.Options
}

//NewVersionOptions creates options with default values
func NewVersionOptions(o *core.Options) *VersionOptions {
	return &VersionOptions{Options: o}
}

//NewVersionCmd creates a new version command
func NewVersionCmd(o *VersionOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Version of the lic CLI",
		Long:  `Prints the version of lic CLI`,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.Run() },
	}

	return cmd
}

//Run runs the command
func (o *VersionOptions) Run() error {

	version := Version
	if version == "" {
		version = "N/A"
	}
	fmt.Printf("lic CLI version: %s\n", version)

	return nil
}
