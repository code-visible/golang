package utils

import (
	"errors"
	"os"
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
