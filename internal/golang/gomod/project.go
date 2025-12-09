package gomod

import (
	"context"
	"fmt"
	"path"

	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/internal/report"
)

// Collector implements the DependencyCollector interface for go.mod files
type Collector struct{}

// NewCollector creates a new go.mod collector
func NewCollector() *Collector {
	return &Collector{}
}

// Name returns the name of this collector
func (c *Collector) Name() string {
	return "go.mod"
}

// CanHandle returns true if a go.mod file exists in the given path
func (c *Collector) CanHandle(prjPath string) bool {
	goModPath := path.Join(prjPath, "go.mod")
	return Exists(goModPath)
}

// Collect initiates collection of imports across given path
func (c *Collector) Collect(ctx context.Context, proj *report.Project, prjPath string) error {
	// Check for cancellation before starting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	goModPath := path.Join(prjPath, "go.mod")
	if !Exists(goModPath) {
		return fmt.Errorf("go.mod does not exist")
	}
	return ReadImports(proj, goModPath)
}

// Collect is the legacy function for backwards compatibility
// Deprecated: Use Collector.Collect instead
func Collect(proj *report.Project, prjPath string) error {
	c := NewCollector()
	return c.Collect(context.Background(), proj, prjPath)
}

// Exists checks that go.mod file exists on given path
func Exists(goModPath string) bool {
	if err := fileop.Exists(goModPath); err == nil {
		return true
	}
	return false
}
