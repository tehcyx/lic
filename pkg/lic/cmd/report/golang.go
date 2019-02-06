package report

import (
	"fmt"

	"github.com/spf13/cobra"
)

//NewGolangReportCmd creates a new report command
func NewGolangReportCmd(o *ReportOptions) *cobra.Command {
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
func (o *ReportOptions) Run() error {

	fmt.Println("Open all *.go files on current path, ")

	return nil
}
