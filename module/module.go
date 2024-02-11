package module

import (
	"fmt"
	"go/token"

	"github.com/code-visible/golang/sourcepkg"
	"github.com/code-visible/golang/utils"
)

type Module struct {
	Name  string   `json:"name"`
	Path  string   `json:"path"`
	Files []string `json:"files"`

	fs   *token.FileSet
	pkgs map[string]*sourcepkg.SourcePkg
}

// initialize module
func NewModule(name string, path string) (*Module, error) {
	// make sure the given path is a directory
	err := utils.MustDir(path)
	if err != nil {
		return nil, err
	}

	// initialize module struct
	m := &Module{
		Name:  name,
		Path:  path,
		Files: nil,
		fs:    token.NewFileSet(),
		pkgs:  make(map[string]*sourcepkg.SourcePkg),
	}

	return m, nil
}

func (m *Module) ScanFiles() {
	dirs := utils.ListDirs(m.Path, true)

	for _, d := range dirs {
		pkg, err := sourcepkg.NewSourcePkg(d, m.fs)
		if err != nil {
			fmt.Printf("meet error while parse package, skipped, error: %s\n", err)
			continue
		}
		m.pkgs[pkg.Name] = pkg
		pkg.ParseFiles()
	}

	// list all exported files
	m.fs.Iterate(func(f *token.File) bool {
		m.Files = append(m.Files, f.Name())
		return true
	})
}
