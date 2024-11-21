package parser

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
	module    string
	path      string
	directory string
	dirs      []*SourceDir
	fs        []*SourceFile
	fset      *token.FileSet
}

func NewSourceMap(project string, directory string) *SourceMap {
	moduleName, err := parseModuleName(filepath.Join(project, "go.mod"))
	if err != nil {
		panic(err)
	}
	sm := &SourceMap{
		module:    moduleName,
		path:      project,
		directory: directory,
		fs:        make([]*SourceFile, 0, 64),
		fset:      token.NewFileSet(),
	}

	return sm
}

func (sm *SourceMap) Scan() {
	sm.walk()
	sm.parseFiles()
}

// Walk scan all the possible sub directories and files of given path.
func (sm *SourceMap) walk() {
	// TODO: handle error
	var current *SourceDir = nil
	err := filepath.WalkDir(sm.directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}

		if d.IsDir() {
			dir := &SourceDir{
				Path:  path,
				Files: 0,
				Pkg:   false,
			}
			current = dir
			sm.dirs = append(sm.dirs, dir)
		} else {
			sm.fs = append(sm.fs, &SourceFile{
				Path: filepath.Dir(path),
				Name: filepath.Base(path),
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
		f.Dir.Files++
		f.Parse2AST(sm.fset)
	}
}

func (sm *SourceMap) Module() string {
	return sm.module
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
