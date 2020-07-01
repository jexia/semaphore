package mock

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/providers"
)

func TestSchemaParsing(t *testing.T) {
	path, err := filepath.Abs("./tests/*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	ctx := instance.NewContext()
	files, err := providers.ResolvePath(ctx, []string{}, path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			ctx := instance.NewContext()
			path := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]

			if strings.HasSuffix(path, fail) {
				return
			}

			var err error

			services := ServicesResolver(file.Path)
			_, err = services(ctx)
			if err != nil {
				t.Errorf("unexpected err while resolving services %s, %v", file.Name(), err)
			}

			schema := SchemaResolver(file.Path)
			_, err = schema(ctx)
			if err != nil {
				t.Errorf("unexpected err while resolving schema %s, %v", file.Name(), err)
			}
		})
	}
}
