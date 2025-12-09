package fileop

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExistsSuccess(t *testing.T) {
	dname, err := os.MkdirTemp("", "")
	defer os.Remove(dname)
	if err != nil {
		t.Fatalf("couldn't create temp dir")
	}
	fname := filepath.Join(dname, "exists.txt")
	err = os.WriteFile(fname, []byte("Bar"), 0644)
	defer os.Remove(fname)
	if err != nil {
		t.Fatalf("couldn't create temp file")
	}

	if Exists(fname) != nil {
		t.Errorf("file should exist\n")
	}
}

func TestExistsFail(t *testing.T) {
	dname, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatalf("couldn't create temp dir")
	}
	defer os.Remove(dname)
	fname := filepath.Join(dname, "exists.txt")

	if Exists(fname) == nil {
		t.Errorf("file should not exist\n")
	}
}
