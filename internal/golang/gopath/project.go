package gopath

import (
	"fmt"

	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/internal/report"
)

// Collect initiates collection of imports accross given path
func Collect(proj *report.Project, prjPath string) error {
	if Exists(prjPath) {
		return ReadImports(proj, prjPath)
	}
	return fmt.Errorf("project folder does not exist")
}

// Exists checks that path exists
func Exists(prjPath string) bool {
	if err := fileop.Exists(prjPath); err == nil {
		return true
	}
	return false
}
