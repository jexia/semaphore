package flow

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/conditions"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/specs/labels"
	"github.com/jexia/semaphore/pkg/specs/types"
)

func NewMockNode(name string, caller Call, rollback Call) *Node {
	ctx := logger.WithLogger(broker.NewBackground())

	return &Node{
		ctx:        ctx,
		Name:       name,
		Call:       caller,
		Revert:     rollback,
		OnError:    NewMockOnError(),
		DependsOn:  map[string]*specs.Node{},
		References: map[string]*specs.PropertyReference{},
	}
}

func NewMockOnError() *specs.OnError {
	return &specs.OnError{
		Response: &specs.ParameterMap{
			Property: &specs.Property{
				Label: labels.Optional,
				Template: specs.Template{
					Message: specs.Message{
						"status": {
							Name:  "status",
							Path:  "status",
							Label: labels.Optional,
							Template: specs.Template{
								Scalar: &specs.Scalar{
									Type: types.Int64,
								},
								Reference: &specs.PropertyReference{
									Resource: "error",
									Path:     "status",
								},
							},
						},
						"message": {
							Name:  "message",
							Path:  "message",
							Label: labels.Optional,
							Template: specs.Template{
								Scalar: &specs.Scalar{
									Type: types.String,
								},
								Reference: &specs.PropertyReference{
									Resource: "error",
									Path:     "message",
								},
							},
						},
					},
				},
			},
		},
		Status: &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type:    types.Int64,
					Default: 500,
				},
			},
		},
		Message: &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Scalar: &specs.Scalar{
					Type:    types.String,
					Default: "mock error message",
				},
			},
		},
	}
}

