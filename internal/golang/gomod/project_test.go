package gomod

import (
	"testing"

	"github.com/tehcyx/lic/internal/report"
)

var (
	modPackageName = `module github.com/tehcyx/imaginary-api

require (
	github.com/tehcyx/imaginary-service
)`
)

func TestCollect(t *testing.T) {
	type args struct {
		proj    *report.Project
		prjPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Collect(tt.args.proj, tt.args.prjPath); (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		goModPath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.args.goModPath); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
