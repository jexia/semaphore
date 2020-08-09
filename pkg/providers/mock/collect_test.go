package mock

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/providers"
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

	ctx := logger.WithLogger(broker.NewContext())
	files, err := providers.ResolvePath(ctx, []string{}, path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewContext())
			path := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]

			resolver := SchemaResolver(file.Path)
			_, err := resolver(ctx)

			if strings.HasSuffix(path, pass) && err != nil {
				t.Errorf("expected test to pass but failed instead %s, %v", file.Name(), err)
			}

			if strings.HasSuffix(path, fail) && err == nil {
				t.Errorf("expected test to fail but passed instead %s", file.Name())
			}
		})
	}
}
