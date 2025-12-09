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
	modLine          = regexp.MustCompile(`^module\s+(?P<package>\S+)$`)
	goDirective      = regexp.MustCompile(`^go\s+\d+\.\d+`)
	modInlineRequire = regexp.MustCompile(`^require\s+(?P<import>\S+)\s+(?P<version>\S+)(\s+//\s+(?P<indirect>indirect))?$`)
	modRequire       = regexp.MustCompile(`^require\s+\($`)

	requireLine = regexp.MustCompile(`^\s*(?P<import>\S+)\s+(?P<version>\S+)(\s+//\s+(?P<indirect>indirect))?$`)

	closingBracket = regexp.MustCompile(`^\s*\)$`)

	//GoFileExtension holds the pattern for the file extensions that should be included for import scans
	GoFileExtension = ".*\\.go$"
)

//ReadImports reads imports on a given filepath with the given regex params for start, end and line
func ReadImports(proj *report.Project, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Errorf("something went wrong opening the file %s: %w", filePath, err)
		log.Printf("%s", msg.Error())
		return msg
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lineText := scanner.Text()
		// Trim leading/trailing whitespace for matching
		trimmedLine := strings.TrimSpace(lineText)

		// Skip empty lines and comments
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "//") {
			continue
		}

		switch {
		case modLine.MatchString(trimmedLine):
			match := modLine.FindStringSubmatch(trimmedLine)
			matchResult := make(map[string]string)
			for i, name := range modLine.SubexpNames() {
				if i != 0 && name != "" {
					matchResult[name] = match[i]
				}
			}
			proj.Name = matchResult["package"]
		case goDirective.MatchString(trimmedLine):
			// Skip go directive (e.g., "go 1.24")
			continue
		case modInlineRequire.MatchString(trimmedLine):
			match := modInlineRequire.FindStringSubmatch(trimmedLine)
			matchResult := make(map[string]string)
			for i, name := range modInlineRequire.SubexpNames() {
				if i != 0 && name != "" {
					matchResult[name] = match[i]
				}
			}
			proj.InsertImport(matchResult["import"], matchResult["version"], "", "", (matchResult["indirect"] != "indirect"))
		case modRequire.MatchString(trimmedLine):
			// Process multi-line require block
			for scanner.Scan() {
				lineText = scanner.Text()
				trimmedLine = strings.TrimSpace(lineText)

				// Skip empty lines
				if trimmedLine == "" {
					continue
				}

				// Check for closing bracket - exit the require block
				if closingBracket.MatchString(trimmedLine) {
					break
				}

				// Check if we hit another require block or other directive - should not happen but be safe
				if modRequire.MatchString(trimmedLine) || modInlineRequire.MatchString(trimmedLine) ||
				   modLine.MatchString(trimmedLine) || goDirective.MatchString(trimmedLine) {
					// This shouldn't happen in well-formed go.mod, but if it does, break out
					break
				}

				// Match dependency lines
				if requireLine.MatchString(trimmedLine) {
					match := requireLine.FindStringSubmatch(trimmedLine)
					matchResult := make(map[string]string)
					for i, name := range requireLine.SubexpNames() {
						if i != 0 && name != "" {
							matchResult[name] = match[i]
						}
					}
					proj.InsertImport(matchResult["import"], matchResult["version"], "", "", (matchResult["indirect"] != "indirect"))
				}
			}
		default:
			continue
		}
	}
	return nil
}
