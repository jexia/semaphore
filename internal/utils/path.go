package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// A FileInfo describes a file
type FileInfo struct {
	os.FileInfo
	Path         string
	AbsolutePath string
}

var nestedWildcard = regexp.MustCompile(`(?m)\*\*\/.+$`)

// CleanPattern removes any nested pattern from the given path
func CleanPattern(path string) string {
	return nestedWildcard.ReplaceAllString(path, "")
}

// ResolvePath resolves the given path and returns the matching pattern files.
func ResolvePath(pattern string) (files []*FileInfo, _ error) {
	dir := filepath.Dir(CleanPattern(pattern))
	pattern = filepath.Clean(pattern)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		matched, err := filepath.Match(pattern, path)
		if err != nil {
			return err
		}

		if !matched {
			return nil
		}

		absolute := strings.Replace(path, dir+"/", "", 1)
		files = append(files, &FileInfo{
			FileInfo:     info,
			Path:         path,
			AbsolutePath: absolute,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
