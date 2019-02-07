package fileop

import (
	"errors"
	"os"
)

//Exists checks if a given filepath exists and is accessible
func Exists(file string) error {
	_, err := os.Stat(file)
	if err == nil {
		// log.Printf("file %s exists", file)
		return nil
	} else if os.IsNotExist(err) {
		// log.Printf("file %s not exists", file)
		return errors.New("File does not exist")
	} else {
		// log.Printf("file %s stat error: %v", file, err)
		return errors.New("File does not exist or permission problem") // can't work with the file, so report failure
	}
}
