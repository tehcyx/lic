package gomod

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	modFile = `module github.com/tehcyx/lic

require (
		github.com/inconshreveable/mousetrap v1.0.0 // indirect
		github.com/pelletier/go-toml v1.3.0
		github.com/spf13/cobra v0.0.3
		github.com/spf13/pflag v1.0.3 // indirect
)

require (github.com/inconshreveable/mousetrap v1.0.0 // indirect)

require github.com/inconshreveable/mousetrap v1.0.0 // indirect`
)

func TestReadImports(t *testing.T) {
	dname, err := os.MkdirTemp("", "")
	defer os.Remove(dname)
	if err != nil {
		t.Fatalf("couldn't create temp dir")
	}
	fname := filepath.Join(dname, "go.mod")
	err = os.WriteFile(fname, []byte(modFile), 0644)
	defer os.Remove(fname)
	if err != nil {
		t.Fatalf("couldn't create temp file")
	}
}
