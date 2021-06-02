package hcl

import (
	"fmt"

	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/discovery/consul"
	"github.com/jexia/semaphore/v2/pkg/specs"
)

func ParseDiscoveryClients(ctx *broker.Context, manifest Manifest) (specs.ServiceDiscoveryClients, error) {
	configs := manifest.DiscoveryServers
	if configs == nil {
		return nil, nil
	}

	clients := specs.ServiceDiscoveryClients{}

	for _, config := range configs {
		// Use discovery name as the provider name if the provider is not set.
		if config.Provider == "" {
			config.Provider = config.Name
		}

		client, err := newServiceDiscoveryClient(config)
		if err != nil {
			return nil, fmt.Errorf("failed to setup service discovery client for '%s': %w", config.Name, err)
		}

		clients[config.Name] = client
	}

	return clients, nil
}

// initialize a resolver for the specific service.
func newServiceDiscoveryClient(dsc Discovery) (specs.ServiceDiscoveryClient, error) {
	switch dsc.Provider {
	case "consul":
		return consul.New(), nil

	default:
		return nil, fmt.Errorf("unknown provider %q", dsc.Provider)
	}
}
