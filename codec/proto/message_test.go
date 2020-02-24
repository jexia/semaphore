package proto

import (
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/refs"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jhump/protoreflect/dynamic"
)

func TestMarshalMessage(t *testing.T) {
	collection, err := protoc.Collect(nil, "./tests")
	if err != nil {
		t.Fatal(err)
	}

	manifest, err := maestro.New(maestro.WithPath("./tests", false), maestro.WithSchemaCollection(collection))
	if err != nil {
		t.Fatal(err)
	}

	method := collection.GetService("proto.Logger").GetMethod("Append")
	schema := method.GetInput().(protoc.Object)
	message := dynamic.NewMessage(schema.GetDescriptor())
	specs := manifest.Flows[0].GetCalls()[0].Request
	store := refs.NewStore(1)

	err = MarshalMessage(message, schema, specs, store)
	if err != nil {
		t.Fatal(err)
	}

	json, _ := message.MarshalJSONPB(&jsonpb.Marshaler{
		EmitDefaults: true,
	})

	t.Log(string(json))
}
