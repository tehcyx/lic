package cmd

import (
	"testing"

	"github.com/tehcyx/lic/pkg/lic/core"
)

func TestNewLicCmd(t *testing.T) {
	opts := core.NewOptions()
	got := NewLicCmd(opts)

	if got == nil {
		t.Fatal("NewLicCmd() returned nil")
	}

	if got.Use != "lic" {
		t.Errorf("NewLicCmd() Use = %v, want 'lic'", got.Use)
	}

	if got.Short == "" {
		t.Error("NewLicCmd() Short description should not be empty")
	}

	if len(got.Commands()) == 0 {
		t.Error("NewLicCmd() should have subcommands")
	}
}
