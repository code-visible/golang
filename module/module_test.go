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
	}
	if m == nil || len(m.Files) == 0 {
		t.Fail()
	}
}
