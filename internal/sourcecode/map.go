package sourcecode

import (
	"errors"
	"go/token"
	"io/fs"
	"path/filepath"
)

// read-only source map of the project
type SourceMap struct {
	path string
	dirs []SourceDir
	fs   []SourceFile
	fset *token.FileSet
}

func NewSourceMap(path string) *SourceMap {
	sm := &SourceMap{
		path: path,
		fs:   make([]SourceFile, 0, 64),
		fset: token.NewFileSet(),
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
			sm.dirs = append(sm.dirs, SourceDir{
				Path: relPath,
			})
		} else {
			sm.fs = append(sm.fs, SourceFile{
				Path: relPath,
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

func (sm *SourceMap) Files() []SourceFile {
	return sm.fs
}

func (sm *SourceMap) Dirs() []SourceDir {
	return sm.dirs
}

func (sm *SourceMap) FileSet() *token.FileSet {
	return sm.fset
}
