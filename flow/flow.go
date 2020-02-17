package flow

import (
	"context"
	"io"
	"sync"

	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/specs"
)

// Error represents a call error
type Error struct {
	ID     string
	Err    error
	Status int
}

// Errors represents a error stack
type Errors []error

// HasErr checks whether there are any errors inside the given error stack
func (errors Errors) HasErr() bool {
	for _, err := range errors {
		if err != nil {
			return true
		}
	}

	return false
}

// Error returns the error stack seperated by a comma
func (errors Errors) Error() (result string) {
	if len(errors) == 0 {
		return result
	}

	result = errors[0].Error()

	for _, err := range errors[1:] {
		result += "," + err.Error()
	}

	return result
}

// Call represents a caller which could be called
type Call func(context.Context, io.Reader) (io.Reader, error)

// Caller represents a caller and the request/response codec parser
type Caller struct {
	Call  Call
	Codec Codec
}

// Step represents a collection of callers and rollbacks which could be executed parallel.
type Step struct {
	callers   []Caller
	rollbacks []Caller
}

// Call executes all the callers inside the given step.
// If any error occures during execution is it appended to the error stack.
func (step *Step) Call(ctx context.Context, refs *refs.Store) Errors {
	return step.Do(ctx, refs, step.callers)
}

// Rollback executes all the rollbacks inside the given step.
// If any error occures during execution is it appended to the error stack.
func (step *Step) Rollback(ctx context.Context, refs *refs.Store) Errors {
	return step.Do(ctx, refs, step.rollbacks)
}

// Do executes the given callers in parallel and unmarshals the response to the given reference store.
// If any error occurs during execution is the returned error appended to the error stack.
func (step *Step) Do(ctx context.Context, refs *refs.Store, callers []Caller) Errors {
	wg := sync.WaitGroup{}
	wg.Add(len(callers))

	errs := make(Errors, len(callers))
	store := sync.Mutex{}

	for index, caller := range callers {
		go func(index int, caller Caller) {
			reader, err := caller.Codec.Marshal(refs)
			if err != nil {
				errs[index] = err
				return
			}

			reader, err = caller.Call(ctx, reader)
			if err != nil {
				errs[index] = err
				return
			}

			store.Lock()
			err = caller.Codec.Unmarshal(reader, refs)
			store.Unlock()
			if err != nil {
				errs[index] = err
				return
			}
		}(index, caller)
	}

	wg.Wait()
	return errs
}

// Manager is responsible for the handling of a flow and it's setps
type Manager struct {
	Refs  *refs.Store
	Flow  *specs.Flow
	Steps []*Step
}

// Call calls all the steps inside the manager if a error is returned is a rollback of all the already executed steps triggered
func (manager *Manager) Call(ctx context.Context, store *refs.Store) error {
	for index, step := range manager.Steps {
		errs := step.Call(ctx, store)
		if errs.HasErr() {
			go manager.Rollback(context.Background(), store, index)
			return errs
		}
	}

	return nil
}

// Rollback calls all the rollbacks in a chronological order from the given position.
// Any errors returned while executing the rollback are logged.
func (manager *Manager) Rollback(ctx context.Context, store *refs.Store, pos int) {
	for index := pos; index >= 0; index-- {
		step := manager.Steps[index]
		errs := step.Call(ctx, store)
		if errs.HasErr() {
			// TODO: log error
		}
	}
}
