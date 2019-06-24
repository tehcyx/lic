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
	}{ // This function just initialites maps
		{"Create new report success", &Project{
			Imports:           map[string]*Import{},
			ValidatedLicenses: map[string]*Import{},
			Violations:        map[string]*Import{},
		}},
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
		{"Name cannot be empty", fields{"Test project", nil, "0.0.1", "master", "rev"}, args{"", "", "", "", true}, true},
		{"Name cannot be empty", fields{"Test project", map[string]*Import{"name": NewImport("name", "", "branch", "rev", true)}, "0.0.1", "master", "rev"}, args{"name", "", "", "", true}, true},
		{"Name cannot be empty", fields{"Test project", map[string]*Import{"name": NewImport("name", "", "branch", "rev", true)}, "0.0.1", "master", "rev"}, args{"newName", "", "", "", true}, false},
		// import already exists [name]
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
	}{ // sets fields to whatever is given, can be empty, no checks applied
		{"Set all fields", args{"name", "version", "branch", "revision", true}, &Import{Name: "name", Version: "version", Branch: "branch", Revision: "revision", IsDirectDependency: true}},
		{"Set no fields", args{"", "", "", "", true}, &Import{IsDirectDependency: true}},
		{"Set no fields - direct dependency false", args{"", "", "", "", false}, &Import{IsDirectDependency: false}},
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
	}{ // output some basic in stdout, just don't break the code
		{"Print output successful", fields{"1", "name", "hash", "version", "branch", "revision", license.Licenses["na"], nil, nil, nil}},
		{"Print output successful", fields{"1", "name", "hash", "version", "branch", "revision", license.Licenses["na"], map[string]*Import{"name": NewImport("name", "version", "branch", "revision", true)}, nil, nil}},
		{"Print output successful", fields{"1", "name", "hash", "version", "branch", "revision", license.Licenses["na"], nil, map[string]*Import{"name": NewImport("name", "version", "branch", "revision", true)}, nil}},
		{"Print output successful", fields{"1", "name", "hash", "version", "branch", "revision", license.Licenses["na"], nil, nil, map[string]*Import{"name": NewImport("name", "version", "branch", "revision", true)}}},
		{"Print output successful", fields{"1", "name", "hash", "version", "branch", "revision", license.Licenses["na"], map[string]*Import{"name": NewImport("name", "version", "branch", "revision", true)}, map[string]*Import{"name": NewImport("name", "version", "branch", "revision", true)}, map[string]*Import{"name": NewImport("name", "version", "branch", "revision", true)}}},
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

func TestImport_GetLicenseInfo(t *testing.T) {
	type fields struct {
		Name               string
		Hash               string
		Version            string
		Branch             string
		Revision           string
		ParsedURL          string
		IsDirectDependency bool
		License            license.License
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"Get license info existing project", fields{"github.com/tehcyx/girc", "", "", "", "", "http://github.com/tehcyx/girc", true, license.Licenses["na"]}},
		{"Get license info nonexistent project don't break it", fields{"nonexistent", "", "", "", "", "http://github.com/tehcyx/nonexistent", true, license.Licenses["na"]}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Import{
				Name:               tt.fields.Name,
				Hash:               tt.fields.Hash,
				Version:            tt.fields.Version,
				Branch:             tt.fields.Branch,
				Revision:           tt.fields.Revision,
				ParsedURL:          tt.fields.ParsedURL,
				IsDirectDependency: tt.fields.IsDirectDependency,
				License:            tt.fields.License,
			}
			i.GetLicenseInfo()
		})
	}
}
