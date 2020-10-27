package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/jexia/semaphore/pkg/discovery"
	"net/url"
	"sync"
)

type Watcher struct {
	*sync.Mutex
	plan          *watch.Plan
	address       string
	services      chan []discovery.Service
	defaultScheme string
	cache         []discovery.Service

	runCh  chan struct{}
	stopCh chan struct{}

	running bool
}

func newWatcher(address string, updates chan []discovery.Service, plan *watch.Plan) *Watcher {
	watcher := &Watcher{
		Mutex:    &sync.Mutex{},
		plan:     plan,
		address:  address,
		services: updates,
		runCh:    make(chan struct{}),
		stopCh:   make(chan struct{}),
	}

	return watcher
}

func (w *Watcher) Run() {
	w.Lock()
	if w.running {
		return
	}
	w.Unlock()

	go w.plan.Run(w.address)

	go func() {
		for {
			select {
			case <-w.stopCh:
				w.plan.Stop()
				w.Lock()
				w.running = false
				w.Unlock()

			case services := <-w.services:
				w.Lock()
				w.cache = services
				w.Unlock()
			}
		}
	}()

	w.Lock()
	w.running = true
	w.Unlock()
}

func (w *Watcher) Resolve() (string, bool) {
	w.Lock()
	defer w.Unlock()

	if len(w.cache) == 0 {
		return "", false
	}

	svc := w.cache[0]

	uri := url.URL{
		Scheme: svc.Scheme,
		Host:   fmt.Sprintf("%s:%d", svc.Host, svc.Port),
	}

	return uri.String(), true
}

func (w *Watcher) Stop() {
	w.Lock()
	defer w.Unlock()

	if !w.running {
		return
	}

	w.stopCh <- struct{}{}
}

func NewWatcherPlan(service string, params map[string]interface{}, updates chan []discovery.Service) (*watch.Plan, error) {
	query := map[string]interface{}{
		"service":     service,
		"type":        "service",
		"passingonly": true,
	}

	for k, v := range params {
		query[k] = v
	}

	plan, err := watch.Parse(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	plan.Handler = func(idx uint64, data interface{}) {
		var services []discovery.Service

		servicesList, ok := data.([]*api.ServiceEntry)
		if !ok {
			return
		}

		for _, serviceEntry := range servicesList {
			svc := serviceEntry.Service

			host := svc.Address
			if host == "" {
				host = serviceEntry.Node.Address
			}

			services = append(services, discovery.Service{
				Host: host,
				Port: svc.Port,
				Name: svc.Service,
				ID:   svc.ID,
			})
		}

		updates <- services
	}

	return plan, nil
}
