package consul

import (
	"fmt"
	"net/url"

	"github.com/jexia/semaphore/v2/pkg/discovery"
)

type Consul struct {
	watchers map[string]*Watcher
}

// New creates a new Consul manager that keeps track of the Consul resolvers.
func New() *Consul {
	return &Consul{
		watchers: make(map[string]*Watcher),
	}
}

// Resolver returns a service resolver based on the continuous watcher (*Watcher type), that
// is subscribed to all the changes related to the service name.
// The (*Watcher).Resolve() is able to return new service address if the address has been changed.
func (c *Consul) Resolver(address string) (discovery.Resolver, error) {
	uri, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the raw service address: %w", err)
	}

	name := uri.Host
	if watcher, ok := c.watchers[name]; ok {
		return watcher, nil
	}

	// Create a new watcher
	ch := make(chan []discovery.Service)

	plan, err := NewWatcherPlan(name, nil, ch)
	if err != nil {
		return nil, fmt.Errorf("failed to build a watching plan: %w", err)
	}

	watcher := newWatcher(uri.Host, uri.Scheme, ch, plan)
	watcher.Run()

	c.watchers[name] = watcher

	return watcher, nil
}

func (c *Consul) Provider() string {
	return "consul"
}
