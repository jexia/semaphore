package mock

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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

	files, err := utils.ReadDir(path, true, ".yaml")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]

			reader, err := os.Open(filepath.Join(file.Path, file.Name()))
			if err != nil {
				t.Error(err)
			}

			_, err = UnmarshalFile(reader)
			if strings.HasSuffix(clean, pass) && err != nil {
				t.Errorf("expected test to pass but failed instead %s, %v", file.Name(), err)
			}

			if strings.HasSuffix(clean, fail) && err == nil {
				t.Errorf("expected test to fail but passed instead %s", file.Name())
			}
		})
	}
}
