package github

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v25/github"
	"golang.org/x/oauth2"
)

const (
	// maxRetries is the maximum number of retry attempts for API calls
	maxRetries = 3
	// baseDelay is the initial delay between retries (exponential backoff)
	baseDelay = 1 * time.Second
	// requestTimeout is the timeout for each API request
	requestTimeout = 10 * time.Second
)

var (
	tokenVar     string
	githubParser = regexp.MustCompile(`github\.com/(?P<owner>[^/]*)/(?P<repo>[^/]*)`)
)

func init() {
	tokenVar = os.Getenv("LIC_GITHUB_ACCESS_TOKEN")

	// https://github.com/mitchellh/go-spdx
	// https://github.com/mitchellh/golicense/
}

// GitHubClient is an interface for GitHub API operations (allows mocking)
type GitHubClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error)
}

// realGitHubClient wraps the actual GitHub client
type realGitHubClient struct {
	client *github.Client
}

func (r *realGitHubClient) GetRepository(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	return r.client.Repositories.Get(ctx, owner, repo)
}

// Provider implements the LicenseProvider interface for GitHub repositories
type Provider struct {
	client GitHubClient
}

// NewProvider creates a new GitHub license provider
func NewProvider() *Provider {
	return &Provider{
		client: nil, // Will be created lazily
	}
}

// NewProviderWithClient creates a provider with a custom client (for testing)
func NewProviderWithClient(client GitHubClient) *Provider {
	return &Provider{
		client: client,
	}
}

// Name returns the name of this provider
func (p *Provider) Name() string {
	return "GitHub"
}

// Supports returns true if the import path is from github.com
func (p *Provider) Supports(importPath string) bool {
	return strings.HasPrefix(importPath, "github.com/")
}

// GetLicense retrieves license information from GitHub API
func (p *Provider) GetLicense(ctx context.Context, importPath, version, branch, url string) (string, error) {
	return p.GetLicenseKey(ctx, importPath)
}

// GetLicenseKey retrieves the license key for a GitHub repository using the provider's client
func (p *Provider) GetLicenseKey(ctx context.Context, name string) (string, error) {
	owner, repository, err := parseRepoOwner(name)
	if err != nil {
		return "", err
	}

	// Create client if not set (lazy initialization for non-test cases)
	if p.client == nil {
		var tc *http.Client
		if tokenVar != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: tokenVar},
			)
			tc = oauth2.NewClient(ctx, ts)
		}
		p.client = &realGitHubClient{client: github.NewClient(tc)}
	}

	// Retry loop with exponential backoff
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Create context with timeout for this attempt
		ctxWithTimeout, cancel := context.WithTimeout(ctx, requestTimeout)

		repo, resp, err := p.client.GetRepository(ctxWithTimeout, owner, repository)
		cancel() // Always cancel context to free resources

		// Handle rate limit errors
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			if attempt < maxRetries {
				// Calculate wait time until rate limit resets
				waitDuration := time.Until(rateLimitErr.Rate.Reset.Time)
				if waitDuration > 0 && waitDuration < 5*time.Minute {
					log.Printf("Rate limit hit for %s/%s. Waiting %v until reset...\n", owner, repository, waitDuration.Round(time.Second))
					time.Sleep(waitDuration)
					continue
				}
			}
			return "", fmt.Errorf("API rate limit exceeded for %s/%s (resets at %v)", owner, repository, rateLimitErr.Rate.Reset.Time)
		}

		// Handle transient errors that should be retried
		if err != nil {
			lastErr = err
			if attempt < maxRetries && isRetriableError(err, resp) {
				delay := calculateBackoff(attempt)
				log.Printf("Retrying request for %s/%s after error (attempt %d/%d, waiting %v): %v\n",
					owner, repository, attempt+1, maxRetries+1, delay, err)
				time.Sleep(delay)
				continue
			}
			return "", fmt.Errorf("failed to get repository %s/%s after %d attempts: %w", owner, repository, attempt+1, err)
		}

		// Check HTTP status code
		if resp.StatusCode != 200 {
			lastErr = fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
			if attempt < maxRetries && resp.StatusCode >= 500 {
				// Retry on 5xx server errors
				delay := calculateBackoff(attempt)
				log.Printf("Retrying request for %s/%s after %d status (attempt %d/%d, waiting %v)\n",
					owner, repository, resp.StatusCode, attempt+1, maxRetries+1, delay)
				time.Sleep(delay)
				continue
			}
			return "", fmt.Errorf("GitHub API returned status %d for %s/%s", resp.StatusCode, owner, repository)
		}

		// Success - extract license
		if repo.License != nil && repo.License.Key != nil {
			return *repo.License.Key, nil
		}
		return "", fmt.Errorf("no license found for %s/%s", owner, repository)
	}

	// All retries exhausted
	return "", fmt.Errorf("failed to get repository %s/%s after %d attempts: %w", owner, repository, maxRetries+1, lastErr)
}

// GetLicenseKey is the legacy function for backward compatibility
// Deprecated: Use Provider.GetLicenseKey instead
func GetLicenseKey(ctx context.Context, name string) (string, error) {
	p := NewProvider()
	return p.GetLicenseKey(ctx, name)
}

// calculateBackoff computes exponential backoff delay for retry attempt
func calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	// attempt 0: 1s, attempt 1: 2s, attempt 2: 4s
	multiplier := math.Pow(2, float64(attempt))
	delay := time.Duration(float64(baseDelay) * multiplier)

	// Cap at 30 seconds
	maxDelay := 30 * time.Second
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}

// isRetriableError determines if an error should trigger a retry
func isRetriableError(err error, resp *github.Response) bool {
	if err == nil {
		return false
	}

	// Retry on timeout or temporary network errors
	if err == context.DeadlineExceeded {
		return true
	}

	// Retry on 5xx server errors (handled separately in main function)
	if resp != nil && resp.StatusCode >= 500 && resp.StatusCode < 600 {
		return true
	}

	return false
}

func parseRepoOwner(name string) (string, string, error) {
	match := githubParser.FindStringSubmatch(name)
	matchResult := make(map[string]string)
	if len(match) > 1 {
		for i, name := range githubParser.SubexpNames() {
			if i != 0 && name != "" {
				matchResult[name] = match[i]
			}
		}
		owner, okOwner := matchResult["owner"]
		repo, okRepo := matchResult["repo"]
		if okOwner && okRepo {
			return owner, repo, nil
		}
	}
	return "", "", fmt.Errorf("Couldn't figure out repository information")

}
