package report

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/tehcyx/lic/internal/fileop"
)

const (
	GoModImportStart = "^require \\(.*"
	GoModImportEnd   = "^\\).*"
	GoModLineImport  = "(\\S+|\\/|\\.)+"

	GoFileImportStart = "^import (\\(|\").*"
	GoFileImportEnd   = "^(\\)|var|func|type).*"
	GoFileLineImport  = "\"(\\S+|\\/|\\.)+\""
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

	cmd.Flags().BoolVarP(&o.Upload, "upload", "u", false, "Upload report to specified report endpoint to capture continuously")
	cmd.Flags().StringVarP(&o.UploadEndpoint, "upload-endpoint", "", "", "URL of the endpoint to report results of the scans")

	cmd.Flags().StringVarP(&o.SrcPath, "src-path", "", "", "Local path of sources to scan")
	cmd.Flags().BoolVarP(&o.HTMLOutput, "html-output", "o", false, "Specifies if results should be published as .html-file stored in current path")

	return cmd
}

//Run runs the command
func (o *ReportOptions) Run() error {
	if o.SrcPath != "" {
		if err := fileop.Exists(o.SrcPath); err != nil {
			log.Printf("Path '%s' does not exist or you don't have the proper access rights.\n", o.SrcPath)
			os.Exit(1)
		}
	} else {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Println("Couldn't get application path, exiting")
			os.Exit(1)
		}
		o.SrcPath = dir
	}

	goModImports, err := fileop.ReadImports(path.Join(o.SrcPath, "go.mod"), GoModImportStart, GoModImportEnd, GoModLineImport)
	if err != nil {
		log.Println("Error reading imports from go.mod file. Reading file tree now.")
	}

	fmt.Printf("%+v\n", goModImports)

	fmt.Println("Check go.mod/go.sum, if they don't exist, open each file and check the imports statement")

	return nil
}
