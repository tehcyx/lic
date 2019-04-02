package fileop

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//FilesInPath given a path and a regex pattern string, the function will return the filepaths that match the pattern
func FilesInPath(rootPath, filePattern string) ([]string, error) {
	var files []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if match, _ := regexp.MatchString(filePattern, path); match {
			if !strings.HasPrefix(path, fmt.Sprintf("%s%s", rootPath, "vendor")) { // skip vendor folder on root path
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return files, errors.New("Something went wrong collecting file information")
	}

	return files, nil
}
