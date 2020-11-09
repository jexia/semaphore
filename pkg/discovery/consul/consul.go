package consul

import (
	"fmt"
	"github.com/jexia/semaphore/pkg/discovery"
	"net/url"
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
func (c *Consul) Resolver(address string) (discovery.Resolver, error) {
	uri, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the raw service address: %w", err)
	}

	name := uri.Host

	if watcher, ok := c.watchers[name]; ok {
		return watcher, nil
	}

	ch := make(chan []discovery.Service)
	plan, err := NewWatcherPlan(name, nil, ch)
	if err != nil {
		return nil, fmt.Errorf("failed to build a watching plan: %c", err)
	}

	watcher := newWatcher(c.address, uri.Scheme, ch, plan)
	watcher.Run()

	return watcher, nil
}

func (c *Consul) Provider() string {
	return "consul"
}
