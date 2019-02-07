package report

import (
	"github.com/spf13/cobra"
	"github.com/tehcyx/lic/pkg/lic/core"
)

//ReportOptions defines available options for the command
type ReportOptions struct {
	*core.Options
	Upload         bool
	UploadEndpoint string
	SrcPath        string
	HTMLOutput     bool
}

//NewReportOptions creates options with default values
func NewReportOptions(o *core.Options) *ReportOptions {
	return &ReportOptions{Options: o}
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
