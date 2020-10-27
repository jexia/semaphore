package consul

import (
	"fmt"
	"github.com/jexia/semaphore/pkg/discovery"
)

type Consul struct {
	address  string
	watchers map[string]*Watcher
}

func New(address string) *Consul {
	return &Consul{
		address: address,
	}
}

// Resolver returns a service resolver based on the continuous watcher (*Watcher type), that
// is subscribed to all the changes related to the service name.
// The (*Watcher).Resolve() is able to return new service address if the address has been changed.
func (w *Consul) Resolver(name string) (discovery.Resolver, error) {
	if watcher, ok := w.watchers[name]; ok {
		return watcher, nil
	}

	ch := make(chan []discovery.Service)
	plan, err := NewWatcherPlan(name, nil, ch)
	if err != nil {
		return nil, fmt.Errorf("failed to build a watching plan: %w", err)
	}

	watcher := newWatcher(w.address, ch, plan)
	watcher.Run()

	return watcher, nil
}

func (w *Consul) Provider() string {
	return "consul"
}
