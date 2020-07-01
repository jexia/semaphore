package maestro

import (
	"errors"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/jexia/maestro/internal/codec/json"
	"github.com/jexia/maestro/pkg/core/api"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/providers"
	"github.com/jexia/maestro/pkg/providers/hcl"
	"github.com/jexia/maestro/pkg/providers/mock"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport/http"
)

func TestNewOptions(t *testing.T) {
	functions := map[string]functions.Intermediate{
		"cdf": nil,
	}

	tests := [][]api.Option{
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

	ctx := instance.NewContext()
	files, err := providers.ResolvePath(ctx, []string{}, path)
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

	ctx := instance.NewContext()
	files, err := providers.ResolvePath(ctx, []string{}, path)
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

func TestErrServe(t *testing.T) {
	path, err := filepath.Abs("./tests/*.hcl")
	if err != nil {
		t.Fatal(err)
	}

	ctx := instance.NewContext()
	files, err := providers.ResolvePath(ctx, []string{}, path)
	if err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	defer listener.Close()

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			schema := filepath.Join(filepath.Dir(file.Path), clean+".yaml")

			client, err := New(
				WithFlows(hcl.FlowsResolver(file.Path)),
				WithServices(hcl.ServicesResolver(file.Path)),
				WithSchema(mock.SchemaResolver(schema)),
				WithCodec(json.NewConstructor()),
				WithListener(http.NewListener(listener.Addr().String(), nil)),
				WithCaller(http.NewCaller()),
				WithLogLevel(logger.Core, "debug"),
			)

			if err != nil {
				t.Fatal(err)
			}

			err = client.Serve()
			if err == nil {
				t.Fatal("unexpected pass expected error to be returned")
			}
		})
	}
}

func TestServeNoListeners(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	err = client.Serve()
	if err == nil {
		t.Fatal("unexpected pass expected error to be returned")
	}
}

func TestNewServiceErr(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.ServicesManifest, error) { return nil, errors.New("unexpected") }
	_, err := New(
		WithServices(resolver),
	)

	if err == nil {
		t.Fatal("unexpected pass expected error to be returned")
	}
}

func TestNewFlowsErr(t *testing.T) {
	resolver := func(instance.Context) ([]*specs.FlowsManifest, error) { return nil, errors.New("unexpected") }
	_, err := New(
		WithFlows(resolver),
	)

	if err == nil {
		t.Fatal("unexpected pass expected error to be returned")
	}
}
