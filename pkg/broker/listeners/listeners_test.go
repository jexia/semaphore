package listeners

import (
	"context"
	"errors"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/transport"
)

type listener struct {
	name string
	err  error
}

func (listener *listener) Name() string { return listener.name }
func (listener *listener) Serve() error { return listener.err }
func (listener *listener) Close() error { return listener.err }
func (listener *listener) Handle(*broker.Context, []*transport.Endpoint, map[string]codec.Constructor) error {
	return listener.err
}

type flow struct {
	name string
	err  error
}

func (flow *flow) NewStore() references.Store                          { return references.NewReferenceStore(0) }
func (flow *flow) GetName() string                                     { return flow.name }
func (flow *flow) Errors() []transport.Error                           { return nil }
func (flow *flow) Do(ctx context.Context, refs references.Store) error { return flow.err }
func (flow *flow) Wait()                                               {}

func TestApplyNil(t *testing.T) {
	err := Apply(nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplyEmpty(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	err := Apply(ctx, codec.Constructors{}, transport.ListenerList{nil}, transport.EndpointList{nil})
	if err != nil {
		t.Fatal(err)
	}
}

func TestApply(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())

	constructors := codec.Constructors{}

	endpoints := transport.EndpointList{
		transport.NewEndpoint("http", &flow{}, nil, nil, nil, nil),
	}

	listeners := transport.ListenerList{
		&listener{name: "http"},
	}

	err := Apply(ctx, constructors, listeners, endpoints)
	if err != nil {
		t.Fatal(err)
	}
}

func TestApplyUnknownListener(t *testing.T) {
	expected := "unknown listener 'http'"

	ctx := logger.WithLogger(broker.NewBackground())

	constructors := codec.Constructors{}

	endpoints := transport.EndpointList{
		transport.NewEndpoint("http", &flow{}, nil, nil, nil, nil),
	}

	listeners := transport.ListenerList{}

	err := Apply(ctx, constructors, listeners, endpoints)
	if err.Error() != expected {
		t.Fatalf("unexpected err %s, expected %s", err.Error(), expected)
	}
}

func TestApplyListenerErr(t *testing.T) {
	expected := errors.New("unexpected err")

	ctx := logger.WithLogger(broker.NewBackground())

	constructors := codec.Constructors{}

	endpoints := transport.EndpointList{
		transport.NewEndpoint("http", &flow{}, nil, nil, nil, nil),
	}

	listeners := transport.ListenerList{
		&listener{name: "http", err: expected},
	}

	err := Apply(ctx, constructors, listeners, endpoints)
	if err != expected {
		t.Fatalf("unexpected err %s, expected %s", err.Error(), expected.Error())
	}
}
