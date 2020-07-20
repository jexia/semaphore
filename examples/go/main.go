package main

import (
	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/codec/proto"
	"github.com/jexia/semaphore/pkg/core/logger"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/specs"
	"github.com/jexia/semaphore/pkg/transport/http"
)

func main() {
	client, err := semaphore.New(
		semaphore.WithLogLevel(logger.Global, "debug"),
		semaphore.WithListener(http.NewListener(":8080", specs.Options{})),
		semaphore.WithFlows(hcl.FlowsResolver("./*.hcl")),
		semaphore.WithEndpoints(hcl.EndpointsResolver("./*.hcl")),
		semaphore.WithSchema(protobuffers.SchemaResolver([]string{"../../", "./proto"}, "./proto/*.proto")),
		semaphore.WithServices(protobuffers.ServiceResolver([]string{"../../", "./proto"}, "./proto/*.proto")),
		semaphore.WithCodec(json.NewConstructor()),
		semaphore.WithCodec(proto.NewConstructor()),
		semaphore.WithCaller(http.NewCaller()),
	)

	if err != nil {
		panic(err)
	}

	err = client.Serve()
	if err != nil {
		panic(err)
	}
}
