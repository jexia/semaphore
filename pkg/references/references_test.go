package references

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/semaphore/pkg/compare"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/mock"
	"github.com/jexia/semaphore/pkg/schema"
)

const (
	pass = "pass"
	fail = "fail"
)

func TestUnmarshalFile(t *testing.T) {
	path, err := filepath.Abs("./tests/*.hcl")
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

			flows, err := hcl.FlowsResolver(file.Path)(ctx)
			if err != nil {
				t.Fatal(err)
			}

			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			path := filepath.Join(filepath.Dir(file.Path), clean+".yaml")

			collection, err := mock.CollectionResolver(path)
			if err != nil {
				t.Fatal(err)
			}

			services, err := mock.ServicesResolver(path)(ctx)
			if err != nil {
				t.Fatal(err)
			}

			schemas, err := mock.SchemaResolver(path)(ctx)
			if err != nil {
				t.Fatal(err)
			}

			err = func() (err error) {
				err = schema.Resolve(ctx, services, schemas, flows)
				if err != nil {
					return err
				}

				err = Resolve(ctx, flows)
				if err != nil {
					return err
				}

				return nil
			}()

			if strings.HasSuffix(clean, pass) && err != nil {
				t.Fatalf("expected test to pass but failed instead %s, %v", file.Name(), err)
			}

			if strings.HasSuffix(clean, fail) && err != nil {
				if err.Error() != collection.Exception.Message {
					t.Fatalf("unexpected error message %s, expected %s", err, collection.Exception.Message)
				}

				return
			}

			err = compare.Types(ctx, services, schemas, flows)

			if strings.HasSuffix(clean, pass) && err != nil {
				t.Fatalf("expected test to pass but failed instead %s, %v", file.Name(), err)
			}

			if strings.HasSuffix(clean, fail) && err == nil {
				t.Fatalf("expected test to fail but passed instead %s", file.Name())
			}

			if strings.HasSuffix(clean, fail) {
				if err.Error() != collection.Exception.Message {
					t.Fatalf("unexpected error message %s, expected %s", err, collection.Exception.Message)
				}
			}
		})
	}
}
