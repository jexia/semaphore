package maestro

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/jexia/maestro/internal/constructor"
	"github.com/jexia/maestro/internal/logger"
	"github.com/jexia/maestro/internal/utils"
	"github.com/jexia/maestro/pkg/codec/json"
	"github.com/jexia/maestro/pkg/definitions/hcl"
	"github.com/jexia/maestro/pkg/definitions/mock"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/transport/http"
)

func TestOptions(t *testing.T) {
	functions := map[string]functions.Intermediate{
		"cdf": nil,
	}

	ctx := instance.NewContext()
	options := NewOptions(ctx, WithFunctions(functions))

	if len(options.Functions) != len(functions) {
		t.Errorf("unexpected functions %+v, expected %+v", options.Functions, functions)
	}
}

func TestNewOptions(t *testing.T) {
	functions := map[string]functions.Intermediate{
		"cdf": nil,
	}

	tests := [][]constructor.Option{
		{WithFlows(nil), WithSchema(nil)},
		{WithFlows(nil)},
		{WithSchema(nil)},
		{WithFunctions(functions)},
		{WithFlows(nil), WithSchema(nil), WithFunctions(functions)},
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
				WithFlows(hcl.FlowsResolver(file.Path)),
				WithServices(hcl.ServicesResolver(file.Path)),
				WithSchema(mock.SchemaResolver(schema)),
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
				WithFlows(hcl.FlowsResolver(file.Path)),
				WithServices(hcl.ServicesResolver(file.Path)),
				WithSchema(mock.SchemaResolver(schema)),
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
