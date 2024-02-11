package module

import (
	"os"
	"testing"
)

func TestCurrentDirectoryAsModule(t *testing.T) {
	_ = os.Chdir("../testdata")
	m, err := NewModule("test", "multifiles")
	if err != nil {
		t.Error(err)
		return
	}
	if m == nil {
		t.FailNow()
	}
	m.ScanFiles()
	if len(m.Files) == 0 {
		t.FailNow()
	}
}
