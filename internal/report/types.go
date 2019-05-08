package report

import (
	"fmt"
)

// projImport holds version information & name, scanned from various files for an import in that file
type projImport struct {
	Name               string
	Hash               string
	Version            string
	Branch             string
	Revision           string
	ParsedURL          string
	IsDirectDependency bool
	License            License
}

// Project holds version information & name, scanned from various files
type Project struct {
	ID                string
	Name              string
	Hash              string
	Version           string
	Branch            string
	Revision          string
	License           License
	Imports           map[string]*projImport
	ValidatedLicenses map[string]*projImport
	Violations        map[string]*projImport
}

// NewProjectReport Creates a new project report
func NewProjectReport() *Project {
	return &Project{
		Imports:           map[string]*projImport{},
		ValidatedLicenses: map[string]*projImport{},
		Violations:        map[string]*projImport{},
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
	p.Imports[name] = newImport(name, version, branch, revision, isDirectDependency)
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
		fmt.Printf("\tImport: %s, Version: %s\n", licen.Name, licen.Version)
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

func newImport(name, version, branch, revision string, isDirectDependency bool) *projImport {
	return &projImport{
		Name:               name,
		Version:            version,
		Branch:             branch,
		Revision:           revision,
		IsDirectDependency: isDirectDependency,
	}
}
