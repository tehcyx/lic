package golang

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/v25/github"
)

func GetLicenseKey(owner, repository, branch string) (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	repo, resp, err := client.Repositories.Get(ctx, owner, repository)
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

func ParseRepoOwner(name string) (string, string, error) {
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
