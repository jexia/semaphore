package openapi3

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/jexia/semaphore/v2"
	"github.com/jexia/semaphore/v2/cmd/semaphore/daemon/providers"
	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/jexia/semaphore/v2/pkg/codec/json"
	"github.com/jexia/semaphore/v2/pkg/functions"
	provider "github.com/jexia/semaphore/v2/pkg/providers"
	"github.com/jexia/semaphore/v2/pkg/providers/hcl"
	"github.com/jexia/semaphore/v2/pkg/providers/protobuffers"
	"github.com/jexia/semaphore/v2/pkg/transport/grpc"
	"github.com/jexia/semaphore/v2/pkg/transport/http"
	"gopkg.in/yaml.v2"
)

func TestOpenAPI3GenerationWithoutReference(t *testing.T) {
	t.Parallel()

	path, err := filepath.Abs("./tests/*.hcl")
	if err != nil {
		t.Fatal(err)
	}

	files, err := provider.ResolvePath(logger.WithLogger(broker.NewBackground()), []string{}, path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		file := file

		t.Run(file.Name(), func(t *testing.T) {
			t.Parallel()

			ctx := logger.WithLogger(broker.NewBackground())

			options, err := hcl.GetOptions(ctx, file.Path)
			if err != nil {
				t.Fatal(err)
			}

			core, err := semaphore.NewOptions(ctx,
				semaphore.WithFlows(hcl.FlowsResolver(file.Path)),
				semaphore.WithCodec(json.NewConstructor()),
				semaphore.WithCaller(http.NewCaller()),
			)
			if err != nil {
				t.Fatal(err)
			}

			arguments := []providers.Option{
				providers.WithServices(hcl.ServicesResolver(file.Path)),
				providers.WithEndpoints(hcl.EndpointsResolver(file.Path)),
				providers.WithListener(http.NewListener(":0")),
				providers.WithListener(grpc.NewListener(":0", nil)),
			}

			for _, proto := range options.Protobuffers {
				arguments = append(arguments, providers.WithSchema(protobuffers.SchemaResolver([]string{"./tests"}, proto)))
			}

			provider, err := providers.NewOptions(ctx, core, arguments...)
			if err != nil {
				t.Fatal(err)
			}

			collection, err := providers.Resolve(ctx, functions.Collection{}, provider)
			if err != nil {
				t.Fatal(err)
			}

			result, err := Generate(collection.EndpointList, collection.FlowListInterface, IncludeNotReferenced)
			if err != nil {
				t.Fatal(err)
			}

			name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			target := name + ".yaml"

			path, err := filepath.Abs("./tests/no_ref/" + target)
			if err != nil {
				t.Fatal(err)
			}

			bb, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			expected := &Object{}
			err = yaml.Unmarshal(bb, expected)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(result, expected); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}

func TestOpenAPI3GenerationWithReference(t *testing.T) {
	t.Parallel()

	path, err := filepath.Abs("./tests/*.hcl")
	if err != nil {
		t.Fatal(err)
	}

	files, err := provider.ResolvePath(logger.WithLogger(broker.NewBackground()), []string{}, path)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		file := file

		t.Run(file.Name(), func(t *testing.T) {
			t.Parallel()

			ctx := logger.WithLogger(broker.NewBackground())

			options, err := hcl.GetOptions(ctx, file.Path)
			if err != nil {
				t.Fatal(err)
			}

			core, err := semaphore.NewOptions(ctx,
				semaphore.WithFlows(hcl.FlowsResolver(file.Path)),
				semaphore.WithCodec(json.NewConstructor()),
				semaphore.WithCaller(http.NewCaller()),
			)
			if err != nil {
				t.Fatal(err)
			}

			arguments := []providers.Option{
				providers.WithServices(hcl.ServicesResolver(file.Path)),
				providers.WithEndpoints(hcl.EndpointsResolver(file.Path)),
				providers.WithListener(http.NewListener(":0")),
				providers.WithListener(grpc.NewListener(":0", nil)),
			}

			for _, proto := range options.Protobuffers {
				arguments = append(arguments, providers.WithSchema(protobuffers.SchemaResolver([]string{"./tests"}, proto)))
			}

			provider, err := providers.NewOptions(ctx, core, arguments...)
			if err != nil {
				t.Fatal(err)
			}

			collection, err := providers.Resolve(ctx, functions.Collection{}, provider)
			if err != nil {
				t.Fatal(err)
			}

			result, err := Generate(collection.EndpointList, collection.FlowListInterface, DefaultOption)
			if err != nil {
				t.Fatal(err)
			}

			name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			target := name + ".yaml"

			path, err := filepath.Abs("./tests/ref/" + target)
			if err != nil {
				t.Fatal(err)
			}

			bb, err := ioutil.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			expected := &Object{}
			err = yaml.Unmarshal(bb, expected)
			if err != nil {
				t.Fatal(err)
			}

			if diff := deep.Equal(result, expected); diff != nil {
				t.Fatal(diff)
			}
		})
	}
}
