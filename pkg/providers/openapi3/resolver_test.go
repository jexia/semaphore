package openapi3

import (
	"encoding/json"
	"testing"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/jexia/semaphore/v2/pkg/broker"
	"github.com/jexia/semaphore/v2/pkg/broker/logger"
	"github.com/stretchr/testify/assert"
)

func Test_getCanonicalName(t *testing.T) {
	type args struct {
		doc  *openapi.Swagger
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"with package property",
			args{
				&openapi.Swagger{
					Info: &openapi.Info{
						ExtensionProps: openapi.ExtensionProps{
							Extensions: map[string]interface{}{
								XPackageExtensionField: json.RawMessage(`"com.semaphore"`),
							},
						},
					},
				},
				"user",
			},
			"com.semaphore.user",
		},
		{
			"no package property",
			args{
				&openapi.Swagger{
					Info: &openapi.Info{
						ExtensionProps: openapi.ExtensionProps{},
					},
				},
				"user",
			},
			"user",
		},
		{
			"no document info",
			args{
				&openapi.Swagger{},
				"user",
			},
			"user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCanonicalName(tt.args.doc, tt.args.name); got != tt.want {
				t.Errorf("getCanonicalName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collect(t *testing.T) {
	ctx := logger.WithLogger(broker.NewBackground())

	t.Run("returns error when cannot resolve imports", func(t *testing.T) {
		_, err := collect(ctx, []string{"./foobar/*.yml"})
		assert.Error(t, err)
	})

	t.Run("returns documents", func(t *testing.T) {
		swaggers, err := collect(ctx, []string{"./fixtures/*.yml"})
		assert.Nil(t, err)

		assert.Len(t, swaggers, 2, "should locate and parse all files")

		_, ok := swaggers["fixtures/petstore.yml"]
		assert.True(t, ok, "should include a document")

		_, ok = swaggers["fixtures/empty.yml"]
		assert.True(t, ok, "should include a document")
	})
}
