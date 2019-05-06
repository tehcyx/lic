package licensereport

import (
	"fmt"
)

type projImport struct {
	Name               string
	Version            string
	Branch             string
	Revision           string
	IsDirectDependency bool
}

type Project struct {
	Name     string
	Imports  map[string]*projImport
	Version  string
	Branch   string
	Revision string
}

func NewProjectReport() *Project {
	return &Project{
		Imports: map[string]*projImport{},
	}
}

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
