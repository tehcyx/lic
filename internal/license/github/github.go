package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/google/go-github/v25/github"
	"golang.org/x/oauth2"
)

// TODO: general idea, to keep track of GitHub API rate limits, since they're passed in every HTTP request to interrupt querying if rate limits are reached.
// X-RateLimit-Limit: 60
// X-RateLimit-Remaining: 56
// X-RateLimit-Reset: 1372700873

var tokenVar string

func init() {
	tokenVar = os.Getenv("LIC_GITHUB_ACCESS_TOKEN")

	// https://github.com/mitchellh/go-spdx
	// https://github.com/mitchellh/golicense/
}

func GetLicenseKey(name string) (string, error) {
	ctx := context.Background()
	var tc *http.Client
	if tokenVar != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: tokenVar},
		)
		tc = oauth2.NewClient(ctx, ts)
	}
	client := github.NewClient(tc)

	repository, owner, err := parseRepoOwner(name)
	if err != nil {
		return "", err
	}

	repo, resp, err := client.Repositories.Get(ctx, repository, owner) // TODO: cleanup the error catching, this is horrible
	if _, ok := err.(*github.RateLimitError); ok {
		return "", fmt.Errorf("API hit rate limit")
	}
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub API returned an error")
	}
	if repo.License != nil {
		return *repo.License.Key, nil
	}
	return "", nil
}

func parseRepoOwner(name string) (string, string, error) {
	githubParser := regexp.MustCompile(`github.com/(?P<owner>[^/]*)/(?P<repo>[^/]*)`)

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
