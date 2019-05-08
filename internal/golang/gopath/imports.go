package gopath

import (
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tehcyx/lic/internal/report"
)

//ReadImports reads imports on a given filepath with the given regex params for start, end and line
func ReadImports(proj *report.Project, filePath string) error {
	imports := make(map[string]int)

	err := filepath.Walk(filePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				fset := token.NewFileSet()

				projAST, err := parser.ParseDir(fset, path, nil, parser.ImportsOnly)
				if err != nil {
					log.Println("Something went wrong")
				}

				for _, v := range projAST {
					for _, vv := range v.Files {
						for _, i := range vv.Imports {
							i.Path.Value = strings.Replace(i.Path.Value, "\"", "", -1)
							imports[i.Path.Value] = 1 // save in map, to skip duplicates
						}
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	for in := range imports {
		proj.InsertImport(in, "n/a", "", "", true)
	}
	return nil
}
