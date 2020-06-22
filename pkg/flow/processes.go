package flow

import (
	"sync"

	"github.com/jexia/maestro/pkg/transport"
)

// NewProcesses constructs a new processes tracker.
// The given delta will be added to the wait group counter.
func NewProcesses(delta int) *Processes {
	processes := &Processes{}
	processes.Add(delta)

	return processes
}

// Processes tracks processes
type Processes struct {
	err   transport.Error
	wg    sync.WaitGroup
	mutex sync.Mutex
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
func (processes *Processes) Add(delta int) {
	processes.wg.Add(delta)
}

// Done marks a given process as done
func (processes *Processes) Done() {
	processes.wg.Done()
}

// Wait awaits till all processes are completed
func (processes *Processes) Wait() {
	processes.wg.Wait()
}

// Err returns the thrown error if thrown
func (processes *Processes) Err() transport.Error {
	processes.mutex.Lock()
	defer processes.mutex.Unlock()
	return processes.err
}

// Fatal marks the given error and is returned on Err()
func (processes *Processes) Fatal(err transport.Error) {
	if err == nil {
		return
	}

	processes.mutex.Lock()
	defer processes.mutex.Unlock()

	if processes.err != nil {
		return
	}

	processes.err = err
}
