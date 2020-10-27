package specs

import (
	"github.com/jexia/semaphore/pkg/discovery"
)

type ServiceDiscoveryClient interface {
	Resolver(host string) (discovery.Resolver, error)
	Provider() string
}

type ServiceDiscoveryClients map[string]ServiceDiscoveryClient
