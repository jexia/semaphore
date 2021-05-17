package providers

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
)

func TempFolderStructure(t *testing.T, files []string, linked map[string]string) string {
	prefix := "semaphore"
	root := filepath.Join(os.TempDir(), prefix+strconv.Itoa(int(time.Now().Unix())))

	result := make([]string, len(files))

	for index, file := range files {
		path := filepath.Join(root, file)
		os.MkdirAll(filepath.Dir(path), os.ModePerm)

		_, err := os.Create(path)
		if err != nil {
			t.Fatal(err)
		}

		result[index] = path
	}

	for old, new := range linked {
		link := filepath.Join(root, new)
		err := os.Symlink(filepath.Join(root, old), link)
		if err != nil {
			t.Fatal(err)
		}

		result = append(result, link)
	}

	t.Cleanup(func() {
		for _, file := range result {
			os.Remove(file)
		}
	})

	return root
}

func TestCleanPattern(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"/mock/path/**/nested": "/mock/path/",
		"**/nested":            "",
		"/mock/**/nested":      "/mock/",
		"/mock**/nested":       "/mock",
	}

	for input, expected := range tests {
		t.Run(input, func(t *testing.T) {
			result := CleanPattern(input)
			if result != expected {
				t.Fatalf("unexpected result %s, expected %s", result, expected)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	files := []string{
		"mock/config.hcl",
		"mock/schemas/schema.hcl",
		"mock/schemas/flow.hcl",
		"mock/proto/main.proto",
		"mock/proto/sub.proto",
		"mock/proto/module.proto",
	}

	symlinks := map[string]string{
		"mock/config.hcl":       "config.hcl",
		"config.hcl":            "linked.hcl",
		"mock/proto/main.proto": "main.proto",
		"mock":                  "symdir",
	}

	root := TempFolderStructure(t, files, symlinks)

	tests := map[string]int{
		"mock/*.hcl":         1,
		"mock/**/*.hcl":      2,
		"mock/**/fl*.hcl":    1,
		"mock/**/sch*.hcl":   1,
		"symdir/**/sch*.hcl": 1,
		"symdir/**/fl*.hcl":  1,
		"mock/proto/*.proto": 3,
		"mock":               0,
		"config.hcl":         1,
		"main.proto":         1,
	}

	for pattern, expected := range tests {
		t.Run(pattern, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())

			files, err := ResolvePath(ctx, []string{}, filepath.Join(root, pattern))
			if err != nil {
				t.Fatal(err)
			}

			if len(files) != expected {
				t.Fatalf("unexpected files %+v, expected %d", files, expected)
			}
		})
	}
}

func TestResolvePathErr(t *testing.T) {
	files := []string{
		"mock/config.hcl",
	}

	symlinks := map[string]string{
		"brokensym": "broken",
	}

	root := TempFolderStructure(t, files, symlinks)

	tests := []string{
		"broken",
	}

	for _, pattern := range tests {
		t.Run(pattern, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())

			_, err := ResolvePath(ctx, []string{}, filepath.Join(root, pattern))
			if err == nil {
				t.Fatal("unexpected pass")
			}
		})
	}
}
