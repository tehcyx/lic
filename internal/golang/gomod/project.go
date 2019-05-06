package gomod

import (
	"fmt"
	"path"

	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/internal/licensereport"
)

// Collect initiates collection of imports accross given path
func Collect(proj *licensereport.Project, prjPath string) error {
	goModPath := path.Join(prjPath, "go.mod")
	if Exists(goModPath) {
		return ReadImports(proj, goModPath)
	}
	return fmt.Errorf("go.mod does not exist")
}

// Exists checks, that go.mod file exists on given path
func Exists(goModPath string) bool {
	if err := fileop.Exists(goModPath); err == nil {
		return true
	}
	return false
}
