package golang

import (
	"context"

	"github.com/tehcyx/lic/internal/report"
)

// DependencyCollector defines the interface for collecting dependencies from different sources
type DependencyCollector interface {
	// CanHandle returns true if this collector can handle the given project path
	CanHandle(path string) bool

	// Collect collects dependencies from the project path and populates the project report
	Collect(ctx context.Context, proj *report.Project, path string) error

	// Name returns the name of this collector for logging purposes
	Name() string
}
