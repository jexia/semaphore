package e2e

import (
	"path/filepath"
	"testing"

	"github.com/jexia/semaphore/v2"
	"github.com/jexia/semaphore/v2/cmd/semaphore/daemon"
	"github.com/jexia/semaphore/v2/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/v2/cmd/semaphore/middleware"
	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/providers/hcl"
	"github.com/jexia/semaphore/v2/pkg/providers/protobuffers"
)

// Config contains the settings of semaphore instance.
type Config struct {
	SemaphoreOptions []semaphore.Option
	ProviderOptions  []providers.Option
}

// Instance creates a new semaphore instance.
func Instance(t *testing.T, flow, schema string, config Config) *daemon.Client { // callers
	ctx := logger.WithLogger(broker.NewContext())

	semaphoreOptions := []semaphore.Option{
		semaphore.WithLogLevel("*", "error"),
		semaphore.WithFlows(hcl.FlowsResolver(flow)),
	}

	core, err := semaphore.NewOptions(ctx, append(semaphoreOptions, config.SemaphoreOptions...)...)
	if err != nil {
		t.Fatalf("cannot instantiate semaphore core: %s", err)
	}

	providerOptions := []providers.Option{
		providers.WithEndpoints(hcl.EndpointsResolver(flow)),
		providers.WithSchema(protobuffers.SchemaResolver([]string{filepath.Dir(schema)}, schema)),
		providers.WithServices(protobuffers.ServiceResolver([]string{filepath.Dir(schema)}, schema)),
		providers.WithAfterConstructor(middleware.ServiceSelector(flow)),
	}

	options, err := providers.NewOptions(ctx, core, append(providerOptions, config.ProviderOptions...)...)
	if err != nil {
		t.Fatalf("unable to configure provider options: %s", err)
	}

	instance, err := daemon.NewClient(ctx, core, options)
	if err != nil {
		t.Fatalf("failed to create a semaphore instance: %s", err)
	}

	return instance
}
