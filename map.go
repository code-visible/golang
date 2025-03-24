package golang

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var errInvalidModuleName = errors.New("invalid module name")

// read-only source map of the project
type SourceMap struct {
	module    string
	path      string
	directory string
	dirs      map[string]*SourceDir
	fs        []*SourceFile
	fset      *token.FileSet
	excludes  map[string]byte
}

func NewSourceMap(project, directory, excludes, module string) *SourceMap {
	var err error
	if module == "" {
		module, err = parseModuleName(filepath.Join(project, "go.mod"))
		if err != nil {
			panic(err)
		}
	}

	es := map[string]byte{}
	excludesSplit := strings.Split(excludes, ",")
	for _, item := range excludesSplit {
		item = strings.Trim(item, " ")
		if item != "" {
			es[item] = 0
		}
	}

	sm := &SourceMap{
		module:    module,
		path:      project,
		directory: filepath.ToSlash(directory),
		dirs:      make(map[string]*SourceDir),
		fs:        make([]*SourceFile, 0, 64),
		fset:      token.NewFileSet(),
		excludes:  es,
	}

	return sm
}

func (sm *SourceMap) Scan() {
	sm.normalizeExcludes()
	sm.walk()
	sm.parseFiles()
}

func (sm *SourceMap) normalizeExcludes() {
	for key := range sm.excludes {
		absp, err := filepath.Abs(key)
		if err != nil {
			fmt.Println("warn: unexpected exclude path: ", err)
			continue
		}
		sm.excludes[absp] = 0
	}
}

// Walk scan all the possible sub directories and files of given path.
func (sm *SourceMap) walk() {
	// TODO: handle error
	err := filepath.WalkDir(sm.directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		absp, _ := filepath.Abs(path)
		parent := filepath.Dir(absp)

		if d.IsDir() {
			if _, ok := sm.excludes[absp]; ok {
				return nil
			}
			if _, ok := sm.excludes[parent]; ok {
				sm.excludes[absp] = 0
				return nil
			}
			dir := &SourceDir{
				Path:  filepath.ToSlash(path),
				Files: 0,
				Pkg:   false,
			}
			sm.dirs[absp] = dir
		} else {
			// ignore the files in the exclude directories
			if _, ok := sm.dirs[parent]; !ok {
				return nil
			}
			sm.fs = append(sm.fs, &SourceFile{
				Path: filepath.ToSlash(parent),
				Name: filepath.ToSlash(filepath.Base(path)),
				Dir:  sm.dirs[parent],
			})
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (sm *SourceMap) parseFiles() {
	wg := sync.WaitGroup{}
	wg.Add(len(sm.fs))
	for _, f := range sm.fs {
		f.Dir.Files++
		go func(file *SourceFile) {
			file.Parse2AST(sm.fset)
			wg.Done()
		}(f)
	}
	wg.Wait()
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
	var dirs = make([]*SourceDir, 0, len(sm.dirs))
	for _, d := range sm.dirs {
		dirs = append(dirs, d)
	}
	return dirs
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
