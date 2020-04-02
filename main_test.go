package maestro

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/constructor"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/instance"
	"github.com/jexia/maestro/logger"
	"github.com/jexia/maestro/schema/mock"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/transport/http"
	"github.com/jexia/maestro/utils"
)

func TestOptions(t *testing.T) {
	functions := map[string]specs.PrepareFunction{
		"cdf": nil,
	}

	ctx := instance.NewContext()
	options := NewOptions(ctx, WithFunctions(functions))

	if len(options.Functions) != len(functions) {
		t.Errorf("unexpected functions %+v, expected %+v", options.Functions, functions)
	}
}

func TestNewOptions(t *testing.T) {
	functions := map[string]specs.PrepareFunction{
		"cdf": nil,
	}

	tests := [][]constructor.Option{
		{WithDefinitions(nil), WithSchema(nil)},
		{WithDefinitions(nil)},
		{WithSchema(nil)},
		{WithFunctions(functions)},
		{WithDefinitions(nil), WithSchema(nil), WithFunctions(functions)},
	}

	for _, input := range tests {
		_, err := New(input...)
		if err != nil {
			t.Fatalf("unexpected fail %+v", err)
		}
	}
}

func TestNewClient(t *testing.T) {
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
			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			schema := filepath.Join(filepath.Dir(file.Path), clean+".yaml")

			_, err = New(
				WithDefinitions(hcl.DefinitionResolver(file.Path)),
				WithSchema(mock.SchemaResolver(schema)),
				WithSchema(hcl.SchemaResolver(file.Path)),
				WithCodec(json.NewConstructor()),
				WithListener(http.NewListener(":0", nil)),
				WithCaller(http.NewCaller()),
				WithLogLevel(logger.Core, "debug"),
			)

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestServe(t *testing.T) {
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
			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			schema := filepath.Join(filepath.Dir(file.Path), clean+".yaml")

			client, err := New(
				WithDefinitions(hcl.DefinitionResolver(file.Path)),
				WithSchema(mock.SchemaResolver(schema)),
				WithSchema(hcl.SchemaResolver(file.Path)),
				WithCodec(json.NewConstructor()),
				WithListener(http.NewListener(":0", nil)),
				WithCaller(http.NewCaller()),
				WithLogLevel(logger.Core, "debug"),
			)

			if err != nil {
				t.Fatal(err)
			}

			go func() {
				time.Sleep(100 * time.Millisecond)
				client.Close()
			}()

			err = client.Serve()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
