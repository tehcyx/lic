package fileop

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

//FilesInPath given a path and a regex pattern string, the function will return the filepaths that match the pattern
func FilesInPath(rootPath, filePattern string) ([]string, error) {
	var files []string
	vendorPath := filepath.Join(rootPath, "vendor")

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Skip vendor directory
		if info.IsDir() && path == vendorPath {
			return filepath.SkipDir
		}
		if match, _ := regexp.MatchString(filePattern, path); match {
			if !strings.HasPrefix(path, vendorPath+string(filepath.Separator)) {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		return files, fmt.Errorf("error collecting file information: %w", err)
	}

	return files, nil
}
