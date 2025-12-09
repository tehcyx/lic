// Package report implements the `lic report golang` (`lic r go`)command.
package report

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/tehcyx/lic/internal/config"
	"github.com/tehcyx/lic/internal/golang"
	"github.com/tehcyx/lic/internal/golang/godep"
	"github.com/tehcyx/lic/internal/golang/gomod"
	"github.com/tehcyx/lic/internal/golang/gopath"
	"github.com/tehcyx/lic/internal/license"
	"github.com/tehcyx/lic/internal/report"

	"github.com/spf13/cobra"

	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/pkg/lic/core"
)

// GolangReportOptions defines available options for the command
type GolangReportOptions struct {
	*Options
	Config *config.Config
}

// Deprecated: Use config.DefaultWhitelistDomains() instead
var DefaultWhitelistResources = config.DefaultWhitelistDomains()

// NewGolangReportOptions creates options with default values
func NewGolangReportOptions(o *core.Options) *GolangReportOptions {
	return &GolangReportOptions{
		Options: NewReportOptions(o),
		Config:  config.Default(),
	}
}

// NewGolangReportCmd creates a new report command
func NewGolangReportCmd(o *GolangReportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "golang",
		Short:   "Generates a report of current working directory or specified path",
		Long:    `Taking in consideration the source on the current path and checking for all licenses, generating a report output in the shell.`,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.Run() },
		Aliases: []string{"go"},
		SilenceUsage: true, // Don't show usage on errors since we already printed the report
	}

	cmd.Flags().BoolVarP(&o.Upload, "upload", "u", false, "Upload report to specified report endpoint to capture continuously")
	cmd.Flags().StringVarP(&o.UploadEndpoint, "upload-endpoint", "", "", "URL of the endpoint to report results of the scans")

	cmd.Flags().StringVarP(&o.SrcPath, "src", "", "", "Local path of sources to scan")
	cmd.Flags().BoolVarP(&o.HTMLOutput, "html-output", "o", false, "Specifies if results should be published as .html-file stored in current path")

	cmd.Flags().StringVarP(&o.ProjectVersion, "project-version", "", "n/a", "Version of scan target")
	cmd.Flags().StringVarP(&o.ProjectName, "project-name", "", "", "Name of scan target")

	cmd.Flags().BoolVarP(&o.StdLib, "stdlib", "s", true, "Should go dependencies be part of the output")

	return cmd
}

// Run runs the command
// Scan has paths this could go:
//  1. If there's a go.mod file, check go.mod file for dependencies and versions of these
//  2. If there's at least one Gopkg.toml/Gopkg.lock file, check Gopkg.lock(s) for all dependencies and versions
//  3. If there's no go.mod file, check $GOPATH and make assumption based on that
func (o *GolangReportOptions) Run() error {
	// Create a context for the entire operation
	ctx := context.Background()

	// Step 1: Validate and set the source path
	if err := o.validatePath(); err != nil {
		return err
	}

	// Step 2: Detect project version from git if needed
	o.detectProjectVersion()

	// Step 3: Collect dependencies using various strategies
	proj, err := o.collectDependencies(ctx)
	if err != nil {
		return err
	}

	// Step 4: Enrich imports with license information
	o.enrichWithLicenses(proj)

	// Step 5: Generate and print the report
	return o.generateReport(proj)
}

// validatePath validates the source path and sets it to current directory if not specified
func (o *GolangReportOptions) validatePath() error {
	if o.SrcPath != "" {
		if err := fileop.Exists(o.SrcPath); err != nil {
			return fmt.Errorf("path '%s' does not exist or you don't have the proper access rights", o.SrcPath)
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("couldn't get current working directory: %w", err)
		}
		o.SrcPath = dir
	}
	return nil
}

// detectProjectVersion attempts to detect the project version from git
func (o *GolangReportOptions) detectProjectVersion() {
	if o.ProjectVersion != "n/a" {
		return // Version already set
	}

	// Check if git is available
	if _, err := exec.LookPath("git"); err == nil {
		cmd := exec.Command("git", "describe", "--tags", "--always")
		cmd.Dir = o.SrcPath
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Printf("Warning: couldn't get git version: %v. Using 'unknown'\n", err)
			o.ProjectVersion = "unknown"
		} else {
			o.ProjectVersion = strings.TrimSpace(out.String())
		}
	} else {
		log.Printf("Warning: git not found in PATH. Using version 'unknown'\n")
		o.ProjectVersion = "unknown"
	}
}

