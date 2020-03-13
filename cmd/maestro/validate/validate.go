package validate

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/spf13/cobra"
)

// Available execution flags
var (
	ProtoPaths []string
	LogLevel   string
)

// Cmd represents the maestro validate command
var Cmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate the flow definitions with the configured schema format(s)",
	Args:  cobra.MinimumNArgs(1),
	RunE:  run,
}

func init() {
	Cmd.PersistentFlags().StringSliceVar(&ProtoPaths, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Cmd.PersistentFlags().StringVar(&LogLevel, "level", "info", "Logging level")
}

func run(cmd *cobra.Command, args []string) error {
	options := []maestro.Option{
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(http.NewCaller()),
	}

	for _, arg := range args {
		options = append(options, maestro.WithDefinitions(hcl.DefinitionResolver(arg)))
	}

	for _, path := range ProtoPaths {
		resolver, err := protoc.Collect(ProtoPaths, path)
		if err != nil {
			return err
		}

		options = append(options, maestro.WithSchema(resolver))
	}

	_, err := maestro.ConstructSpecs(maestro.NewOptions(options...))
	if err != nil {
		return err
	}

	return nil
}
