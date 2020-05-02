package mock

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/maestro/pkg/definitions"
	"github.com/jexia/maestro/pkg/instance"
)

func TestSchemaParsing(t *testing.T) {
	path, err := filepath.Abs("./tests/*.yaml")
	if err != nil {
		t.Fatal(err)
	}

	files, err := definitions.ResolvePath(path)
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
