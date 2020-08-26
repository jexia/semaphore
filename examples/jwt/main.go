package main

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/daemon"
	"github.com/jexia/semaphore/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/codec/proto"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/functions/lib/jwt"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/transport/http"
)

// Claims is a custom implementation of lib/jwt.Claims interface.
type Claims struct{ jwtgo.StandardClaims }

// Subject is a method returning subject (to satisfy Claims interface).
func (c Claims) Subject() string { return c.StandardClaims.Subject }

// NewClaims instantiates a new claims object.
func NewClaims() jwt.Claims { return new(Claims) }

func main() {
	var (
		reader    = jwt.HMAC("very-strong-secret")
		ctx       = logger.WithLogger(broker.NewContext())
		functions = functions.Custom{
			"jwt": jwt.New(reader, NewClaims),
		}
	)

	core, err := semaphore.NewOptions(ctx,
		semaphore.WithLogLevel("*", "debug"),
		semaphore.WithFlows(hcl.FlowsResolver("./*.hcl")),
		semaphore.WithCodec(json.NewConstructor()),
		semaphore.WithCodec(proto.NewConstructor()),
		semaphore.WithCaller(http.NewCaller()),
		semaphore.WithFunctions(functions),
	)

	if err != nil {
		panic(err)
	}

	options, err := providers.NewOptions(ctx, core,
		providers.WithEndpoints(hcl.EndpointsResolver("./*.hcl")),
		providers.WithSchema(protobuffers.SchemaResolver([]string{"./proto"}, "./proto/*.proto")),
		providers.WithServices(protobuffers.ServiceResolver([]string{"./proto"}, "./proto/*.proto")),
		providers.WithListener(http.NewListener(":8080")),
	)

	if err != nil {
		panic(err)
	}

	client, err := daemon.NewClient(ctx, core, options)
	if err != nil {
		panic(err)
	}

	if err := client.Serve(); err != nil {
		panic(err)
	}
}
