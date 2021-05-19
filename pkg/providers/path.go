package providers

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"go.uber.org/zap"
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
func ResolvePath(ctx *broker.Context, ignore []string, pattern string) (files []*FileInfo, _ error) {
	resolved := map[string]struct{}{}
	for _, path := range ignore {
		resolved[path] = struct{}{}
	}

	logger.Debug(ctx, "resolve path", zap.String("pattern", pattern))

	dir := filepath.Dir(CleanPattern(pattern))
	return walk(ctx, dir, dir, resolved, pattern)
}

func walk(ctx *broker.Context, absolute string, target string, resolved map[string]struct{}, pattern string) (files []*FileInfo, _ error) {
	pattern = filepath.Clean(pattern)
	logger.Debug(ctx, "resolve pattern", zap.String("dir", target))

	err := filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fname, err := filepath.Rel(target, path); err == nil {
			path = filepath.Join(absolute, fname)
		}

		logger.Debug(ctx, "matching path",
			zap.String("path", path),
			zap.String("pattern", pattern),
		)

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

		logger.Debug(ctx, "file matched pattern", zap.String("path", path))

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
