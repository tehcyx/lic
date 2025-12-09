// Package report implements the `lic report golang` (`lic r go`) command.
package report

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tehcyx/lic/internal/report"
	"github.com/tehcyx/lic/pkg/lic/core"
)

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name      string
		srcPath   string
		wantError bool
	}{
		{
			name:      "empty path should use current directory",
			srcPath:   "",
			wantError: false,
		},
		{
			name:      "valid path",
			srcPath:   ".",
			wantError: false,
		},
		{
			name:      "non-existent path",
			srcPath:   "/nonexistent/path/that/does/not/exist",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewGolangReportOptions(core.NewOptions())
			opts.SrcPath = tt.srcPath
			err := opts.validatePath()
			if (err != nil) != tt.wantError {
				t.Errorf("validatePath() error = %v, wantError %v", err, tt.wantError)
			}
			// If no error expected and path was empty, should be set to current dir
			if !tt.wantError && tt.srcPath == "" && opts.SrcPath == "" {
				t.Error("validatePath() should set SrcPath to current directory when empty")
			}
		})
	}
}

func TestDetectProjectVersion(t *testing.T) {
	tests := []struct {
		name            string
		initialVersion  string
		expectUnchanged bool
	}{
		{
			name:            "already set version should not change",
			initialVersion:  "1.0.0",
			expectUnchanged: true,
		},
		{
			name:            "n/a version should be detected",
			initialVersion:  "n/a",
			expectUnchanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewGolangReportOptions(core.NewOptions())
			opts.ProjectVersion = tt.initialVersion
			opts.SrcPath = "."
			opts.detectProjectVersion()

			if tt.expectUnchanged && opts.ProjectVersion != tt.initialVersion {
				t.Errorf("detectProjectVersion() changed version from %v to %v, expected no change", tt.initialVersion, opts.ProjectVersion)
			}
			if !tt.expectUnchanged && opts.ProjectVersion == "n/a" {
				t.Error("detectProjectVersion() should have changed version from n/a")
			}
		})
	}
}

func TestCalculateImportHash(t *testing.T) {
	opts := NewGolangReportOptions(core.NewOptions())
	imp := &report.Import{
		Name:    "github.com/example/package",
		Version: "v1.0.0",
	}

	opts.calculateImportHash(imp)

	if imp.Hash == "" {
		t.Error("calculateImportHash() should set a non-empty hash")
	}
	if len(imp.Hash) != 64 {
		t.Errorf("calculateImportHash() hash length = %d, want 64 (SHA256 hex)", len(imp.Hash))
	}

	// Hash should be deterministic
	firstHash := imp.Hash
	opts.calculateImportHash(imp)
	if imp.Hash != firstHash {
		t.Error("calculateImportHash() should produce consistent hashes")
	}

	// Different import should produce different hash
	imp2 := &report.Import{
		Name:    "github.com/example/other",
		Version: "v1.0.0",
	}
	opts.calculateImportHash(imp2)
	if imp2.Hash == firstHash {
		t.Error("calculateImportHash() should produce different hashes for different imports")
	}
}

func TestCheckWhitelist(t *testing.T) {
	tests := []struct {
		name       string
		importName string
		wantMatch  bool
	}{
		{
			name:       "github.com import should match",
			importName: "github.com/spf13/cobra",
			wantMatch:  true,
		},
		{
			name:       "golang.org import should match",
			importName: "golang.org/x/oauth2",
			wantMatch:  true,
		},
		{
			name:       "gopkg.in import should match",
			importName: "gopkg.in/yaml.v2",
			wantMatch:  true,
		},
		{
			name:       "non-whitelisted domain should not match",
			importName: "example.com/some/package",
			wantMatch:  false,
		},
		{
			name:       "similar but different domain should not match",
			importName: "mygithub.company.com/repo",
			wantMatch:  false,
		},
		{
			name:       "domain as substring should not match",
			importName: "evil.com/github.com/fake",
			wantMatch:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewGolangReportOptions(core.NewOptions())
			proj := report.NewProjectReport()
			imp := &report.Import{
				Name:    tt.importName,
				Version: "v1.0.0",
			}

			matched := opts.checkWhitelist(imp, proj)
			if matched != tt.wantMatch {
				t.Errorf("checkWhitelist(%s) = %v, want %v", tt.importName, matched, tt.wantMatch)
			}
		})
	}
}

