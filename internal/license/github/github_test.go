package github

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-github/v25/github"
)

// mockGitHubClient implements GitHubClient interface for testing
type mockGitHubClient struct {
	repos map[string]*github.Repository
	err   error
}

func (m *mockGitHubClient) GetRepository(ctx context.Context, owner, repo string) (*github.Repository, *github.Response, error) {
	if m.err != nil {
		return nil, &github.Response{Response: &http.Response{StatusCode: 404}}, m.err
	}

	key := owner + "/" + repo
	if r, ok := m.repos[key]; ok {
		return r, &github.Response{Response: &http.Response{StatusCode: 200}}, nil
	}

	return nil, &github.Response{Response: &http.Response{StatusCode: 404}}, fmt.Errorf("repository not found")
}

func newMockRepo(licenseKey string) *github.Repository {
	key := licenseKey
	return &github.Repository{
		License: &github.License{
			Key: &key,
		},
	}
}

func TestProvider_GetLicenseKey(t *testing.T) {
	tests := []struct {
		name       string
		importPath string
		mockRepos  map[string]*github.Repository
		mockErr    error
		want       string
		wantErr    bool
	}{
		{
			name:       "success - apache license",
			importPath: "github.com/example/repo",
			mockRepos: map[string]*github.Repository{
				"example/repo": newMockRepo("apache-2.0"),
			},
			want:    "apache-2.0",
			wantErr: false,
		},
		{
			name:       "success - mit license",
			importPath: "github.com/owner/project",
			mockRepos: map[string]*github.Repository{
				"owner/project": newMockRepo("mit"),
			},
			want:    "mit",
			wantErr: false,
		},
		{
			name:       "repository not found",
			importPath: "github.com/nonexistent/repo",
			mockRepos:  map[string]*github.Repository{},
			want:       "",
			wantErr:    true,
		},
		{
			name:       "repository without license",
			importPath: "github.com/example/nolicense",
			mockRepos: map[string]*github.Repository{
				"example/nolicense": &github.Repository{License: nil},
			},
			want:    "",
			wantErr: true,
		},
		{
			name:       "invalid import path",
			importPath: "github.com/invalid",
			mockRepos:  map[string]*github.Repository{},
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockGitHubClient{
				repos: tt.mockRepos,
				err:   tt.mockErr,
			}
			provider := NewProviderWithClient(mockClient)

			ctx := context.Background()
			got, err := provider.GetLicenseKey(ctx, tt.importPath)

			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.GetLicenseKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Provider.GetLicenseKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_GetLicense(t *testing.T) {
	mockClient := &mockGitHubClient{
		repos: map[string]*github.Repository{
			"example/repo": newMockRepo("mit"),
		},
	}
	provider := NewProviderWithClient(mockClient)

	ctx := context.Background()
	got, err := provider.GetLicense(ctx, "github.com/example/repo", "", "", "")

	if err != nil {
		t.Errorf("Provider.GetLicense() unexpected error = %v", err)
		return
	}
	if got != "mit" {
		t.Errorf("Provider.GetLicense() = %v, want 'mit'", got)
	}
}

func TestProvider_Supports(t *testing.T) {
	provider := NewProvider()

	tests := []struct {
		name       string
		importPath string
		want       bool
	}{
		{"github repository", "github.com/user/repo", true},
		{"github with path", "github.com/user/repo/pkg", true},
		{"gitlab repository", "gitlab.com/user/repo", false},
		{"gopkg.in", "gopkg.in/yaml.v2", false},
		{"golang.org", "golang.org/x/oauth2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := provider.Supports(tt.importPath); got != tt.want {
				t.Errorf("Provider.Supports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProvider_Name(t *testing.T) {
	provider := NewProvider()
	if got := provider.Name(); got != "GitHub" {
		t.Errorf("Provider.Name() = %v, want 'GitHub'", got)
	}
}

func Test_parseRepoOwner(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantRepo  string
		wantOwner string
		wantErr   bool
	}{
		{name: "Repo owner from repo - success", args: args{name: "github.com/repo/owner"}, wantRepo: "repo", wantOwner: "owner", wantErr: false},
		{name: "Repo owner from repo - error no owner", args: args{name: "github.com/repo"}, wantRepo: "", wantOwner: "", wantErr: true},
		{name: "Repo owner from repo - error invalid", args: args{name: "test/test1/test2"}, wantRepo: "", wantOwner: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseRepoOwner(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRepoOwner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantRepo {
				t.Errorf("parseRepoOwner() got = %v, want %v", got, tt.wantRepo)
			}
			if got1 != tt.wantOwner {
				t.Errorf("parseRepoOwner() got1 = %v, want %v", got1, tt.wantOwner)
			}
		})
	}
}

func Test_calculateBackoff(t *testing.T) {
	tests := []struct {
		name          string
		attempt       int
		wantMinDelay  int // in seconds
		wantMaxDelay  int // in seconds
	}{
		{
			name:         "first retry - 1 second",
			attempt:      0,
			wantMinDelay: 1,
			wantMaxDelay: 1,
		},
		{
			name:         "second retry - 2 seconds",
			attempt:      1,
			wantMinDelay: 2,
			wantMaxDelay: 2,
		},
		{
			name:         "third retry - 4 seconds",
			attempt:      2,
			wantMinDelay: 4,
			wantMaxDelay: 4,
		},
		{
			name:         "fourth retry - 8 seconds",
			attempt:      3,
			wantMinDelay: 8,
			wantMaxDelay: 8,
		},
		{
			name:         "large attempt - capped at 30 seconds",
			attempt:      10,
			wantMinDelay: 30,
			wantMaxDelay: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateBackoff(tt.attempt)
			gotSeconds := int(got.Seconds())

			if gotSeconds < tt.wantMinDelay || gotSeconds > tt.wantMaxDelay {
				t.Errorf("calculateBackoff(%d) = %v seconds, want between %d and %d seconds",
					tt.attempt, gotSeconds, tt.wantMinDelay, tt.wantMaxDelay)
			}
		})
	}
}

func Test_isRetriableError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		resp *github.Response
		want bool
	}{
		{
			name: "nil error should not retry",
			err:  nil,
			resp: nil,
			want: false,
		},
		{
			name: "context deadline exceeded should retry",
			err:  context.DeadlineExceeded,
			resp: nil,
			want: true,
		},
		{
			name: "500 status should retry",
			err:  fmt.Errorf("server error"),
			resp: &github.Response{Response: &http.Response{StatusCode: 500}},
			want: true,
		},
		{
			name: "502 status should retry",
			err:  fmt.Errorf("bad gateway"),
			resp: &github.Response{Response: &http.Response{StatusCode: 502}},
			want: true,
		},
		{
			name: "503 status should retry",
			err:  fmt.Errorf("service unavailable"),
			resp: &github.Response{Response: &http.Response{StatusCode: 503}},
			want: true,
		},
		{
			name: "404 status should not retry",
			err:  fmt.Errorf("not found"),
			resp: &github.Response{Response: &http.Response{StatusCode: 404}},
			want: false,
		},
		{
			name: "400 status should not retry",
			err:  fmt.Errorf("bad request"),
			resp: &github.Response{Response: &http.Response{StatusCode: 400}},
			want: false,
		},
		{
			name: "nil response should not retry",
			err:  fmt.Errorf("some error"),
			resp: nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetriableError(tt.err, tt.resp)
			if got != tt.want {
				t.Errorf("isRetriableError() = %v, want %v", got, tt.want)
			}
		})
	}
}
