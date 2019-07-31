package report

import (
	"fmt"

	"github.com/tehcyx/lic/internal/license"
)

// Import holds version information & name, scanned from various files for an import in that file
type Import struct {
	Name               string
	Hash               string
	Version            string
	Branch             string
	Revision           string
	ParsedURL          string
	IsDirectDependency bool
	License            license.License
}

// Project holds version information & name, scanned from various files
type Project struct {
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

// NewProjectReport Creates a new project report
func NewProjectReport() *Project {
	return &Project{
		Imports:           map[string]*Import{},
		ValidatedLicenses: map[string]*Import{},
		Violations:        map[string]*Import{},
	}
}

// NewImport creates a new project import
func NewImport(name, version, branch, revision string, isDirectDependency bool) *Import {
	// TODO: introduce error handling, name should not be possible to be set to empty, either version, branch or revision have to be set
	return &Import{
		Name:               name,
		Version:            version,
		Branch:             branch,
		Revision:           revision,
		IsDirectDependency: isDirectDependency,
	}
}

// InsertImport creates a new import entry on the project
func (p *Project) InsertImport(name, version, branch, revision string, isDirectDependency bool) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if _, ok := p.Imports[name]; ok {
		return fmt.Errorf("import already exists")
	}
	p.Imports[name] = NewImport(name, version, branch, revision, isDirectDependency)
	return nil
}

// PrintReport outputs the generated report to stdout
func (p *Project) PrintReport() {
	fmt.Printf("Report for %s %s\n", p.Name, p.Version)
	fmt.Printf("Generated project hash: %s\n", p.Hash)
	fmt.Println("")
	numberLicenses := len(p.ValidatedLicenses)
	var wasWere, dependencyDependencies string
	if len(p.ValidatedLicenses) == 1 {
		wasWere = "was"
		dependencyDependencies = "dependency"
	} else {
		wasWere = "were"
		dependencyDependencies = "dependencies"
	}
	fmt.Printf("During the scan there %s %d %s found:\n", wasWere, numberLicenses, dependencyDependencies)

	for _, licen := range p.ValidatedLicenses {
		fmt.Printf("\tImport: %s, Version: %s, License: %s (%s)\n", licen.Name, licen.Version, licen.License.Name, licen.License.ShortName)
	}

	var blacklistImport string
	if len(p.Violations) == 1 {
		wasWere = "was"
		blacklistImport = "blackisted import"
	} else {
		wasWere = "were"
		blacklistImport = "blacklisted imports"
	}
	fmt.Printf("Additionally %d %s %s found:\n", len(p.Violations), blacklistImport, wasWere)
	for _, viol := range p.Violations {
		fmt.Printf("\tImport: %s, Version: %s\n", viol.Name, viol.Version)
	}
}

func (i *Import) GetLicenseInfo() {
	lic := license.Get(i.Name, i.Version, i.Branch, i.ParsedURL)
	i.License = lic
}
