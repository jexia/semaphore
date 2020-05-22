package definitions

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/logger"
)

// A FileInfo describes a file
type FileInfo struct {
	os.FileInfo
	Path string
}

var nestedWildcard = regexp.MustCompile(`(?m)\*\*\/.+$`)

// CleanPattern removes any nested pattern from the given path
func CleanPattern(path string) string {
	return nestedWildcard.ReplaceAllString(path, "")
}

// ResolvePath resolves the given path and returns the matching pattern files.
func ResolvePath(ctx instance.Context, ignore []string, pattern string) (files []*FileInfo, _ error) {
	resolved := map[string]struct{}{}
	for _, path := range ignore {
		resolved[path] = struct{}{}
	}

	return walk(ctx, resolved, pattern)
}

func walk(ctx instance.Context, resolved map[string]struct{}, pattern string) (files []*FileInfo, _ error) {
	dir := filepath.Dir(CleanPattern(pattern))
	pattern = filepath.Clean(pattern)

	ctx.Logger(logger.Core).WithField("dir", dir).Debug("Resolve pattern")

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		ctx.Logger(logger.Core).WithField("path", path).Debug("Matching path")

		_, has := resolved[path]
		if has {
			return nil
		}

		resolved[path] = struct{}{}

		if info.IsDir() {
			return nil
		}

		matched, err := filepath.Match(pattern, path)
		if err != nil {
			return err
		}

		if !matched {
			return nil
		}

		ctx.Logger(logger.Core).WithField("path", path).Debug("File matched pattern")

		files = append(files, &FileInfo{
			FileInfo: info,
			Path:     path,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
