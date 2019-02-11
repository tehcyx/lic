package fileop

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

//FilesInPath given a path and a regex pattern string, the function will return the filepaths that match the pattern
func FilesInPath(rootPath, filePattern string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if match, _ := regexp.MatchString(filePattern, path); match {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return files, errors.New("Something went wrong collecting file information")
	}

	return files, nil
}
