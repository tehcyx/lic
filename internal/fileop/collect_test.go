package fileop

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFilesInPathSuccess(t *testing.T) {
	dname, err := ioutil.TempDir("", "")
	defer os.Remove(dname)
	if err != nil {
		t.Fatalf("couldn't create temp dir")
	}
	fname := filepath.Join(dname, "foo.txt")
	err = ioutil.WriteFile(fname, []byte("Bar"), 0644)
	defer os.Remove(fname)
	if err != nil {
		t.Fatalf("couldn't create temp file")
	}

	fname = filepath.Join(dname, "foo1.txt")
	err = ioutil.WriteFile(fname, []byte("Bar"), 0644)
	defer os.Remove(fname)
	if err != nil {
		t.Fatalf("couldn't create temp file")
	}
	fname = filepath.Join(dname, "foo.md")
	err = ioutil.WriteFile(fname, []byte("Bar"), 0644)
	defer os.Remove(fname)
	if err != nil {
		t.Fatalf("couldn't create temp file")
	}

	textFiles, err := FilesInPath(dname, ".*\\.txt")
	if err != nil {
		t.Errorf("error on file tree\n")
	}

	if len(textFiles) != 2 {
		t.Errorf("two txt files on tree\n")
	}

	markdownFiles, err := FilesInPath(dname, ".*\\.md")
	if err != nil {
		t.Errorf("error on file tree\n")
	}

	if len(markdownFiles) != 1 {
		t.Errorf("one md file on tree\n")
	}

	javaFiles, err := FilesInPath(dname, ".*\\.java")
	if err != nil {
		t.Errorf("error on file tree\n")
	}

	if len(javaFiles) != 0 {
		t.Errorf("no java file on tree\n")
	}
}

func TestFilesInPathFail(t *testing.T) {
	dname, err := ioutil.TempDir("", "")
	defer os.Remove(dname)
	_, err = FilesInPath(dname+"notexistent", ".*\\.txt")
	if err == nil {
		t.Errorf("fail\n")
	}
}
