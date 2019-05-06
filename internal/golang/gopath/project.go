package gopath

import (
	"fmt"
	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/internal/licensereport"
)

func Collect(proj *licensereport.Project, prjPath string) error {
	if Exists(prjPath) {
		return ReadImports(proj, prjPath)
	}
	return fmt.Errorf("project folder does not exist")
}

func Exists(prjPath string) bool {
	if err := fileop.Exists(prjPath); err == nil {
		return true
	}
	return false
}
