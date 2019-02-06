package report

import (
	"fmt"

	"github.com/spf13/cobra"
)

//GolangReportOptions defines available options for the command
type GolangReportOptions struct {
	*core.Options
	Upload         bool
	UploadEndpoint string
	SrcPath        string
	HTMLOutput     bool
}

//NewGolangReportOptions creates options with default values
func NewGolangReportOptions(o *core.Options) *GolangReportOptions {
	return &GolangReportOptions{Options: o}
}

//NewGolangReportCmd creates a new report command
func NewGolangReportCmd(o *GolangReportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "golang",
		Short:   "Generates a report of current working directory",
		Long:    `Taking in consideration the source on the current path and checking for all licenses, generating a report output in the shell.`,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
		Aliases: []string{"go"},
	}

	return cmd
}

//Run runs the command
func (o *GolangReportOptions) Run() error {

	fmt.Println("Open all *.go files on current path, ")

	return nil
}
