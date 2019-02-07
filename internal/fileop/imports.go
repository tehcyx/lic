package fileop

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

func ReadImports(filePath, importStart, importEnd, importInline string) ([]string, error) {
	var imports []string

	if err := Exists(filePath); err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			msg := fmt.Errorf("Something went wrong opening the file %s", filePath)
			log.Printf(msg.Error())
			return []string{}, msg
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
					return imports, nil
				} else {
					re := regexp.MustCompile(importInline) //TODO: this problably doesn't cover single line imports
					stringMatch := re.FindStringSubmatch(lineText)
					imports = append(imports, stringMatch[0])
					continue
				}
			}
		}
	}

	return imports, nil
}