func BenchmarkSingleNodeCallingJSONCodecParallel(b *testing.B) {
	ctx := logger.WithLogger(broker.NewBackground())
	constructor := json.NewConstructor()

	req, err := constructor.New("first.request", &specs.ParameterMap{
		Property: &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Message: specs.Message{
					"key": {
						Name:  "key",
						Path:  "key",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type:    types.String,
								Default: "message",
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		b.Fatal(err)
	}

	res, err := constructor.New("first.response", &specs.ParameterMap{
		Property: &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Message: specs.Message{
					"key": {
						Name:  "key",
						Path:  "key",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type:    types.String,
								Default: "message",
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		b.Fatal(err)
	}

	options := &CallOptions{
		Request: &Request{
			Codec: req,
		},
		Response: &Request{
			Codec: res,
		},
	}

	call := NewCall(ctx, nil, options)
	node := NewMockNode("first", call, nil)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Background()
			tracker := NewTracker("", 1)
			processes := NewProcesses(1)
			refs := references.NewReferenceStore(0)

			node.Do(ctx, tracker, processes, refs)
		}
	})
}

func BenchmarkSingleNodeCallingJSONCodecSerial(b *testing.B) {
	ctx := logger.WithLogger(broker.NewBackground())
	constructor := json.NewConstructor()

	req, err := constructor.New("first.request", &specs.ParameterMap{
		Property: &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Message: specs.Message{
					"key": {
						Name:  "key",
						Path:  "key",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type:    types.String,
								Default: "message",
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		b.Fatal(err)
	}

	res, err := constructor.New("first.response", &specs.ParameterMap{
		Property: &specs.Property{
			Label: labels.Optional,
			Template: specs.Template{
				Message: specs.Message{
					"key": {
						Name:  "key",
						Path:  "key",
						Label: labels.Optional,
						Template: specs.Template{
							Scalar: &specs.Scalar{
								Type:    types.String,
								Default: "message",
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		b.Fatal(err)
	}

	options := &CallOptions{
		Request: &Request{
			Codec: req,
		},
		Response: &Request{
			Codec: res,
		},
	}

	call := NewCall(ctx, nil, options)
	node := NewMockNode("first", call, nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		tracker := NewTracker("", 1)
		processes := NewProcesses(1)
		refs := references.NewReferenceStore(0)

		node.Do(ctx, tracker, processes, refs)
	}
}

func BenchmarkSingleNodeCallingParallel(b *testing.B) {
	caller := &mocker{}
	node := NewMockNode("first", caller, nil)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Background()
			tracker := NewTracker("", 1)
			processes := NewProcesses(1)
			refs := references.NewReferenceStore(0)

			node.Do(ctx, tracker, processes, refs)
		}
	})
}

func BenchmarkSingleNodeCallingSerial(b *testing.B) {
	caller := &mocker{}
	node := NewMockNode("first", caller, nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		tracker := NewTracker("", 1)
		processes := NewProcesses(1)
		refs := references.NewReferenceStore(0)

		node.Do(ctx, tracker, processes, refs)
	}
}

func BenchmarkBranchedNodeCallingParallel(b *testing.B) {
	caller := &mocker{}
	nodes := []*Node{
		NewMockNode("first", caller, nil),
		NewMockNode("second", caller, nil),
		NewMockNode("third", caller, nil),
	}

	nodes[0].Next = []*Node{nodes[1]}
	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[2]}
	nodes[2].Previous = []*Node{nodes[1]}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Background()
			tracker := NewTracker("", len(nodes))
			processes := NewProcesses(1)
			refs := references.NewReferenceStore(0)

			nodes[0].Do(ctx, tracker, processes, refs)
		}
	})
}

func BenchmarkBranchedNodeCallingSerial(b *testing.B) {
	caller := &mocker{}
	nodes := []*Node{
		NewMockNode("first", caller, nil),
		NewMockNode("second", caller, nil),
		NewMockNode("third", caller, nil),
	}

	nodes[0].Next = []*Node{nodes[1]}
	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[2]}
	nodes[2].Previous = []*Node{nodes[1]}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		tracker := NewTracker("", len(nodes))
		processes := NewProcesses(1)
		refs := references.NewReferenceStore(0)

		nodes[0].Do(ctx, tracker, processes, refs)
	}
}

func TestConstructingNode(t *testing.T) {
	type test struct {
		Node     *specs.Node
		Call     Call
		Rollback Call
		Expected int
	}

	tests := map[string]*test{
		"node call": {
			Expected: 2,
			Node: &specs.Node{
				Call: &specs.Call{
					Request: &specs.ParameterMap{
						Property: &specs.Property{
							Template: specs.Template{
								Message: specs.Message{
									"first": {
										Name: "first",
										Path: "first",
										Template: specs.Template{
											Reference: &specs.PropertyReference{
												Resource: "input",
												Path:     "first",
											},
										},
									},
									"second": {
										Name: "second",
										Path: "second",
										Template: specs.Template{
											Reference: &specs.PropertyReference{
												Resource: "input",
												Path:     "second",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"combination": {
			Expected: 1,
			Node: &specs.Node{
				Call: &specs.Call{
					Request: &specs.ParameterMap{
						Property: &specs.Property{
							Template: specs.Template{
								Message: specs.Message{
									"first": {
										Name: "first",
										Path: "first",
										Template: specs.Template{
											Reference: &specs.PropertyReference{
												Resource: "input",
												Path:     "first",
											},
										},
									},
									"second": {
										Name: "second",
										Path: "second",
										Template: specs.Template{
											Reference: &specs.PropertyReference{
												Resource: "input",
												Path:     "first",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())
			result := NewNode(ctx, test.Node, WithCall(test.Call), WithRollback(test.Rollback))

			if len(result.References) != test.Expected {
				t.Fatalf("unexpected amount of references %d, expected %d", len(result.References), test.Expected)
			}
		})
	}
}

func TestConstructingNodeReferences(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	call := &mocker{}
	rollback := &mocker{}

	node := &specs.Node{
		ID: "mock",
	}

	result := NewNode(ctx, node, WithCall(call), WithRollback(rollback))
	if result == nil {
		t.Fatal("nil node returned")
	}
}

func TestNodeHas(t *testing.T) {
	nodes := make(Nodes, 2)

	nodes[0] = &Node{Name: "first"}
	nodes[1] = &Node{Name: "second"}

	if !nodes.Has("first") {
		t.Fatal("unexpected result, expected 'first' to be available")
	}

	if nodes.Has("unexpected") {
		t.Fatal("unexpected result, expected 'unexpected' to be unavailable")
	}
}

func TestNodeCalling(t *testing.T) {
	caller := &mocker{}

	nodes := []*Node{
		NewMockNode("first", caller, nil),
		NewMockNode("second", caller, nil),
		NewMockNode("third", caller, nil),
	}

	nodes[0].Next = []*Node{nodes[1]}
	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[2]}
	nodes[2].Previous = []*Node{nodes[1]}

	tracker := NewTracker("", len(nodes))
	processes := NewProcesses(1)
	refs := references.NewReferenceStore(0)

	nodes[0].Do(context.Background(), tracker, processes, refs)
	processes.Wait()

	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if caller.Counter != len(nodes) {
		t.Errorf("unexpected counter total %d, expected %d", caller.Counter, len(nodes))
	}
}

func TestSlowNodeAbortingOnErr(t *testing.T) {
	slow := &mocker{name: "slow"}
	failed := &mocker{name: "failed", Err: errors.New("unexpected")}
	caller := &mocker{}

	nodes := []*Node{
		NewMockNode("first", caller, nil),
		NewMockNode("second", slow, nil),
		NewMockNode("third", failed, nil),
		NewMockNode("fourth", caller, nil),
	}

	nodes[0].Next = []*Node{nodes[1], nodes[2]}

	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[3]}

	nodes[2].Previous = []*Node{nodes[0]}
	nodes[2].Next = []*Node{nodes[3]}

	nodes[3].Previous = []*Node{nodes[1], nodes[2]}

	tracker := NewTracker("", len(nodes))
	processes := NewProcesses(1)
	refs := references.NewReferenceStore(0)

	slow.mutex.Lock()
	failed.mutex.Lock()

	go func() {
		failed.mutex.Unlock()
		time.Sleep(100 * time.Millisecond)
		slow.mutex.Unlock()
	}()

	nodes[0].Do(context.Background(), tracker, processes, refs)

	processes.Wait()

	counter := (caller.Counter + slow.Counter + failed.Counter)
	if counter != 3 {
		t.Fatalf("unexpected counter total %d, expected %d", counter, 3)
	}
}

func TestNodeRevert(t *testing.T) {
	rollback := &mocker{}

	nodes := []*Node{
		NewMockNode("first", nil, rollback),
		NewMockNode("second", nil, rollback),
		NewMockNode("third", nil, rollback),
	}

	nodes[0].Next = []*Node{nodes[1]}
	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[2]}
	nodes[2].Previous = []*Node{nodes[1]}

	tracker := NewTracker("", len(nodes))
	processes := NewProcesses(1)
	refs := references.NewReferenceStore(0)

	nodes[len(nodes)-1].Rollback(context.Background(), tracker, processes, refs)
	processes.Wait()

	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if rollback.Counter != len(nodes) {
		t.Errorf("unexpected counter total %d, expected %d", rollback.Counter, len(nodes))
	}
}

func TestNodeBranchesCalling(t *testing.T) {
	caller := &mocker{}

	nodes := []*Node{
		NewMockNode("first", caller, nil),
		NewMockNode("second", caller, nil),
		NewMockNode("third", caller, nil),
		NewMockNode("fourth", caller, nil),
	}

	nodes[0].Next = []*Node{nodes[1], nodes[2]}

	nodes[1].Previous = []*Node{nodes[0]}
	nodes[1].Next = []*Node{nodes[3]}
	nodes[2].Previous = []*Node{nodes[0]}
	nodes[2].Next = []*Node{nodes[3]}

	nodes[3].Previous = []*Node{nodes[1], nodes[2]}

	tracker := NewTracker("", len(nodes))
	processes := NewProcesses(1)
	refs := references.NewReferenceStore(0)

	nodes[0].Do(context.Background(), tracker, processes, refs)
	processes.Wait()

	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if caller.Counter != len(nodes) {
		t.Errorf("unexpected counter total %d, expected %d", caller.Counter, len(nodes))
	}
}

func TestBeforeDoNode(t *testing.T) {
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.BeforeDo = func(ctx context.Context, node *Node, tracker Tracker, processes *Processes, store references.Store) (context.Context, error) {
		counter++
		return ctx, nil
	}

	processes := NewProcesses(1)
	node.Do(context.Background(), NewTracker("", 1), processes, nil)
	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after rollback function to be called", counter)
	}
}

func TestBeforeDoNodeErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.BeforeDo = func(ctx context.Context, node *Node, tracker Tracker, processes *Processes, store references.Store) (context.Context, error) {
		counter++
		return ctx, expected
	}

	processes := NewProcesses(1)
	node.Do(context.Background(), NewTracker("", 1), processes, nil)
	if !errors.Is(processes.Err(), expected) {
		t.Errorf("unexpected err '%s', expected '%s' to be thrown", processes.Err(), expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after rollback function to be called", counter)
	}
}

func TestAfterDoNode(t *testing.T) {
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.AfterDo = func(ctx context.Context, node *Node, tracker Tracker, processes *Processes, store references.Store) (context.Context, error) {
		counter++
		return ctx, nil
	}

	processes := NewProcesses(1)
	node.Do(context.Background(), NewTracker("", 1), processes, nil)
	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after rollback function to be called", counter)
	}
}

func TestAfterDoNodeErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.AfterDo = func(ctx context.Context, node *Node, tracker Tracker, processes *Processes, store references.Store) (context.Context, error) {
		counter++
		return ctx, expected
	}

	processes := NewProcesses(1)
	node.Do(context.Background(), NewTracker("", 1), processes, nil)
	if !errors.Is(processes.Err(), expected) {
		t.Errorf("unexpected err '%s', expected '%s' to be thrown", processes.Err(), expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after rollback function to be called", counter)
	}
}

func TestBeforeRevertNode(t *testing.T) {
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.BeforeRollback = func(ctx context.Context, node *Node, tracker Tracker, processes *Processes, store references.Store) (context.Context, error) {
		counter++
		return ctx, nil
	}

	processes := NewProcesses(1)
	node.Rollback(context.Background(), NewTracker("", 1), processes, nil)
	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after revert function to be called", counter)
	}
}

func TestBeforeRevertNodeErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.BeforeRollback = func(ctx context.Context, node *Node, tracker Tracker, processes *Processes, store references.Store) (context.Context, error) {
		counter++
		return ctx, expected
	}

	processes := NewProcesses(1)
	node.Rollback(context.Background(), NewTracker("", 1), processes, nil)
	if !errors.Is(processes.Err(), expected) {
		t.Errorf("unexpected err '%s', expected '%s' to be thrown", processes.Err(), expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after revert function to be called", counter)
	}
}

func TestAfterRevertNode(t *testing.T) {
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.AfterRollback = func(ctx context.Context, node *Node, tracker Tracker, processes *Processes, store references.Store) (context.Context, error) {
		counter++
		return ctx, nil
	}

	processes := NewProcesses(1)
	node.Rollback(context.Background(), NewTracker("", 1), processes, nil)
	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after revert function to be called", counter)
	}
}

func TestAfterRevertNodeErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.AfterRollback = func(ctx context.Context, node *Node, tracker Tracker, processes *Processes, store references.Store) (context.Context, error) {
		counter++
		return ctx, expected
	}

	processes := NewProcesses(1)
	node.Rollback(context.Background(), NewTracker("", 1), processes, nil)
	if !errors.Is(processes.Err(), expected) {
		t.Errorf("unexpected err '%s', expected '%s' to be thrown", processes.Err(), expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected after revert function to be called", counter)
	}
}

func TestNodeDoFunctions(t *testing.T) {
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.Functions = functions.Stack{
		"hash": &functions.Function{
			Fn: func(store references.Store) error {
				counter++
				return nil
			},
		},
	}

	processes := NewProcesses(1)
	node.Do(context.Background(), NewTracker("", 1), processes, nil)
	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected node function to be called", counter)
	}
}

func TestNodeDoFunctionsErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.Functions = functions.Stack{
		"hash": &functions.Function{
			Fn: func(store references.Store) error {
				counter++
				return expected
			},
		},
	}

	processes := NewProcesses(1)
	node.Do(context.Background(), NewTracker("", 1), processes, nil)
	if !errors.Is(processes.Err(), expected) {
		t.Errorf("unexpected err '%s', expected '%s' to be thrown", processes.Err(), expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected node function to be called", counter)
	}
}

func TestNodeDoConditionFunctions(t *testing.T) {
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.Condition = &Condition{
		stack: functions.Stack{
			"hash": &functions.Function{
				Fn: func(store references.Store) error {
					counter++
					return nil
				},
			},
		},
	}

	processes := NewProcesses(1)
	node.Do(context.Background(), NewTracker("", 1), processes, nil)
	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected condition function to be called", counter)
	}
}

func TestNodeDoConditionFunctionsErr(t *testing.T) {
	expected := errors.New("unexpected err")
	counter := 0
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	node.Condition = &Condition{
		stack: functions.Stack{
			"hash": &functions.Function{
				Fn: func(store references.Store) error {
					counter++
					return expected
				},
			},
		},
	}

	processes := NewProcesses(1)
	node.Do(context.Background(), NewTracker("", 1), processes, nil)
	if !errors.Is(processes.Err(), expected) {
		t.Errorf("unexpected err '%s', expected '%s' to be thrown", processes.Err(), expected)
	}

	if counter != 1 {
		t.Fatalf("unexpected counter %d, expected condition functions to be called", counter)
	}
}

func TestNodeDoConditionBreaks(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	call := &mocker{}
	node := NewMockNode("mock", call, nil)

	expression, err := conditions.NewEvaluableExpression(ctx, "false")
	if err != nil {
		t.Fatal(err)
	}

	node.Condition = &Condition{
		expression: expression,
	}

	processes := NewProcesses(1)
	tracker := NewTracker("", 1)

	node.Do(context.Background(), tracker, processes, nil)
	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	if !tracker.Skipped(node) {
		t.Fatal("tracker has not skipped the node")
	}

	if tracker.Met(node) {
		t.Fatal("tracker has met a skipped node")
	}
}

func TestNodeDoSkipChildren(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())
	call := &mocker{}

	expression, err := conditions.NewEvaluableExpression(ctx, "false")
	if err != nil {
		t.Fatal(err)
	}

	condition := &Condition{
		expression: expression,
	}

	first := NewMockNode("first", call, nil)
	second := NewMockNode("second", call, nil)
	third := NewMockNode("third", call, nil)

	first.Next = Nodes{second}
	second.Next = Nodes{third}

	first.Condition = condition
	second.Condition = condition

	processes := NewProcesses(1)
	tracker := NewTracker("", 1)

	first.Do(context.Background(), tracker, processes, nil)
	if processes.Err() != nil {
		t.Error(processes.Err())
	}

	result := Nodes{
		first,
		second,
		third,
	}

	for _, node := range result {
		t.Run(node.Name, func(t *testing.T) {
			if !tracker.Skipped(node) {
				t.Fatal("tracker has not skipped the node")
			}

			if tracker.Met(node) {
				t.Fatal("tracker has met the skipped node")
			}
		})
	}
}

func TestWithCall(t *testing.T) {
	expected := &caller{}
	options := NodeOptions{}

	option := WithCall(expected)

	option(&options)
	if options.call == nil {
		t.Fatal("call not set")
	}

	if options.call != expected {
		t.Fatal("unexpected call")
	}
}

func TestWithRollback(t *testing.T) {
	expected := &caller{}
	options := NodeOptions{}

	option := WithRollback(expected)

	option(&options)
	if options.rollback == nil {
		t.Fatal("rollback not set")
	}

	if options.rollback != expected {
		t.Fatal("unexpected rollback")
	}
}

func TestWithCondition(t *testing.T) {
	expected := &Condition{}
	options := NodeOptions{}

	option := WithCondition(expected)

	option(&options)
	if options.condition == nil {
		t.Fatal("confition not set")
	}

	if options.condition != expected {
		t.Fatal("unexpected condition")
	}
}

func TestWithFunctions(t *testing.T) {
	expected := functions.Stack{}
	options := NodeOptions{}

	option := WithFunctions(expected)

	option(&options)
	if options.functions == nil {
		t.Fatal("functions not set")
	}
}

func TestWithNodeMiddleware(t *testing.T) {
	expected := NodeMiddleware{
		AfterDo: func(ctx context.Context, node *Node, tracker Tracker, processes *Processes, store references.Store) (context.Context, error) {
			return nil, nil
		},
	}

	options := NodeOptions{}

	option := WithNodeMiddleware(expected)

	option(&options)
	if options.middleware.AfterDo == nil {
		t.Fatal("middleware not set")
	}
}

func TestSettingNodeArguments(t *testing.T) {
	arguments := NodeArguments{}
	arguments.Set(nil)

	if len(arguments) != 0 {
		t.Fatal("expected empty arguments")
	}

	arguments.Set(WithCondition(nil))

	if len(arguments) != 1 {
		t.Fatal("arguments not set properly")
	}
}
