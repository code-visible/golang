package utils

import (
	"fmt"
	"strings"
)

func FormatPath(path_ string) string {
	if path_ == "." {
		return "/"
	}
	if !strings.HasPrefix(path_, "/") {
		return fmt.Sprintf("/%s", path_)
	}
	return path_
}
