package hcl

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/utils"
)

const (
	pass = "pass"
	fail = "fail"
)

// TestParseSpecs reads all available test cases inside the tests directory.
// The test is expected to pass/fail based on the file extension.
func TestParseSpecs(t *testing.T) {
	path, err := filepath.Abs("./tests/*.hcl")
	if err != nil {
		t.Fatal(err)
	}

	files, err := utils.ResolvePath(path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			ctx := context.Background()
			ctx = logger.WithValue(ctx)

			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			reader, err := os.Open(file.Path)
			if err != nil {
				t.Error(err)
			}

			manifests, err := UnmarshalHCL(ctx, file.Name(), reader)
			if strings.HasSuffix(clean, pass) && err != nil {
				t.Errorf("expected test to pass but failed instead %s, %v", file.Name(), err)
			}

			if strings.HasSuffix(clean, fail) && err == nil {
				t.Errorf("expected test to fail but passed instead %s", file.Name())
			}

			_, err = ParseSpecs(ctx, manifests, nil)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
