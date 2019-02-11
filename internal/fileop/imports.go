package fileop

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

//ReadImports reads imports on a given filepath with the given regex params for start, end and line
func ReadImports(filePath, importStart, importEnd, importInline string) (map[string]int, error) {
	var imports map[string]int
	imports = make(map[string]int)

	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Errorf("Something went wrong opening the file %s", filePath)
		log.Printf(msg.Error())
		return nil, msg
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var importfound bool
	for scanner.Scan() {
		lineText := scanner.Text()
		if importfound == false {
			if match, _ := regexp.MatchString(importStart, lineText); match {
				importfound = true
			}
		} else {
			if match, _ := regexp.MatchString(importEnd, lineText); match {
				break
			} else {
				re := regexp.MustCompile(importInline) //TODO: this problably doesn't cover single line imports
				stringMatch := re.FindStringSubmatch(lineText)
				if len(stringMatch) > 0 {
					imports[stringMatch[0]] = 1
				}
				continue
			}
		}
	}

	return imports, nil
}

//ReadModPackageName reads package name from go.mod file
func ReadModPackageName(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Errorf("Something went wrong opening the file %s", filePath)
		log.Printf(msg.Error())
		return "", msg
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		if strings.Contains(lineText, "module ") {
			re := regexp.MustCompile("^module (.*)$")
			text := re.FindStringSubmatch(lineText)
			if len(text) > 1 {
				return text[1], nil
			}
		}
	}
	return "", fmt.Errorf("Couldn't find module name on path")
}
