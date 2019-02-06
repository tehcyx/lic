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
func NewReportCmd(o *ReportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "report",
		Short:   "Creates a report of sources",
		Aliases: []string{"r"},
	}

	cmd.Flags().BoolVarP(&o.Upload, "upload", "u", false, "Upload report to specified report endpoint to capture continuously")
	cmd.Flags().StringVarP(&o.UploadEndpoint, "upload-endpoint", "", "", "URL of the endpoint to report results of the scans")

	cmd.Flags().StringVarP(&o.SrcPath, "src-path", "", "", "Local path of sources to scan")
	cmd.Flags().BoolVarP(&o.HTMLOutput, "html-output", "o", false, "Specifies if results should be published as .html-file stored in current path")

	return cmd
}
