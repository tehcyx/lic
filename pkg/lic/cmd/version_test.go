package cmd

import (
	"testing"

	"github.com/tehcyx/lic/pkg/lic/core"
)

func TestNewVersionOptions(t *testing.T) {
	opts := core.NewOptions()
	got := NewVersionOptions(opts)

	if got == nil {
		t.Fatal("NewVersionOptions() returned nil")
	}

	if got.Options != opts {
		t.Error("NewVersionOptions() should wrap the provided Options")
	}
}

func TestNewVersionCmd(t *testing.T) {
	opts := NewVersionOptions(core.NewOptions())
	got := NewVersionCmd(opts)

	if got == nil {
		t.Fatal("NewVersionCmd() returned nil")
	}

	if got.Use != "version" {
		t.Errorf("NewVersionCmd() Use = %v, want 'version'", got.Use)
	}

	if got.Short == "" {
		t.Error("NewVersionCmd() Short description should not be empty")
	}

	if got.Long == "" {
		t.Error("NewVersionCmd() Long description should not be empty")
	}
}

func TestVersionOptions_Run(t *testing.T) {
	opts := NewVersionOptions(core.NewOptions())

	// Run() should never return an error - it just prints version info
	err := opts.Run()
	if err != nil {
		t.Errorf("VersionOptions.Run() unexpected error = %v", err)
	}
}
