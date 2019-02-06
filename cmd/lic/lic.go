package main

import (
	"fmt"
	"os"
)

func main() {
	command := cmd.NewLicCmd(core.NewOptions())

	err := command.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