// getCollectors returns the list of dependency collectors in priority order
func (o *GolangReportOptions) getCollectors() []golang.DependencyCollector {
	return []golang.DependencyCollector{
		gomod.NewCollector(),   // Priority 1: go.mod
		godep.NewCollector(),   // Priority 2: Gopkg.lock
		gopath.NewCollector(),  // Priority 3: GOPATH fallback
	}
}

// collectDependencies collects dependencies using the first available collector
func (o *GolangReportOptions) collectDependencies(ctx context.Context) (*report.Project, error) {
	proj := report.NewProjectReport()

	collectors := o.getCollectors()
	var lastErr error

	// Try each collector in order until one succeeds
	for _, collector := range collectors {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if collector.CanHandle(o.SrcPath) {
			log.Printf("Info: Using %s collector", collector.Name())
			err := collector.Collect(ctx, proj, o.SrcPath)
			if err != nil {
				log.Printf("Info: %s collector failed: %v. Trying next collector.", collector.Name(), err)
				lastErr = err
				continue
			}

			// If we successfully collected dependencies, we're done
			if len(proj.Imports) > 0 {
				break
			}

			log.Printf("Info: %s collector found no dependencies. Trying next collector.", collector.Name())
		}
	}

	// Ensure we found at least some dependencies
	if len(proj.Imports) == 0 {
		if lastErr != nil {
			return nil, fmt.Errorf("can't run on source folder: '%s' - no dependencies found (last error: %w)", o.SrcPath, lastErr)
		}
		return nil, fmt.Errorf("can't run on source folder: '%s' - no dependencies found", o.SrcPath)
	}

	// Generate project hash and set version
	h := sha256.New()
	h.Write([]byte(proj.Name + proj.Version))
	proj.Hash = fmt.Sprintf("%x", (h.Sum(nil)))
	proj.Version = o.ProjectVersion

	return proj, nil
}

// enrichWithLicenses enriches each import with license information and validates against whitelist
func (o *GolangReportOptions) enrichWithLicenses(proj *report.Project) {
	for _, imp := range proj.Imports {
		imp.License = license.Licenses["na"]

		// Check if this is a standard library package
		if o.Config.Golang.IsStdLib(imp.Name) {
			imp.Version = "Standard Library"
			proj.ValidatedLicenses[imp.Name] = imp
			o.calculateImportHash(imp)
			continue
		}

		// Check if import matches any whitelisted domain
		isWhitelisted := o.checkWhitelist(imp, proj)
		if !isWhitelisted {
			proj.Violations[imp.Name] = imp
		}

		o.calculateImportHash(imp)
	}
}

// checkWhitelist checks if an import matches the whitelist and fetches license info
func (o *GolangReportOptions) checkWhitelist(imp *report.Import, proj *report.Project) bool {
	for _, whitelistDomain := range o.Config.Golang.WhitelistDomains {
		// Use HasPrefix to ensure the domain is at the start of the import path
		// This prevents matching "mygithub.company.com" when whitelist is "github.com"
		if strings.HasPrefix(imp.Name, whitelistDomain+"/") || imp.Name == whitelistDomain {
			parsedURL, err := url.Parse("https://" + imp.Name)
			if err != nil {
				log.Printf("Warning: invalid URL format for import %s: %v\n", imp.Name, err)
				continue
			}
			imp.ParsedURL = parsedURL.String()
			imp.GetLicenseInfo()

			proj.ValidatedLicenses[imp.Name] = imp
			return true
		}
	}
	return false
}

// calculateImportHash generates a SHA256 hash for an import
func (o *GolangReportOptions) calculateImportHash(imp *report.Import) {
	h := sha256.New()
	h.Write([]byte(imp.Name + imp.Version))
	imp.Hash = fmt.Sprintf("%x", (h.Sum(nil)))
}

// generateReport prints the report and returns an error if violations are found
func (o *GolangReportOptions) generateReport(proj *report.Project) error {
	proj.PrintReport()
	if len(proj.Violations) > 0 {
		return fmt.Errorf("license violations found: %d packages not in whitelist", len(proj.Violations))
	}
	return nil
}
