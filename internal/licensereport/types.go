package licensereport

import (
	"fmt"
)

// projImport holds version information & name, scanned from various files for an import in that file
type projImport struct {
	Name               string
	Version            string
	Branch             string
	Revision           string
	IsDirectDependency bool
}

// Project holds version information & name, scanned from various files
type Project struct {
	Name     string
	Imports  map[string]*projImport
	Version  string
	Branch   string
	Revision string
}

// NewProjectReport Creates a new project report
func NewProjectReport() *Project {
	return &Project{
		Imports: map[string]*projImport{},
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

func newImport(name, version, branch, revision string, isDirectDependency bool) *projImport {
	return &projImport{
		Name:               name,
		Version:            version,
		Branch:             branch,
		Revision:           revision,
		IsDirectDependency: isDirectDependency,
	}
}
