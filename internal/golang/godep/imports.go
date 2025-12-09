package godep

import (
	"fmt"
	"log"

	"github.com/pelletier/go-toml"
	"github.com/tehcyx/lic/internal/fileop"
	"github.com/tehcyx/lic/internal/report"
)

//ReadImports reads imports on a given filepath with the given regex params for start, end and line
func ReadImports(proj *report.Project, filePath string) error {
	tomlFiles, err := fileop.FilesInPath(filePath, "(?i)/Gopkg.lock$")
	if err != nil {
		return fmt.Errorf("couldn't read files in path %s: %w", filePath, err)
	}
	for _, f := range tomlFiles {
		projToml, err := toml.LoadFile(f)
		if err != nil {
			log.Println(err)
		}
		res := gopkgLock{}
		projToml.Unmarshal(&res)
		for _, prj := range res.Projects {
			proj.InsertImport(prj.Name, prj.Version, prj.Branch, prj.Revision, true)
		}
	}
	return nil
}

type gopkgToml struct {
	Required   []string          `toml:"required"`
	Constraint []gopkgConstraint `toml:"constraint"`
}

type gopkgConstraint struct {
	Name    string
	Version string
	Branch  string
}

type gopkgLock struct {
	Projects []gopkgLockProjects `toml:"projects"`
}

type gopkgLockProjects struct {
	Name     string
	Version  string
	Branch   string
	Packages []string
	Revision string
}
