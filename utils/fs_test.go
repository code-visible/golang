package utils

import (
	"os"
	"testing"
)

func TestMustDir(t *testing.T) {
	if MustDir("fs.go") == nil {
		t.Fail()
	}
	if MustDir("..") != nil {
		t.Fail()
	}
}

func TestCountGoFiles(t *testing.T) {
	list, err := ListGoFiles("../utils", true)
	if err != nil {
		t.Error(err)
	}
	if len(list) != 2 {
		t.Fail()
	}
	list, err = ListGoFiles("../utils", false)
	if err != nil {
		t.Error(err)
	}
	if len(list) != 1 {
		t.Fail()
	}
}

func TestListDirs(t *testing.T) {
	_ = os.Chdir("../testdata")
	list := ListDirs("multifiles", true)
	if len(list) != 1 {
		t.Fail()
	}
}
