package utils

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var ErrNotDirectory = errors.New("project path should be a directory")

// MustDir promise the given path is a directory
func MustDir(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return ErrNotDirectory
	}
	return nil
}

// list all go files in given directory (not recursively)
// result should be relative path
func ListGoFiles(dir string, includeTestFile bool) ([]string, error) {
	var list []string
	fs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		fname := f.Name()
		if strings.HasSuffix(fname, ".go") {
			if !includeTestFile && strings.HasSuffix(fname, "_test.go") {
				continue
			}
			list = append(list, filepath.Join(dir, fname))
		}
	}

	return list, nil
}

// list all directories
// if parameter "dir" is not a directory, return empty list
func ListDirs(dir string, recursive bool) []string {
	var list []string
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			list = append(list, path)
		}
		return nil
	})
	return list
}
