package godep

import "github.com/tehcyx/lic/internal/report"

// Collect initiates collection of imports accross given path
func Collect(proj *report.Project, prjPath string) error {
	return ReadImports(proj, prjPath)
}
