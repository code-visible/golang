package sourcepkg

import (
	"go/token"
	"os"
	"testing"
)

func TestCurrentDirectoryAsPkg(t *testing.T) {
	_ = os.Chdir("../testdata")
	p, err := NewSourcePkg("multifiles", token.NewFileSet())
	if err != nil {
		t.Error(err)
	}
	if p.Name != "main" {
		t.Fail()
	}
	if len(p.Files) != 2 || len(p.parsed.Files) != 2 {
		t.Fail()
	}
}
