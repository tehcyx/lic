package godep

import (
	"github.com/tehcyx/lic/internal/licensereport"
)

// Collect initiates collection of imports accross given path
func Collect(proj *licensereport.Project, prjPath string) error {
	return ReadImports(proj, prjPath)
}
