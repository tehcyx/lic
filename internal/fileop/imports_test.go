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

	modPackageName = `module github.com/tehcyx/imaginary-api

require (
	github.com/tehcyx/imaginary-service
)`
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

	imports, err := ReadImports(fname, GoFileImportStart, GoFileImportEnd, GoFileLineImport)
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

	imports, err := ReadImports(fname, GoFileImportStart, GoFileImportEnd, GoFileLineImport)
	if err == nil {
		t.Errorf("file does not exist, test should fail")
	}
	if len(imports) != 0 {
		t.Errorf("file does not exist, imports should be nil")
	}
}

func TestReadModPackageNameSuccess(t *testing.T) {
	dname, err := ioutil.TempDir("", "")
	defer os.Remove(dname)
	if err != nil {
		t.Fatalf("couldn't create temp dir")
	}
	fname := filepath.Join(dname, "go.mod")
	err = ioutil.WriteFile(fname, []byte(modPackageName), 0644)
	defer os.Remove(fname)
	if err != nil {
		t.Fatalf("couldn't create temp file")
	}

	packageName, err := ReadModPackageName(fname)
	if err != nil {
		t.Errorf("couldn't process file")
	}
	if packageName != "github.com/tehcyx/imaginary-api" {
		t.Errorf("package name should be '%s', something went wrong", "github.com/tehcyx/imaginary-api")
	}
}

func TestReadModPackageNameFailFileEmpty(t *testing.T) {
	dname, err := ioutil.TempDir("", "")
	defer os.Remove(dname)
	if err != nil {
		t.Fatalf("couldn't create temp dir")
	}
	fname := filepath.Join(dname, "go.mod")
	err = ioutil.WriteFile(fname, []byte(""), 0644)
	defer os.Remove(fname)
	if err != nil {
		t.Fatalf("couldn't create temp file")
	}

	packageName, err := ReadModPackageName(fname)
	if err == nil {
		t.Errorf("something went wrong")
	}
	if packageName != "" {
		t.Errorf("there's no package name in this file")
	}
}

func TestReadModPackageNameFailFileNotExistent(t *testing.T) {
	dname, err := ioutil.TempDir("", "")
	defer os.Remove(dname)
	if err != nil {
		t.Fatalf("couldn't create temp dir")
	}
	fname := filepath.Join(dname, "go.mod")

	packageName, err := ReadModPackageName(fname)
	if err == nil {
		t.Errorf("the file does not exist, error should be there")
	}
	if packageName != "" {
		t.Errorf("file does not exist, so this should be empty")
	}
}
