// Package report implements the `lic report golang` (`lic r go`)command.
package report

import (
	"testing"

	"github.com/tehcyx/lic/pkg/lic/core"
)

func TestRun(t *testing.T) { // TODO refactor Run() method into multiple function calls, to isolate functionality
	opts := NewGolangReportOptions(core.NewOptions())
	opts.ProjectVersion = "0.1.0"
	opts.SrcPath = "testdata/"
	// runErr := opts.Run()
	// if runErr != nil {
	// 	t.Errorf("Run exited with an error")
	// }
}
