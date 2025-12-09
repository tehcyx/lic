// Package report implements the `lic report` command.
package report

import (
	"testing"

	"github.com/tehcyx/lic/pkg/lic/core"
)

func TestNewReportOptions(t *testing.T) {
	opts := core.NewOptions()
	got := NewReportOptions(opts)

	if got == nil {
		t.Fatal("NewReportOptions() returned nil")
	}

	if got.Options != opts {
		t.Error("NewReportOptions() should wrap the provided Options")
	}

	// Verify default values for report-specific fields
	if got.Upload {
		t.Error("NewReportOptions() Upload should default to false")
	}

	if got.HTMLOutput {
		t.Error("NewReportOptions() HTMLOutput should default to false")
	}

	if got.StdLib {
		t.Error("NewReportOptions() StdLib should default to false")
	}
}

func TestNewReportCmd(t *testing.T) {
	got := NewReportCmd()

	if got == nil {
		t.Fatal("NewReportCmd() returned nil")
	}

	if got.Use != "report" {
		t.Errorf("NewReportCmd() Use = %v, want 'report'", got.Use)
	}

	if got.Short == "" {
		t.Error("NewReportCmd() Short description should not be empty")
	}

	// Verify aliases are set
	if len(got.Aliases) == 0 {
		t.Error("NewReportCmd() should have at least one alias")
	}

	expectedAlias := "r"
	found := false
	for _, alias := range got.Aliases {
		if alias == expectedAlias {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("NewReportCmd() should have alias %q", expectedAlias)
	}
}
