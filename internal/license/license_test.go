package license

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestInit(t *testing.T) {
	if len(Licenses) == 0 {
		t.Errorf("initialization of licenses failed")
	}

	if _, ok := Licenses["na"]; !ok {
		t.Errorf("Fallback not initialized, something went wrong")
	}

	// Verify we have a comprehensive SPDX license list (100+ licenses)
	if len(Licenses) < 100 {
		t.Errorf("Expected at least 100 licenses, got %d", len(Licenses))
	}
}

func TestLicenseMap_CommonLicenses(t *testing.T) {
	// Test that commonly used licenses are present
	commonLicenses := []string{
		// Permissive
		"mit", "apache-2.0", "bsd-2-clause", "bsd-3-clause", "isc",
		// Copyleft
		"gpl-2.0", "gpl-3.0", "lgpl-2.1", "lgpl-3.0", "agpl-3.0",
		"mpl-2.0", "epl-1.0", "epl-2.0",
		// Creative Commons
		"cc0-1.0", "cc-by-4.0", "cc-by-sa-4.0",
		// Other
		"unlicense", "other", "na",
	}

	for _, key := range commonLicenses {
		if _, ok := Licenses[key]; !ok {
			t.Errorf("Common license %q not found in license map", key)
		}
	}
}

func TestLicenseMap_SPDXIdentifiers(t *testing.T) {
	// Verify all license keys are valid SPDX-style identifiers (lowercase with hyphens)
	invalidKeys := []string{}
	for key := range Licenses {
		// SPDX identifiers should be lowercase and use hyphens
		if strings.ToLower(key) != key {
			invalidKeys = append(invalidKeys, key)
		}
	}

	if len(invalidKeys) > 0 {
		t.Errorf("Found non-SPDX-compliant license keys: %v", invalidKeys)
	}
}

func TestLicenseMap_AllEntriesValid(t *testing.T) {
	// Verify all license entries have required fields
	for key, lic := range Licenses {
		if lic.Name == "" {
			t.Errorf("License %q has empty Name field", key)
		}
		if lic.ShortName == "" && key != "na" {
			t.Errorf("License %q has empty ShortName field", key)
		}
		// ShortName should match the map key
		if lic.ShortName != "" && lic.ShortName != key {
			t.Errorf("License %q has mismatched ShortName: %q", key, lic.ShortName)
		}
	}
}

// mockProvider is a test provider for testing the provider interface
type mockProvider struct {
	prefix      string
	licenseKey  string
	shouldError bool
}

func (m *mockProvider) Supports(importPath string) bool {
	return strings.HasPrefix(importPath, m.prefix)
}

func (m *mockProvider) GetLicense(ctx context.Context, importPath, version, branch, url string) (string, error) {
	if m.shouldError {
		return "", fmt.Errorf("mock error")
	}
	return m.licenseKey, nil
}

func (m *mockProvider) Name() string {
	return "Mock"
}

func TestProvider_Interface(t *testing.T) {
	// Test that mockProvider implements Provider interface
	var _ Provider = &mockProvider{}

	tests := []struct {
		name        string
		provider    *mockProvider
		importPath  string
		wantSupport bool
		wantKey     string
		wantErr     bool
	}{
		{
			name:        "supports matching prefix",
			provider:    &mockProvider{prefix: "example.com/", licenseKey: "mit"},
			importPath:  "example.com/user/repo",
			wantSupport: true,
			wantKey:     "mit",
			wantErr:     false,
		},
		{
			name:        "does not support non-matching prefix",
			provider:    &mockProvider{prefix: "example.com/", licenseKey: "mit"},
			importPath:  "other.com/user/repo",
			wantSupport: false,
		},
		{
			name:        "returns error when configured",
			provider:    &mockProvider{prefix: "example.com/", shouldError: true},
			importPath:  "example.com/user/repo",
			wantSupport: true,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.provider.Supports(tt.importPath); got != tt.wantSupport {
				t.Errorf("Supports() = %v, want %v", got, tt.wantSupport)
			}

			if tt.wantSupport {
				ctx := context.Background()
				key, err := tt.provider.GetLicense(ctx, tt.importPath, "", "", "")
				if (err != nil) != tt.wantErr {
					t.Errorf("GetLicense() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !tt.wantErr && key != tt.wantKey {
					t.Errorf("GetLicense() key = %v, want %v", key, tt.wantKey)
				}
			}
		})
	}
}

func TestGetProviders(t *testing.T) {
	providers := getProviders()

	if len(providers) == 0 {
		t.Error("getProviders() returned empty list")
	}

	// Check that GitHub provider is included
	foundGitHub := false
	for _, p := range providers {
		if p.Name() == "GitHub" {
			foundGitHub = true
			// Test that it supports github.com imports
			if !p.Supports("github.com/user/repo") {
				t.Error("GitHub provider should support github.com imports")
			}
			// Test that it doesn't support non-github imports
			if p.Supports("example.com/user/repo") {
				t.Error("GitHub provider should not support non-github imports")
			}
		}
	}

	if !foundGitHub {
		t.Error("getProviders() should include GitHub provider")
	}
}

func TestGet(t *testing.T) {
	type args struct {
		name    string
		version string
		branch  string
		url     string
	}
	tests := []struct {
		name string
		args args
		want License
	}{
		{name: "Get success", args: args{name: "github.com/tehcyx/girc", version: "0.0.1", branch: "", url: "http://github.com/tehcyx/girc"}, want: Licenses["apache-2.0"]},
		{name: "Get no license", args: args{name: "noproject", version: "", branch: "", url: "noprojecturl"}, want: Licenses["na"]},
		{name: "Get no license", args: args{name: "github.com/nouser/noproject", version: "", branch: "", url: "noprojecturl"}, want: Licenses["na"]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(tt.args.name, tt.args.version, tt.args.branch, tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet_UnsupportedImport(t *testing.T) {
	// Test import path with no supporting provider
	got := Get("unsupported.example.com/user/repo", "v1.0.0", "", "")
	want := Licenses["na"]

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Get() for unsupported import = %v, want %v", got, want)
	}
}
