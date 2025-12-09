package config

import (
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg == nil {
		t.Fatal("Default() returned nil")
	}

	if cfg.Golang.WhitelistDomains == nil {
		t.Error("Default() should initialize WhitelistDomains")
	}

	if cfg.Golang.StdLibPackages == nil {
		t.Error("Default() should initialize StdLibPackages")
	}

	// Check default whitelist domains
	expectedDomains := []string{"github.com", "gopkg.in", "golang.org"}
	if len(cfg.Golang.WhitelistDomains) != len(expectedDomains) {
		t.Errorf("Default() WhitelistDomains length = %d, want %d", len(cfg.Golang.WhitelistDomains), len(expectedDomains))
	}

	for i, domain := range expectedDomains {
		if cfg.Golang.WhitelistDomains[i] != domain {
			t.Errorf("Default() WhitelistDomains[%d] = %s, want %s", i, cfg.Golang.WhitelistDomains[i], domain)
		}
	}

	// Check that stdlib packages are populated
	if len(cfg.Golang.StdLibPackages) == 0 {
		t.Error("Default() should populate StdLibPackages")
	}
}

func TestDefaultStdLibPackages(t *testing.T) {
	packages := DefaultStdLibPackages()

	if len(packages) == 0 {
		t.Error("DefaultStdLibPackages() returned empty map")
	}

	// Test some known standard library packages
	knownPackages := []string{
		"fmt", "io", "os", "net/http", "encoding/json",
		"context", "errors", "strings", "bytes", "time",
		// Go 1.16+ packages
		"embed", "io/fs",
		// Go 1.21+ packages
		"log/slog",
	}

	for _, pkg := range knownPackages {
		if _, ok := packages[pkg]; !ok {
			t.Errorf("DefaultStdLibPackages() missing known package: %s", pkg)
		}
	}

	// Test that deprecated packages are still included for backwards compatibility
	if _, ok := packages["io/ioutil"]; !ok {
		t.Error("DefaultStdLibPackages() should include deprecated package io/ioutil")
	}
}

func TestGolangConfig_IsStdLib(t *testing.T) {
	cfg := &GolangConfig{
		StdLibPackages: DefaultStdLibPackages(),
	}

	tests := []struct {
		name    string
		pkg     string
		wantStd bool
	}{
		{
			name:    "fmt is stdlib",
			pkg:     "fmt",
			wantStd: true,
		},
		{
			name:    "net/http is stdlib",
			pkg:     "net/http",
			wantStd: true,
		},
		{
			name:    "github.com package is not stdlib",
			pkg:     "github.com/spf13/cobra",
			wantStd: false,
		},
		{
			name:    "unknown package is not stdlib",
			pkg:     "example.com/unknown",
			wantStd: false,
		},
		{
			name:    "embed is stdlib (Go 1.16+)",
			pkg:     "embed",
			wantStd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.IsStdLib(tt.pkg)
			if got != tt.wantStd {
				t.Errorf("IsStdLib(%s) = %v, want %v", tt.pkg, got, tt.wantStd)
			}
		})
	}
}

func TestGolangConfig_IsWhitelisted(t *testing.T) {
	cfg := &GolangConfig{
		WhitelistDomains: DefaultWhitelistDomains(),
	}

	tests := []struct {
		name          string
		domain        string
		wantWhitelist bool
	}{
		{
			name:          "github.com is whitelisted",
			domain:        "github.com",
			wantWhitelist: true,
		},
		{
			name:          "gopkg.in is whitelisted",
			domain:        "gopkg.in",
			wantWhitelist: true,
		},
		{
			name:          "golang.org is whitelisted",
			domain:        "golang.org",
			wantWhitelist: true,
		},
		{
			name:          "example.com is not whitelisted",
			domain:        "example.com",
			wantWhitelist: false,
		},
		{
			name:          "gitlab.com is not whitelisted",
			domain:        "gitlab.com",
			wantWhitelist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.IsWhitelisted(tt.domain)
			if got != tt.wantWhitelist {
				t.Errorf("IsWhitelisted(%s) = %v, want %v", tt.domain, got, tt.wantWhitelist)
			}
		})
	}
}

func TestDefaultWhitelistDomains(t *testing.T) {
	domains := DefaultWhitelistDomains()

	if len(domains) == 0 {
		t.Error("DefaultWhitelistDomains() returned empty slice")
	}

	expected := map[string]bool{
		"github.com":  true,
		"gopkg.in":    true,
		"golang.org":  true,
	}

	for _, domain := range domains {
		if !expected[domain] {
			t.Errorf("DefaultWhitelistDomains() contains unexpected domain: %s", domain)
		}
		delete(expected, domain)
	}

	if len(expected) > 0 {
		for domain := range expected {
			t.Errorf("DefaultWhitelistDomains() missing expected domain: %s", domain)
		}
	}
}

func TestConfigImmutability(t *testing.T) {
	// Test that modifying default config doesn't affect new configs
	cfg1 := Default()
	cfg1.Golang.WhitelistDomains = append(cfg1.Golang.WhitelistDomains, "example.com")

	cfg2 := Default()

	// cfg2 should not have the modification from cfg1
	for _, domain := range cfg2.Golang.WhitelistDomains {
		if domain == "example.com" {
			t.Error("Modifying one config affected another - configs are not independent")
		}
	}
}
