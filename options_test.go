package semaphore

import (
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport/http"
	"go.uber.org/zap/zapcore"
)

func TestWithFlowsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.FlowListInterface, error) { return nil, nil }

	result, err := NewOptions(ctx, WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.FlowResolvers) != 1 {
		t.Fatal("unexpected result expected flow resolver to be set")
	}
}

func TestWithMultipleFlowsOption(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	resolver := func(*broker.Context) (specs.FlowListInterface, error) { return nil, nil }

	result, err := NewOptions(ctx, WithFlows(resolver), WithFlows(resolver))
	if err != nil {
		t.Fatal(err)
	}

	if len(result.FlowResolvers) != 2 {
		t.Fatal("unexpected result expected multiple flow resolvers to be set")
	}
}

func TestWithLogLevel(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, WithLogLevel("*", "debug"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithFunctions(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, WithFunctions(functions.Custom{"mock": nil}))
	if err != nil {
		t.Fatal(err)
	}

	if options.Functions == nil {
		t.Fatal("functions not set")
	}

	_, has := options.Functions["mock"]
	if !has {
		t.Fatal("mock function does not exist")
	}
}

func TestWithCaller(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, WithCaller(http.NewCaller()))
	if err != nil {
		t.Fatal(err)
	}

	if options.Callers == nil {
		t.Fatal("callers not set")
	}

	caller := options.Callers.Get("http")
	if caller == nil {
		t.Fatal("HTTP caller does not exist")
	}
}

func TestWithCodec(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	options, err := NewOptions(ctx, WithCodec(json.NewConstructor()))
	if err != nil {
		t.Fatal(err)
	}

	if options.Codec == nil {
		t.Fatal("codecs not set")
	}

	codec := options.Codec.Get("json")
	if codec == nil {
		t.Fatal("JSON codec does not exist")
	}
}

func TestWithInvalidLogLevel(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, WithLogLevel("*", "unkown"))
	if err != nil {
		t.Fatal(err)
	}

	if ctx.Atom.Level() != zapcore.ErrorLevel {
		t.Fatalf("unexpected atom level %s, expected %s", ctx.Atom.Level(), zapcore.ErrorLevel)
	}
}

func TestWithInvalidLogLevelPattern(t *testing.T) {
	ctx := logger.WithLogger(broker.WithModule(broker.NewBackground(), "x"))
	_, err := NewOptions(ctx, WithLogLevel("[x-]", "error"))
	if err != nil {
		t.Fatal("unexpected pass")
	}

	if ctx.Atom.Level() != zapcore.ErrorLevel {
		t.Fatalf("unexpected atom level %s, expected %s", ctx.Atom.Level(), zapcore.ErrorLevel)
	}
}

func TestNewOptions(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewOptionsNil(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewOptionsMiddleware(t *testing.T) {
	middleware := MiddlewareFunc(func(*broker.Context) ([]Option, error) {
		options := []Option{
			WithLogLevel("*", "warn"),
		}

		return options, nil
	})

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, WithMiddleware(middleware))
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewOptionsMiddlewareErr(t *testing.T) {
	expected := errors.New("unexpected err")
	middleware := MiddlewareFunc(func(*broker.Context) ([]Option, error) {
		return nil, expected
	})

	ctx := logger.WithLogger(broker.NewBackground())
	_, err := NewOptions(ctx, WithMiddleware(middleware))
	if err != expected {
		t.Fatalf("unexpected err (%+v), expected (%+v)", err, expected)
	}
}
