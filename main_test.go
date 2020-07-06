package maestro

import (
	"context"
	"errors"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/jexia/maestro/internal/codec/json"
	"github.com/jexia/maestro/internal/flow"
	"github.com/jexia/maestro/pkg/core/api"
	"github.com/jexia/maestro/pkg/core/instance"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/providers"
	"github.com/jexia/maestro/pkg/providers/hcl"
	"github.com/jexia/maestro/pkg/providers/mock"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport"
	"github.com/jexia/maestro/pkg/transport/http"
)

func TestNewOptions(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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

func TestNewClientNilOptions(t *testing.T) {
	t.Parallel()

	maestro, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}

	if maestro == nil {
		t.Fatal("nil client returned")
	}
}

func TestNewClientMiddlewareError(t *testing.T) {
	t.Parallel()

	expected := errors.New("middleware")
	_, err := New(WithMiddleware(func(ctx instance.Context) ([]api.Option, error) {
		return nil, expected
	}))

	if err == nil {
		t.Fatal("unexpected pass")
	}

	if err != expected {
		t.Fatalf("unexpected err %s, expected %s", err, expected)
	}
}

func TestNewClientMiddlewareOptions(t *testing.T) {
	t.Parallel()

	expected := "mock"
	client, err := New(WithMiddleware(func(ctx instance.Context) ([]api.Option, error) {
		result := []api.Option{
			WithFunctions(functions.Custom{
				expected: func(args ...*specs.Property) (*specs.Property, functions.Exec, error) {
					return nil, nil, nil
				},
			}),
		}

		return result, nil
	}))

	if err != nil {
		t.Fatal(err)
	}

	_, has := client.Options.Functions[expected]
	if !has {
		t.Fatal("expected function was not set")
	}
}

func TestServe(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

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
	t.Parallel()

	resolver := func(instance.Context) ([]*specs.ServicesManifest, error) { return nil, errors.New("unexpected") }
	_, err := New(
		WithServices(resolver),
	)

	if err == nil {
		t.Fatal("unexpected pass expected error to be returned")
	}
}

func TestNewFlowsErr(t *testing.T) {
	t.Parallel()

	resolver := func(instance.Context) ([]*specs.FlowsManifest, error) { return nil, errors.New("unexpected") }
	_, err := New(
		WithFlows(resolver),
	)

	if err == nil {
		t.Fatal("unexpected pass expected error to be returned")
	}
}

func TestNewGetCollection(t *testing.T) {
	t.Parallel()

	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	result := client.Collection()
	if result == nil {
		t.Fatal("unexpected empty collection")
	}
}

func TestClosingRunningFlows(t *testing.T) {
	t.Parallel()

	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	timeout := 100 * time.Microsecond
	run := make(chan struct{})

	ctx := instance.NewContext()
	manager := flow.NewManager(ctx, "", nil, nil, nil, &flow.ManagerMiddleware{
		AfterDo: func(ctx context.Context, manager *flow.Manager, store refs.Store) (context.Context, error) {
			close(run)
			time.Sleep(timeout)
			return ctx, nil
		},
	})

	client.transporters = []*transport.Endpoint{
		{
			Flow: manager,
		},
	}

	go manager.Do(context.Background(), refs.NewReferenceStore(0))
	<-run

	start := time.Now()
	client.Close()
	diff := time.Now().Sub(start)
	if diff < timeout/2 {
		t.Fatalf("close did not wait for flow to finish execution, diff %+v", diff)
	}
}

func TestClosingEmptyFlows(t *testing.T) {
	t.Parallel()

	client, err := New()
	if err != nil {
		t.Fatal(err)
	}

	client.transporters = []*transport.Endpoint{
		{},
		{},
		{},
	}

	client.Close()
}
