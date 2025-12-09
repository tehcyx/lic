package license

import "context"

// Provider defines the interface for retrieving license information from different sources
type Provider interface {
	// Supports returns true if this provider can handle the given import path
	Supports(importPath string) bool

	// GetLicense retrieves license information for the given import
	// Returns the license key (SPDX identifier) and any error encountered
	GetLicense(ctx context.Context, importPath, version, branch, url string) (string, error)

	// Name returns the name of this provider for logging purposes
	Name() string
}
