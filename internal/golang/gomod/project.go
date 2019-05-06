package gomod

import (
	"fmt"
	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/internal/licensereport"
	"path"
)

func Collect(proj *licensereport.Project, prjPath string) error {
	goModPath := path.Join(prjPath, "go.mod")
	if Exists(goModPath) {
		return ReadImports(proj, goModPath)
	}
	return fmt.Errorf("go.mod does not exist")
}

func Exists(goModPath string) bool {
	if err := fileop.Exists(goModPath); err == nil {
		return true
	}
	return false
}
