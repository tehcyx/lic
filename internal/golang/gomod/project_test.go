package gomod

import (
	"testing"
)

var (
	modPackageName = `module github.com/tehcyx/imaginary-api

require (
	github.com/tehcyx/imaginary-service
)`
)

func TestReadPackageNameSuccess(t *testing.T) {
	// dname, err := ioutil.TempDir("", "")
	// defer os.Remove(dname)
	// if err != nil {
	// 	t.Fatalf("couldn't create temp dir")
	// }
	// fname := filepath.Join(dname, "go.mod")
	// err = ioutil.WriteFile(fname, []byte(modPackageName), 0644)
	// defer os.Remove(fname)
	// if err != nil {
	// 	t.Fatalf("couldn't create temp file")
	// }

	// packageName, err := ReadPackageName(fname)
	// if err != nil {
	// 	t.Errorf("couldn't process file")
	// }
	// if packageName != "github.com/tehcyx/imaginary-api" {
	// 	t.Errorf("package name should be '%s', something went wrong", "github.com/tehcyx/imaginary-api")
	// }
}

func TestReadPackageNameFailFileEmpty(t *testing.T) {
	// dname, err := ioutil.TempDir("", "")
	// defer os.Remove(dname)
	// if err != nil {
	// 	t.Fatalf("couldn't create temp dir")
	// }
	// fname := filepath.Join(dname, "go.mod")
	// err = ioutil.WriteFile(fname, []byte(""), 0644)
	// defer os.Remove(fname)
	// if err != nil {
	// 	t.Fatalf("couldn't create temp file")
	// }

	// packageName, err := ReadPackageName(fname)
	// if err == nil {
	// 	t.Errorf("something went wrong")
	// }
	// if packageName != "" {
	// 	t.Errorf("there's no package name in this file")
	// }
}

func TestReadPackageNameFailFileNotExistent(t *testing.T) {
	// dname, err := ioutil.TempDir("", "")
	// defer os.Remove(dname)
	// if err != nil {
	// 	t.Fatalf("couldn't create temp dir")
	// }
	// fname := filepath.Join(dname, "go.mod")

	// packageName, err := ReadPackageName(fname)
	// if err == nil {
	// 	t.Errorf("the file does not exist, error should be there")
	// }
	// if packageName != "" {
	// 	t.Errorf("file does not exist, so this should be empty")
	// }
}
