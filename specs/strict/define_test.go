package strict

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/schema/mock"
	"github.com/jexia/maestro/utils"
)

const (
	pass = "pass"
	fail = "fail"
)

func TestUnmarshalFile(t *testing.T) {
	path, err := filepath.Abs("./tests")
	if err != nil {
		t.Fatal(err)
	}

	files, err := utils.ReadDir(path, true, ".hcl")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			reader, err := os.Open(filepath.Join(file.Path, file.Name()))
			if err != nil {
				t.Fatal(err)
			}

			definition, err := hcl.UnmarshalHCL(file.Name(), reader)
			if err != nil {
				t.Fatal(err)
			}

			manifest, err := hcl.ParseManifest(definition, nil)
			if err != nil {
				t.Fatal(err)
			}

			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			file, err := os.Open(filepath.Join(file.Path, clean+".yaml"))
			if err != nil {
				t.Fatal(err)
			}

			collection, err := mock.UnmarshalFile(file)
			if err != nil {
				t.Fatal(err)
			}

			err = Define(collection, manifest)
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
