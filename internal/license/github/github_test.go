package github

import "testing"

// TODO: use github mock instead, to not use up API rate limits

func TestGetLicenseKey(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "GetLicenseKey success", args: args{name: "github.com/tehcyx/girc"}, want: "apache-2.0", wantErr: false},
		{name: "GetLicenseKey failure", args: args{name: "github.com/tehcyx/doesnotexist"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetLicenseKey(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLicenseKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetLicenseKey() = %v, want %v", got, tt.want)
			}
		})
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
