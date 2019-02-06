package cmd

import (
	"time"

	"github.com/tehcyx/lic/pkg/cmd/report"
	"github.com/tehcyx/lic/pkg/core"
	"github.com/spf13/cobra"
)

const (
	sleep = 10 * time.Second
)

//NewLicCmd creates a new kyma CLI command
func NewLicCmd(o *core.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lic",
		Short: "Enables scanning of sources.",
		Long: `Lic is an easily extensible & flexible report generator to statically analyse your local sources
and create a report on the fly or upload said report to a server.
Find more information at: https://github.com/tehcyx/lic
`,
	}

	cmd.PersistentFlags().BoolVarP(&o.Verbose, "verbose", "v", false, "verbose output")

	versionCmd := NewVersionCmd(NewVersionOptions(o))
	cmd.AddCommand(versionCmd)

	reportCmd := NewReportCmd()
	cmd.AddCommand(reportCmd)

	reportGolangCmd := report.NewReportGolangCmd(report.NewGolangReportOptions(o))
	reportCmd.AddCommand(reportGolangCmd)

	return cmd
}