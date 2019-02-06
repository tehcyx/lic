package main

import (
	"fmt"
	"os"

	"github.com/tehcyx/lic/pkg/lic/cmd"
	"github.com/tehcyx/lic/pkg/lic/core"
)

func main() {
	command := cmd.NewLicCmd(core.NewOptions())

	err := command.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
