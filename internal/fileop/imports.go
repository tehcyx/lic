package fileop

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

//ReadImports reads imports on a given filepath with the given regex params for start, end and line
func ReadImports(filePath, importStart, importEnd, importLine, importInline string) (map[string]string, error) {
	var imports map[string]string
	imports = make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		msg := fmt.Errorf("Something went wrong opening the file %s", filePath)
		log.Printf(msg.Error())
		return nil, msg
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var importfound bool

	reInlineImport := regexp.MustCompile(importInline)
	reLineImport := regexp.MustCompile(importLine)
	for scanner.Scan() {
		lineText := scanner.Text()
		if importfound == false {
			if match, _ := regexp.MatchString(importStart, lineText); match {
				importfound = true
			}
			stringMatch := reInlineImport.FindStringSubmatch(lineText)
			if len(stringMatch) > 1 {
				imports[stringMatch[1]] = stringMatch[2]
			}
		} else {
			if match, _ := regexp.MatchString(importEnd, lineText); match {
				break
			} else {
				stringMatch := reLineImport.FindStringSubmatch(lineText)
				if len(stringMatch) > 0 {
					imports[stringMatch[0]] = ""
				}
				continue
			}
		}
	}
	return imports, nil
}
