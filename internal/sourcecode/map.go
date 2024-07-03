package sourcecode

import (
	"bytes"
	"errors"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
)

var errInvalidModuleName = errors.New("invalid module name")

// read-only source map of the project
type SourceMap struct {
	name   string
	path   string
	dirs   []*SourceDir
	fs     []*SourceFile
	fset   *token.FileSet
	dirIdx map[string]int
}

func NewSourceMap(gomod string, path string) *SourceMap {
	if gomod == "" {
		gomod = filepath.Join(path, "go.mod")
	}
	moduleName, err := parseModuleName(gomod)
	if err != nil {
		panic(err)
	}
	sm := &SourceMap{
		name: moduleName,
		path: path,
		fs:   make([]*SourceFile, 0, 64),
		fset: token.NewFileSet(),
	}

	return sm
}

func (sm *SourceMap) Scan() {
	// TODO
	_ = os.Chdir(sm.path)
	sm.walk()
	sm.parseFiles()
}

// Walk scan all the possible sub directories and files of given path.
func (sm *SourceMap) walk() {
	// TODO: handle error
	var current int
	err := filepath.WalkDir(sm.path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		relPath, err := filepath.Rel(sm.path, path)
		if err != nil {
			return errors.New("unknown error while parsing relative path")
		}

		if d.IsDir() {
			current = len(sm.dirs)
			sm.dirs = append(sm.dirs, &SourceDir{
				Path: relPath,
			})
		} else {
			sm.fs = append(sm.fs, &SourceFile{
				Path: filepath.Dir(relPath),
				Name: filepath.Base(relPath),
				Dir:  current,
			})
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (sm *SourceMap) parseFiles() {
	for _, f := range sm.fs {
		f.Parse2AST(sm.fset)
	}
}

func (sm *SourceMap) Path() string {
	return sm.path
}

func (sm *SourceMap) Files() []*SourceFile {
	return sm.fs
}

func (sm *SourceMap) Dirs() []*SourceDir {
	return sm.dirs
}

func (sm *SourceMap) FileSet() *token.FileSet {
	return sm.fset
}

func parseModuleName(gomod string) (string, error) {
	var moduleName = []byte("module")
	mod, err := os.ReadFile(gomod)
	if err != nil {
		panic(err)
	}
	if !bytes.HasPrefix(mod, moduleName) {
		return "", errInvalidModuleName
	}
	mod = mod[len(moduleName)+1:]
	idx := bytes.IndexByte(mod, '\n')
	if idx < 0 {
		return "", errInvalidModuleName
	}
	return string(mod[:idx]), nil
}
