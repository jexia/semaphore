package mock

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/utils"
)

const (
	pass = "pass"
	fail = "fail"
)

func TestUnmarshalFile(t *testing.T) {
	path, err := filepath.Abs("./tests/*.yaml")
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

			path := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]

			resolver := SchemaResolver(file.Path)
			err := resolver(ctx, schema.NewStore(ctx))

			if strings.HasSuffix(path, pass) && err != nil {
				t.Errorf("expected test to pass but failed instead %s, %v", file.Name(), err)
			}

			if strings.HasSuffix(path, fail) && err == nil {
				t.Errorf("expected test to fail but passed instead %s", file.Name())
			}
		})
	}
}
