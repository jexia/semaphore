package definitions

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/sirupsen/logrus"
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

	dir := filepath.Dir(CleanPattern(pattern))
	return walk(ctx, dir, dir, resolved, pattern)
}

func walk(ctx instance.Context, absolute string, target string, resolved map[string]struct{}, pattern string) (files []*FileInfo, _ error) {
	pattern = filepath.Clean(pattern)
	ctx.Logger(logger.Core).WithField("dir", target).Debug("Resolve pattern")

	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fname, err := filepath.Rel(target, path); err == nil {
			path = filepath.Join(absolute, fname)
		}

		ctx.Logger(logger.Core).WithFields(logrus.Fields{
			"path":    path,
			"pattern": pattern,
		}).Debug("Matching path")

		_, has := resolved[path]
		if has {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			link, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}

			info, err := os.Lstat(link)
			if err != nil {
				return err
			}

			if info.IsDir() {
				result, err := walk(ctx, path, link, resolved, pattern)
				if err != nil {
					return err
				}

				files = append(files, result...)
				return nil
			}
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
