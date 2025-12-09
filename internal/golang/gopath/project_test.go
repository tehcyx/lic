package gopath

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tehcyx/lic/internal/report"
)

func TestCollect(t *testing.T) {
	proj := report.NewProjectReport()

	// Test with non-existent path
	err := Collect(proj, "/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("Collect() should return error for non-existent path")
	}

	// Test with valid path (using current directory)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Test legacy function works (it should handle errors from ReadImports)
	err = Collect(proj, wd)
	// We don't assert error here because ReadImports might fail if there are no Go files
	// The important thing is that Collect() doesn't panic

	// Test with context cancellation
	c := NewCollector()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err = c.Collect(ctx, proj, wd)
	if err != context.Canceled {
		t.Errorf("Collect() with cancelled context should return context.Canceled, got %v", err)
	}
}

func TestExists(t *testing.T) {
	// Test with non-existent path
	if Exists("/nonexistent/path/that/does/not/exist") {
		t.Error("Exists() should return false for non-existent path")
	}

	// Test with valid path (temp dir)
	tmpDir := t.TempDir()
	if !Exists(tmpDir) {
		t.Error("Exists() should return true for existing directory")
	}

	// Test with file path
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if !Exists(tmpFile) {
		t.Error("Exists() should return true for existing file")
	}
}

func TestCollector_Name(t *testing.T) {
	c := NewCollector()
	if got := c.Name(); got != "GOPATH" {
		t.Errorf("Collector.Name() = %v, want 'GOPATH'", got)
	}
}

func TestCollector_CanHandle(t *testing.T) {
	c := NewCollector()

	// GOPATH collector always returns true for existing paths (it's the fallback)
	tmpDir := t.TempDir()
	if !c.CanHandle(tmpDir) {
		t.Error("CanHandle() should return true for existing directory")
	}

	// Should return false for non-existent path
	if c.CanHandle("/nonexistent/path/that/does/not/exist") {
		t.Error("CanHandle() should return false for non-existent path")
	}
}
