package flow

import (
	"context"
	"sync"
	"testing"

	"github.com/jexia/maestro/pkg/codec/json"
	"github.com/jexia/maestro/pkg/functions"
	"github.com/jexia/maestro/pkg/instance"
	"github.com/jexia/maestro/pkg/refs"
	"github.com/jexia/maestro/pkg/specs"
)

type fncounter struct {
	counter int
	mutex   sync.Mutex
	err     error
}

func (counter *fncounter) handle(refs.Store) error {
	counter.mutex.Lock()
	counter.counter++
	counter.mutex.Unlock()
	return counter.err
}

func TestNewRequest(t *testing.T) {
	request := NewRequest(nil, nil, nil)
	if request == nil {
		t.Fatal("unexpected result, expected a request to be returned")
	}
}

func TestNewCall(t *testing.T) {
	ctx := instance.NewContext()
	node := &specs.Node{}
	options := &CallOptions{}

	result := NewCall(ctx, node, options)
	if result == nil {
		t.Fatal("unexpected result, expected a call to be constructed")
	}
}

func TestCallExecution(t *testing.T) {
	type test struct {
		node    *specs.Node
		options *CallOptions
		store   refs.Store
	}

	constructor := json.NewConstructor()
	codec, _ := constructor.New("", &specs.ParameterMap{})

	tests := map[string]*test{
		"empty": {
			node: &specs.Node{},
			options: &CallOptions{
				Request:  &Request{},
				Response: &Request{},
			},
		},
		"request codec": {
			node: &specs.Node{},
			options: &CallOptions{
				Request: &Request{
					codec: codec,
				},
				Response: &Request{},
			},
		},
		"response codec": {
			node: &specs.Node{},
			options: &CallOptions{
				Request: &Request{},
				Response: &Request{
					codec: codec,
				},
			},
		},
		"request functions": {
			node: &specs.Node{},
			options: &CallOptions{
				Request: &Request{
					functions: functions.Stack{
						"sample": &functions.Function{
							Fn: func(store refs.Store) error { return nil },
						},
					},
				},
				Response: &Request{},
			},
		},
		"response functions": {
			node: &specs.Node{},
			options: &CallOptions{
				Request: &Request{},
				Response: &Request{
					functions: functions.Stack{
						"sample": &functions.Function{
							Fn: func(store refs.Store) error { return nil },
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()

			result := NewCall(ctx, test.node, test.options)
			if result == nil {
				t.Fatal("unexpected result, expected a call to be constructed")
			}

			err := result.Do(context.Background(), test.store)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestCallFunctionsExecution(t *testing.T) {
	type test struct {
		fn       *fncounter
		expected int
		node     *specs.Node
		options  *CallOptions
		store    refs.Store
	}

	tests := map[string]func() *test{
		"request": func() *test {
			counter := &fncounter{}

			return &test{
				fn:       counter,
				expected: 1,
				node:     &specs.Node{},
				options: &CallOptions{
					Request: &Request{
						functions: functions.Stack{
							"sample": &functions.Function{
								Fn: counter.handle,
							},
						},
					},
					Response: &Request{},
				},
			}
		},
		"response": func() *test {
			counter := &fncounter{}

			return &test{
				fn:       counter,
				expected: 1,
				node:     &specs.Node{},
				options: &CallOptions{
					Request: &Request{},
					Response: &Request{
						functions: functions.Stack{
							"sample": &functions.Function{
								Fn: counter.handle,
							},
						},
					},
				},
			}
		},
		"multiple": func() *test {
			counter := &fncounter{}

			return &test{
				fn:       counter,
				expected: 4,
				node:     &specs.Node{},
				options: &CallOptions{
					Request: &Request{
						functions: functions.Stack{
							"first": &functions.Function{
								Fn: counter.handle,
							},
							"second": &functions.Function{
								Fn: counter.handle,
							},
						},
					},
					Response: &Request{
						functions: functions.Stack{
							"first": &functions.Function{
								Fn: counter.handle,
							},
							"second": &functions.Function{
								Fn: counter.handle,
							},
						},
					},
				},
			}
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := instance.NewContext()
			options := test()

			result := NewCall(ctx, options.node, options.options)
			if result == nil {
				t.Fatal("unexpected result, expected a call to be constructed")
			}

			err := result.Do(context.Background(), options.store)
			if err != nil {
				t.Fatal(err)
			}

			if options.fn.counter != options.expected {
				t.Fatalf("unexpected result %d, expected %d functions to be called", options.fn.counter, options.expected)
			}
		})
	}
}