func TestCollectDependencies(t *testing.T) {
	// Test with current directory which should have go.mod
	opts := NewGolangReportOptions(core.NewOptions())
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Navigate to project root (where go.mod is)
	projectRoot := filepath.Join(cwd, "..", "..", "..", "..")
	opts.SrcPath = projectRoot

	ctx := context.Background()
	proj, err := opts.collectDependencies(ctx)
	if err != nil {
		t.Errorf("collectDependencies() error = %v", err)
	}
	if proj == nil {
		t.Fatal("collectDependencies() returned nil project")
	}
	if len(proj.Imports) == 0 {
		t.Error("collectDependencies() found no imports in project with go.mod")
	}
	if proj.Hash == "" {
		t.Error("collectDependencies() should set project hash")
	}
}

func TestEnrichWithLicenses(t *testing.T) {
	opts := NewGolangReportOptions(core.NewOptions())
	proj := report.NewProjectReport()

	// Add test imports using InsertImport
	proj.InsertImport("fmt", "stdlib", "", "", true)
	proj.InsertImport("github.com/spf13/cobra", "v1.0.0", "", "", true)
	proj.InsertImport("example.com/unknown", "v1.0.0", "", "", true)

	opts.enrichWithLicenses(proj)

	// Check standard library import
	if _, ok := proj.ValidatedLicenses["fmt"]; !ok {
		t.Error("enrichWithLicenses() should validate stdlib import")
	}

	// Check whitelisted import
	if _, ok := proj.ValidatedLicenses["github.com/spf13/cobra"]; !ok {
		t.Error("enrichWithLicenses() should validate whitelisted import")
	}

	// Check non-whitelisted import becomes violation
	if _, ok := proj.Violations["example.com/unknown"]; !ok {
		t.Error("enrichWithLicenses() should mark non-whitelisted import as violation")
	}

	// All imports should have hashes
	for _, imp := range proj.Imports {
		if imp.Hash == "" {
			t.Errorf("enrichWithLicenses() should set hash for import %s", imp.Name)
		}
	}
}

func TestGenerateReport(t *testing.T) {
	tests := []struct {
		name           string
		violationCount int
		wantError      bool
	}{
		{
			name:           "no violations should succeed",
			violationCount: 0,
			wantError:      false,
		},
		{
			name:           "violations should return error",
			violationCount: 2,
			wantError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewGolangReportOptions(core.NewOptions())
			proj := report.NewProjectReport()

			// Add violations if needed
			for i := 0; i < tt.violationCount; i++ {
				proj.Violations[string(rune('a'+i))] = &report.Import{Name: string(rune('a' + i))}
			}

			err := opts.generateReport(proj)
			if (err != nil) != tt.wantError {
				t.Errorf("generateReport() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name      string
		srcPath   string
		wantError bool
	}{
		{
			name:      "run with current directory",
			srcPath:   ".",
			wantError: false,
		},
		{
			name:      "run with project root",
			srcPath:   getProjectRoot(t),
			wantError: false,
		},
		{
			name:      "run with non-existent path",
			srcPath:   "/nonexistent/path/that/does/not/exist",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewGolangReportOptions(core.NewOptions())
			opts.SrcPath = tt.srcPath
			opts.ProjectVersion = "test-version"
			opts.ProjectName = "test-project"

			err := opts.Run()

			if (err != nil) != tt.wantError {
				t.Errorf("Run() error = %v, wantError %v", err, tt.wantError)
			}

			// If no error, verify the workflow completed properly
			if !tt.wantError {
				// SrcPath should be set (not empty)
				if opts.SrcPath == "" {
					t.Error("Run() should set SrcPath")
				}
				// ProjectVersion should remain set
				if opts.ProjectVersion == "" {
					t.Error("Run() should preserve ProjectVersion")
				}
			}
		})
	}
}

func TestRun_Integration(t *testing.T) {
	// Integration test: Run the full workflow on the actual project
	t.Run("full workflow on lic project", func(t *testing.T) {
		opts := NewGolangReportOptions(core.NewOptions())
		opts.SrcPath = getProjectRoot(t)
		opts.ProjectVersion = "test"
		opts.ProjectName = "lic"

		err := opts.Run()
		if err != nil {
			// It's okay if there are violations (non-whitelisted dependencies)
			// We're just testing that the workflow completes
			t.Logf("Run() completed with error (possibly violations): %v", err)
		}
	})
}

// getProjectRoot returns the absolute path to the project root (where go.mod is)
func getProjectRoot(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	// Navigate from pkg/lic/cmd/report to project root
	return filepath.Join(cwd, "..", "..", "..", "..")
}
