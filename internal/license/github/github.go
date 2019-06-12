package github

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/v25/github"
)

func init() {
	// TODO: if there's a github secret env, use it to initialize authenticated requests

	// https://github.com/mitchellh/go-spdx
	// https://github.com/mitchellh/golicense/
}

func GetLicenseKey(name string) (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	repository, owner, err := parseRepoOwner(name)
	if err != nil {
		return "", err
	}

	repo, resp, err := client.Repositories.Get(ctx, repository, owner)
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
	return "", "", fmt.Errorf("Couldn't figure out repository information")

}
