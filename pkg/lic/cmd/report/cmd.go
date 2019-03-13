// Package report implements the `lic report` command.
package report

import (
	"github.com/spf13/cobra"
	"github.com/tehcyx/lic/pkg/lic/core"
)

//Options defines available options for the command
type Options struct {
	*core.Options
	Upload         bool
	UploadEndpoint string
	SrcPath        string
	HTMLOutput     bool
}

//NewReportOptions creates options with default values
func NewReportOptions(o *core.Options) *Options {
	return &Options{Options: o}
}

//NewReportCmd creates a new report command
func NewReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "report",
		Short:   "Creates a report of sources",
		Aliases: []string{"r"},
	}
	return cmd
}
