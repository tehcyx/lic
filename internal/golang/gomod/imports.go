package gomod

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tehcyx/lic/internal/report"
)

var (
	modLine          = regexp.MustCompile(`^module (?P<package>\S+|\/|\.)+$`)
	modInlineRequire = regexp.MustCompile(`^require (?P<import>\S+|\/|\.)+ (?P<version>\S+|\/|\.)+`)
	modRequire       = regexp.MustCompile(`^require \($`)

	requireLine = regexp.MustCompile(`(?P<import>\S+|\/|\.)+ (?P<version>\S+|\/|\.)+( \/\/ (?P<indirect>indirect)){0,1}`)

	closingBracket = regexp.MustCompile(`^\)$`)

	//GoFileExtension holds the pattern for the file extensions that should be included for import scans
	GoFileExtension = ".*\\.go$"
)

//ReadImports reads imports on a given filepath with the given regex params for start, end and line
func ReadImports(proj *report.Project, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Errorf("Something went wrong opening the file %s", filePath)
		log.Printf(msg.Error())
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineText := scanner.Text()
		switch {
		case modLine.MatchString(lineText):
			lineText = strings.TrimSpace(lineText)
			match := modLine.FindStringSubmatch(lineText)
			matchResult := make(map[string]string)
			for i, name := range modLine.SubexpNames() {
				if i != 0 && name != "" {
					matchResult[name] = match[i]
				}
			}
			proj.Name = matchResult["package"]
		case modInlineRequire.MatchString(lineText):
			lineText = strings.TrimSpace(lineText)
			match := modInlineRequire.FindStringSubmatch(lineText)
			matchResult := make(map[string]string)
			for i, name := range modInlineRequire.SubexpNames() {
				if i != 0 && name != "" {
					matchResult[name] = match[i]
				}
			}
			proj.InsertImport(matchResult["import"], matchResult["version"], "", "", (matchResult["indirect"] == ""))
		case modRequire.MatchString(lineText):
			for scanner.Scan() {
				lineText = scanner.Text()
				switch {
				case requireLine.MatchString(lineText):
					lineText = strings.TrimSpace(lineText)
					match := requireLine.FindStringSubmatch(lineText)
					matchResult := make(map[string]string)
					for i, name := range requireLine.SubexpNames() {
						if i != 0 && name != "" {
							matchResult[name] = match[i]
						}
					}
					proj.InsertImport(matchResult["import"], matchResult["version"], "", "", (matchResult["indirect"] == ""))
				case closingBracket.MatchString(lineText):
					break
				}
			}
		default:
			continue
		}
	}
	return nil
}
