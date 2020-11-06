package dependencies

import "fmt"

// ErrCircularDependency is returned when circular dependency is detected.
type ErrCircularDependency struct {
	Flow, From, To string
}

func (e ErrCircularDependency) Error() string {
	return fmt.Sprintf("circular resource dependency detected: %s.%s <-> %s.%s", e.Flow, e.From, e.Flow, e.To)
}
