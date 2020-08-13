package semaphore

import (
	"context"
	"errors"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/config"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/flow"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/mock"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport"
	"github.com/jexia/semaphore/pkg/transport/http"
)

func TestNewOptions(t *testing.T) {
	t.Parallel()
	ctx := logger.WithLogger(broker.NewBackground())

	functions := map[string]functions.Intermediate{
		"cdf": nil,
	}

	tests := [][]config.Option{
		{WithFlows(nil), WithSchema(nil)},
		{WithFlows(nil)},
		{WithSchema(nil)},
		{WithFunctions(functions)},
		{WithFlows(nil), WithSchema(nil), WithFunctions(functions)},
	}

	for _, input := range tests {
		_, err := New(ctx, input...)
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

	ctx := logger.WithLogger(broker.NewBackground())
	files, err := providers.ResolvePath(ctx, []string{}, path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			schema := filepath.Join(filepath.Dir(file.Path), clean+".yaml")

			_, err = New(ctx,
				WithFlows(hcl.FlowsResolver(file.Path)),
				WithServices(hcl.ServicesResolver(file.Path)),
				WithSchema(mock.SchemaResolver(schema)),
				WithCodec(json.NewConstructor()),
				WithListener(http.NewListener(":0")),
				WithCaller(http.NewCaller()),
				WithLogLevel("*", "debug"),
			)

			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestNewClientNilOptions(t *testing.T) {
	t.Parallel()

	ctx := logger.WithLogger(broker.NewBackground())
	semaphore, err := New(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	if semaphore == nil {
		t.Fatal("nil client returned")
	}
}

func TestNewClientNilContext(t *testing.T) {
	t.Parallel()

	_, err := New(nil)
	if err == nil {
		t.Fatal("unexpected pass")
	}
}

func TestNewClientMiddlewareError(t *testing.T) {
	t.Parallel()
	ctx := logger.WithLogger(broker.NewBackground())

	expected := errors.New("middleware")
	_, err := New(ctx, WithMiddleware(func(ctx *broker.Context) ([]config.Option, error) {
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
	ctx := logger.WithLogger(broker.NewBackground())

	expected := "mock"
	client, err := New(ctx, WithMiddleware(func(ctx *broker.Context) ([]config.Option, error) {
		result := []config.Option{
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

	ctx := logger.WithLogger(broker.NewBackground())
	files, err := providers.ResolvePath(ctx, []string{}, path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			clean := file.Name()[:len(file.Name())-len(filepath.Ext(file.Name()))]
			schema := filepath.Join(filepath.Dir(file.Path), clean+".yaml")

			client, err := New(ctx,
				WithFlows(hcl.FlowsResolver(file.Path)),
				WithServices(hcl.ServicesResolver(file.Path)),
				WithSchema(mock.SchemaResolver(schema)),
				WithCodec(json.NewConstructor()),
				WithListener(http.NewListener(":0")),
				WithCaller(http.NewCaller()),
				WithLogLevel("*", "debug"),
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

	ctx := logger.WithLogger(broker.NewBackground())
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

			client, err := New(ctx,
				WithFlows(hcl.FlowsResolver(file.Path)),
				WithServices(hcl.ServicesResolver(file.Path)),
				WithSchema(mock.SchemaResolver(schema)),
				WithCodec(json.NewConstructor()),
				WithListener(http.NewListener(listener.Addr().String())),
				WithCaller(http.NewCaller()),
				WithLogLevel("*", "debug"),
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
	ctx := logger.WithLogger(broker.NewBackground())

	client, err := New(ctx)
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
	ctx := logger.WithLogger(broker.NewBackground())

	resolver := func(*broker.Context) (specs.ServiceList, error) { return nil, errors.New("unexpected") }
	_, err := New(ctx,
		WithServices(resolver),
	)

	if err == nil {
		t.Fatal("unexpected pass expected error to be returned")
	}
}

func TestNewFlowsErr(t *testing.T) {
	t.Parallel()
	ctx := logger.WithLogger(broker.NewBackground())

	resolver := func(*broker.Context) (specs.FlowListInterface, error) { return nil, errors.New("unexpected") }
	_, err := New(ctx,
		WithFlows(resolver),
	)

	if err == nil {
		t.Fatal("unexpected pass expected error to be returned")
	}
}

func TestNewGetEndpoints(t *testing.T) {
	t.Parallel()
	ctx := logger.WithLogger(broker.NewBackground())

	client, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result := client.GetEndpoints()
	if result == nil {
		t.Fatal("unexpected empty endpoints")
	}
}

func TestNewGetFlows(t *testing.T) {
	t.Parallel()
	ctx := logger.WithLogger(broker.NewBackground())

	client, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result := client.GetFlows()
	if result == nil {
		t.Fatal("unexpected empty flows")
	}
}

func TestNewGetServices(t *testing.T) {
	t.Parallel()
	ctx := logger.WithLogger(broker.NewBackground())

	client, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result := client.GetServices()
	if result == nil {
		t.Fatal("unexpected empty services")
	}
}

func TestNewGetSchemas(t *testing.T) {
	t.Parallel()
	ctx := logger.WithLogger(broker.NewBackground())

	client, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	result := client.GetSchemas()
	if result == nil {
		t.Fatal("unexpected empty schemas")
	}
}

func TestClosingRunningFlows(t *testing.T) {
	t.Parallel()
	ctx := logger.WithLogger(broker.NewBackground())

	client, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	timeout := 100 * time.Millisecond
	run := make(chan struct{})

	manager := flow.NewManager(ctx, "", nil, nil, nil, &flow.ManagerMiddleware{
		AfterDo: func(ctx context.Context, manager *flow.Manager, store references.Store) (context.Context, error) {
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

	go manager.Do(context.Background(), references.NewReferenceStore(0))
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
	ctx := logger.WithLogger(broker.NewBackground())

	client, err := New(ctx)
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
