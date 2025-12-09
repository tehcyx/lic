package gopath

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tehcyx/lic/internal/report"
)

func TestReadImports(t *testing.T) {
	proj := report.NewProjectReport()

	// Create a temporary directory with a simple Go file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	goCode := `package main

import (
	"fmt"
	"os"
	"github.com/example/repo"
)

func main() {
	fmt.Println("test")
}
`
	if err := os.WriteFile(testFile, []byte(goCode), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test reading imports from the directory
	err := ReadImports(proj, tmpDir)
	if err != nil {
		t.Errorf("ReadImports() unexpected error = %v", err)
	}

	// Verify that imports were found
	if len(proj.Imports) == 0 {
		t.Error("ReadImports() should have found imports in the test file")
	}

	// Check that specific imports are present
	foundFmt := false
	foundOs := false
	foundExternal := false

	for _, imp := range proj.Imports {
		if imp.Name == "fmt" {
			foundFmt = true
		}
		if imp.Name == "os" {
			foundOs = true
		}
		if imp.Name == "github.com/example/repo" {
			foundExternal = true
		}
	}

	if !foundFmt {
		t.Error("ReadImports() should have found 'fmt' import")
	}
	if !foundOs {
		t.Error("ReadImports() should have found 'os' import")
	}
	if !foundExternal {
		t.Error("ReadImports() should have found 'github.com/example/repo' import")
	}
}
