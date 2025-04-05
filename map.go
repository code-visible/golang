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
	types     map[string]byte
	dot       bool
}

func NewSourceMap(project, directory, excludes, module, types string) *SourceMap {
	var err error
	if module == "" {
		module, err = parseModuleName(filepath.Join(project, "go.mod"))
		if err != nil {
			panic(err)
		}
	}

	dot := true
	es := map[string]byte{}
	excludesSplit := strings.Split(excludes, ",")
	for _, item := range excludesSplit {
		item = strings.Trim(item, " ")
		if item == ".*" {
			dot = false
			continue
		}
		if item != "" {
			es[item] = 0
		}
	}

	typs := map[string]byte{}
	typesSplit := strings.Split(types, ",")
	for _, item := range typesSplit {
		item = strings.Trim(item, " ")
		if item != "" {
			typs[item] = 0
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
		types:     typs,
		dot:       dot,
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
	err := filepath.Walk(sm.directory, func(path string, d fs.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		absp, _ := filepath.Abs(path)
		parent := filepath.Dir(absp)
		baseName := filepath.Base(path)
		// ignore dot files and directories
		if !sm.dot && baseName != "." && strings.HasPrefix(baseName, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		// ignore current directories or files
		if _, ok := sm.excludes[absp]; ok {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			dir := &SourceDir{
				Path:  filepath.ToSlash(path),
				Files: 0,
				Pkg:   false,
			}
			sm.dirs[absp] = dir
		} else {
			// check if the file should be parsed
			if !sm.shouldParseFile(baseName) {
				return nil
			}
			sm.fs = append(sm.fs, &SourceFile{
				Path: filepath.ToSlash(parent),
				Name: filepath.ToSlash(baseName),
				Dir:  sm.dirs[parent],
			})
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (sm *SourceMap) shouldParseFile(fileName string) bool {
	if _, ok := sm.types[fileName]; ok {
		return true
	}
	ext := filepath.Ext(fileName)
	if ext == ".go" {
		return true
	}
	if _, ok := sm.types[fmt.Sprintf("*%s", ext)]; ok {
		return true
	}
	return false
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
