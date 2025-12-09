package gopath

import (
	"context"
	"fmt"

	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/internal/report"
)

// Collector implements the DependencyCollector interface for GOPATH-based projects
type Collector struct{}

// NewCollector creates a new GOPATH collector
func NewCollector() *Collector {
	return &Collector{}
}

// Name returns the name of this collector
func (c *Collector) Name() string {
	return "GOPATH"
}

// CanHandle returns true if the path exists (GOPATH is the fallback collector)
func (c *Collector) CanHandle(prjPath string) bool {
	return Exists(prjPath)
}

// Collect initiates collection of imports across given path
func (c *Collector) Collect(ctx context.Context, proj *report.Project, prjPath string) error {
	// Check for cancellation before starting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if !Exists(prjPath) {
		return fmt.Errorf("project folder does not exist")
	}
	return ReadImports(proj, prjPath)
}

// Collect is the legacy function for backwards compatibility
// Deprecated: Use Collector.Collect instead
func Collect(proj *report.Project, prjPath string) error {
	c := NewCollector()
	return c.Collect(context.Background(), proj, prjPath)
}

// Exists checks that path exists
func Exists(prjPath string) bool {
	if err := fileop.Exists(prjPath); err == nil {
		return true
	}
	return false
}
