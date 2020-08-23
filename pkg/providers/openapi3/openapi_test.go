package openapi3

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jexia/semaphore"
	"github.com/jexia/semaphore/cmd/semaphore/daemon"
	"github.com/jexia/semaphore/pkg/broker"
	"github.com/jexia/semaphore/pkg/broker/logger"
	"github.com/jexia/semaphore/pkg/broker/providers"
	"github.com/jexia/semaphore/pkg/codec/json"
	"github.com/jexia/semaphore/pkg/functions"
	provider "github.com/jexia/semaphore/pkg/providers"
	"github.com/jexia/semaphore/pkg/providers/hcl"
	"github.com/jexia/semaphore/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/pkg/transport/grpc"
	"github.com/jexia/semaphore/pkg/transport/http"
	"gopkg.in/yaml.v2"
)

func TestOpenAPI3Generation(t *testing.T) {
	t.Parallel()

	path, err := filepath.Abs("./tests/*.hcl")
	if err != nil {
		t.Fatal(err)
	}

	ctx := logger.WithLogger(broker.NewBackground())
	files, err := provider.ResolvePath(ctx, []string{}, path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			ctx := logger.WithLogger(broker.NewBackground())

			options, err := hcl.GetOptions(ctx, file.Path)
			if err != nil {
				t.Fatal(err)
			}

			arguments := []semaphore.Option{
				semaphore.WithFlows(hcl.FlowsResolver(file.Path)),
				semaphore.WithServices(hcl.ServicesResolver(file.Path)),
				semaphore.WithEndpoints(hcl.EndpointsResolver(file.Path)),
				semaphore.WithCodec(json.NewConstructor()),
				semaphore.WithListener(http.NewListener(":0")),
				semaphore.WithListener(grpc.NewListener(":0", nil)),
				semaphore.WithCaller(http.NewCaller()),
			}

			for _, proto := range options.Protobuffers {
				arguments = append(arguments, semaphore.WithSchema(protobuffers.SchemaResolver([]string{"./tests"}, proto)))
			}

			client, err := daemon.New(ctx, arguments...)

			if err != nil {
				t.Fatal(err)
			}

			collection, err := providers.Resolve(ctx, functions.Collection{}, client.Options)
			if err != nil {
				t.Fatal(err)
			}

			result, err := Generate(collection.EndpointList, collection.FlowListInterface)
			if err != nil {
				t.Fatal(err)
			}

			name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			target := name + ".yaml"

			path, err := filepath.Abs("./tests/" + target)
			if err != nil {
				t.Fatal(err)
			}

			bb, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			expected := Object{}

			err = yaml.Unmarshal(bb, &expected)
			if err != nil {
				t.Fatal(err)
			}

			if reflect.DeepEqual(result, expected) {
				t.Fatal("unexpected result")
			}
		})
	}
}
