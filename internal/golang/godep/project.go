package godep

import (
	"github.com/tehcyx/lic/internal/licensereport"
)

func Collect(proj *licensereport.Project, prjPath string) error {
	return ReadImports(proj, prjPath)
}
