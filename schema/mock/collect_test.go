package mock

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/maestro/internal/instance"
	"github.com/jexia/maestro/internal/utils"
	"github.com/jexia/maestro/schema"
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
			ctx := instance.NewContext()
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
