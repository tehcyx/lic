package report

import (
	"reflect"
	"testing"

	"github.com/tehcyx/lic/internal/license"
)

func TestNewProjectReport(t *testing.T) {
	tests := []struct {
		name string
		want *Project
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProjectReport(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProjectReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProject_InsertImport(t *testing.T) {
	type fields struct {
		Name     string
		Imports  map[string]*Import
		Version  string
		Branch   string
		Revision string
	}
	type args struct {
		name               string
		version            string
		branch             string
		revision           string
		isDirectDependency bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				Name:     tt.fields.Name,
				Imports:  tt.fields.Imports,
				Version:  tt.fields.Version,
				Branch:   tt.fields.Branch,
				Revision: tt.fields.Revision,
			}
			if err := p.InsertImport(tt.args.name, tt.args.version, tt.args.branch, tt.args.revision, tt.args.isDirectDependency); (err != nil) != tt.wantErr {
				t.Errorf("Project.InsertImport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewImport(t *testing.T) {
	type args struct {
		name               string
		version            string
		branch             string
		revision           string
		isDirectDependency bool
	}
	tests := []struct {
		name string
		args args
		want *Import
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImport(tt.args.name, tt.args.version, tt.args.branch, tt.args.revision, tt.args.isDirectDependency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProject_PrintReport(t *testing.T) {
	type fields struct {
		ID                string
		Name              string
		Hash              string
		Version           string
		Branch            string
		Revision          string
		License           license.License
		Imports           map[string]*Import
		ValidatedLicenses map[string]*Import
		Violations        map[string]*Import
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				ID:                tt.fields.ID,
				Name:              tt.fields.Name,
				Hash:              tt.fields.Hash,
				Version:           tt.fields.Version,
				Branch:            tt.fields.Branch,
				Revision:          tt.fields.Revision,
				License:           tt.fields.License,
				Imports:           tt.fields.Imports,
				ValidatedLicenses: tt.fields.ValidatedLicenses,
				Violations:        tt.fields.Violations,
			}
			p.PrintReport()
		})
	}
}
