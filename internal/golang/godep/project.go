package godep

import (
	"context"
	"path"

	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/internal/report"
)

// Collector implements the DependencyCollector interface for Gopkg.lock files
type Collector struct{}

// NewCollector creates a new godep collector
func NewCollector() *Collector {
	return &Collector{}
}

// Name returns the name of this collector
func (c *Collector) Name() string {
	return "Gopkg.lock"
}

// CanHandle returns true if a Gopkg.lock file exists in the given path
func (c *Collector) CanHandle(prjPath string) bool {
	gopkgPath := path.Join(prjPath, "Gopkg.lock")
	return fileop.Exists(gopkgPath) == nil
}

// Collect initiates collection of imports across given path
func (c *Collector) Collect(ctx context.Context, proj *report.Project, prjPath string) error {
	// Check for cancellation before starting
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return ReadImports(proj, prjPath)
}

// Collect is the legacy function for backwards compatibility
// Deprecated: Use Collector.Collect instead
func Collect(proj *report.Project, prjPath string) error {
	c := NewCollector()
	return c.Collect(context.Background(), proj, prjPath)
}
