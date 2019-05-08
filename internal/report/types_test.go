package report

import (
	"reflect"
	"testing"
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
		Imports  map[string]*projImport
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

func Test_newImport(t *testing.T) {
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
		want *projImport
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newImport(tt.args.name, tt.args.version, tt.args.branch, tt.args.revision, tt.args.isDirectDependency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newImport() = %v, want %v", got, tt.want)
			}
		})
	}
}
