package fileop

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	//GoFileImportStart indicates a multiline or single line import, either "import (" or "import \""
	GoFileImportStart = "^import (\\(|\").*"
	//GoFileImportEnd indicates the end of imports, either by a closing bracket ")", a variable definition, a function definition or struct definition
	GoFileImportEnd = "^(\\)|var|func|type).*"
	//GoFileLineImport indicates the import found in a multiline import between the double quotes
	GoFileLineImport = "\"(\\S+|\\/|\\.)+\""
	//GoModInlineImport will cover single line imports that are just "require github.com/user/repo".
	GoFileInlineImport = "^require (\\S+|\\/|\\.)+ (\\S+|\\/|\\.)+"

	importStyleOne = `package main

import (
	"math"
	m "math"
	. "math"
	_ "github.com/tehcyx/imaginary-api"
	"github.com/tehcyx/imaginary-service"
	"fmt"
)

func main() {
	fmt.Println("test")
}
`
)

func TestReadImportsSuccess(t *testing.T) {
	dname, err := ioutil.TempDir("", "")
	defer os.Remove(dname)
	if err != nil {
		t.Fatalf("couldn't create temp dir")
	}
	fname := filepath.Join(dname, "test.go")
	err = ioutil.WriteFile(fname, []byte(importStyleOne), 0644)
	defer os.Remove(fname)
	if err != nil {
		t.Fatalf("couldn't create temp file")
	}

	imports, err := ReadImports(fname, GoFileImportStart, GoFileImportEnd, GoFileLineImport, GoFileInlineImport)
	if err != nil {
		t.Errorf("couldn't process file")
	}
	if len(imports) != 4 {
		t.Errorf("There should be 4 imports, as it's collecting imports in a map. Actual number is %d", len(imports))
	}
}

func TestReadImportsFail(t *testing.T) {
	dname, err := ioutil.TempDir("", "")
	defer os.Remove(dname)
	if err != nil {
		t.Fatalf("couldn't create temp dir")
	}
	fname := filepath.Join(dname, "test.go")

	imports, err := ReadImports(fname, GoFileImportStart, GoFileImportEnd, GoFileLineImport, GoFileInlineImport)
	if err == nil {
		t.Errorf("file does not exist, test should fail")
	}
	if len(imports) != 0 {
		t.Errorf("file does not exist, imports should be nil")
	}
}
