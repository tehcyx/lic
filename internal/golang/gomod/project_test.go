package gomod

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tehcyx/lic/internal/report"
)

var (
	modPackageName = `module github.com/tehcyx/imaginary-api

require (
	github.com/tehcyx/imaginary-service
)`
)

func TestCollect(t *testing.T) {
	proj := report.NewProjectReport()

	// Test with non-existent path
	err := Collect(proj, "/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("Collect() should return error for non-existent go.mod")
	}

	// Test with valid go.mod file
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	goModContent := `module github.com/test/project

go 1.24

require (
	github.com/example/dep1 v1.0.0
	github.com/example/dep2 v2.1.0
)
`
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create test go.mod: %v", err)
	}

	err = Collect(proj, tmpDir)
	if err != nil {
		t.Errorf("Collect() unexpected error = %v", err)
	}

	// Verify imports were collected
	if len(proj.Imports) == 0 {
		t.Error("Collect() should have found dependencies in go.mod")
	}

	// Test with context cancellation
	c := NewCollector()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	proj2 := report.NewProjectReport()
	err = c.Collect(ctx, proj2, tmpDir)
	if err != context.Canceled {
		t.Errorf("Collect() with cancelled context should return context.Canceled, got %v", err)
	}
}

func TestExists(t *testing.T) {
	// Test with non-existent file
	if Exists("/nonexistent/path/go.mod") {
		t.Error("Exists() should return false for non-existent go.mod")
	}

	// Test with existing go.mod
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("Failed to create test go.mod: %v", err)
	}

	if !Exists(goModPath) {
		t.Error("Exists() should return true for existing go.mod")
	}
}

func TestCollector_Name(t *testing.T) {
	c := NewCollector()
	if got := c.Name(); got != "go.mod" {
		t.Errorf("Collector.Name() = %v, want 'go.mod'", got)
	}
}

func TestCollector_CanHandle(t *testing.T) {
	c := NewCollector()

	// Test with directory containing go.mod
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	if !c.CanHandle(tmpDir) {
		t.Error("CanHandle() should return true for directory with go.mod")
	}

	// Test with directory without go.mod
	tmpDir2 := t.TempDir()
	if c.CanHandle(tmpDir2) {
		t.Error("CanHandle() should return false for directory without go.mod")
	}
}
