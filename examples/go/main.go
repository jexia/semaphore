package main

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/pkg/codec"
	"github.com/jexia/maestro/pkg/core/logger"
	"github.com/jexia/maestro/pkg/definitions/hcl"
	"github.com/jexia/maestro/pkg/definitions/proto"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport/http"
)

func main() {
	client, err := maestro.New(
		maestro.WithLogLevel(logger.Global, "debug"),
		maestro.WithListener(http.NewListener(":8080", specs.Options{})),
		maestro.WithFlows(hcl.FlowsResolver("./*.hcl")),
		maestro.WithEndpoints(hcl.EndpointsResolver("./*.hcl")),
		maestro.WithSchema(proto.SchemaResolver([]string{"../../", "./proto"}, "./proto/*.proto")),
		maestro.WithServices(proto.ServiceResolver([]string{"../../", "./proto"}, "./proto/*.proto")),
		maestro.WithCodec(codec.JSON()),
		maestro.WithCodec(codec.Proto()),
		maestro.WithCaller(http.NewCaller()),
	)

	if err != nil {
		panic(err)
	}

	err = client.Serve()
	if err != nil {
		panic(err)
	}
}
