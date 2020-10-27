package e2e

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/daemon"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/cmd/semaphore/middleware"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	codecJSON "github.com/jexia/semaphore/pkg/codec/json"
	codecProto "github.com/jexia/semaphore/pkg/codec/proto"
	codecXML "github.com/jexia/semaphore/pkg/codec/xml"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/specs"
	transportGRPC "github.com/jexia/semaphore/pkg/transport/grpc"
	transportHTTP "github.com/jexia/semaphore/pkg/transport/http"
)

const (
	SemaphoreGRPCPort = 5051
	SemaphoreHTTPPort = 8080
	SemaphoreHost     = "localhost"
)

// Instance creates a new semaphore instance.
func Instance(t *testing.T, flow, schema string) *daemon.Client {
	ctx := logger.WithLogger(broker.NewContext())

	core, err := semaphore.NewOptions(ctx,
		semaphore.WithLogLevel("*", "error"),
		semaphore.WithFlows(hcl.FlowsResolver(flow)),
		semaphore.WithCodec(codecXML.NewConstructor()),
		semaphore.WithCodec(codecJSON.NewConstructor()),
		semaphore.WithCodec(codecProto.NewConstructor()),
		semaphore.WithCaller(transportHTTP.NewCaller()),
		semaphore.WithCaller(transportGRPC.NewCaller()),
	)

	if err != nil {
		t.Fatalf("cannot instantiate semaphore core: %s", err)
	}

	options, err := providers.NewOptions(ctx, core,
		providers.WithEndpoints(hcl.EndpointsResolver(flow)),
		providers.WithSchema(protobuffers.SchemaResolver([]string{filepath.Dir(schema)}, schema)),
		providers.WithServices(protobuffers.ServiceResolver([]string{filepath.Dir(schema)}, schema)),
		providers.WithListener(transportHTTP.NewListener(fmt.Sprintf(":%d", SemaphoreHTTPPort))),
		providers.WithListener(transportGRPC.NewListener(fmt.Sprintf(":%d", SemaphoreGRPCPort), specs.Options{})),
		providers.WithAfterConstructor(middleware.ServiceSelector(flow)),
	)

	if err != nil {
		t.Fatalf("unable to configure provider options: %s", err)
	}

	client, err := daemon.NewClient(ctx, core, options)
	if err != nil {
		t.Fatalf("failed to create a semaphore instance: %s", err)
	}

	return client
}
