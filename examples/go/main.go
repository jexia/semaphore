package main

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/pkg/codec/json"
	"github.com/jexia/maestro/pkg/codec/proto"
	"github.com/jexia/maestro/pkg/definitions/hcl"
	"github.com/jexia/maestro/pkg/definitions/protoc"
	"github.com/jexia/maestro/pkg/logger"
	"github.com/jexia/maestro/pkg/specs"
	"github.com/jexia/maestro/pkg/transport/http"
)

func main() {
	client, err := maestro.New(
		maestro.WithLogLevel(logger.Global, "debug"),
		maestro.WithListener(http.NewListener(":8080", specs.Options{})),
		maestro.WithFlows(hcl.FlowsResolver("./*.hcl")),
		maestro.WithEndpoints(hcl.EndpointsResolver("./*.hcl")),
		maestro.WithSchema(protoc.SchemaResolver([]string{"../../../", "./proto"}, "./proto/*.proto")),
		maestro.WithServices(protoc.ServiceResolver([]string{"../../../", "./proto"}, "./proto/*.proto")),
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
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
