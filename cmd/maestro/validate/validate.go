package validate

import (
	"github.com/jexia/maestro"
	"github.com/jexia/maestro/cmd/maestro/config"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var global = config.New()

// Cmd represents the maestro validate command
var Cmd = &cobra.Command{
	Use:          "validate",
	Short:        "Validate the flow definitions with the configured schema format(s)",
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Cmd.PersistentFlags().StringP("config", "c", "", "Config file path")
	Cmd.PersistentFlags().StringSliceVar(&global.Protobuffers, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Cmd.PersistentFlags().StringSliceVar(&global.Flows, "flow", []string{}, "If set are all flow definitions inside the given path passed as flow definitions")
	Cmd.PersistentFlags().StringVar(&global.LogLevel, "level", "error", "Logging level")
}

func run(cmd *cobra.Command, args []string) error {
	err := config.Read(cmd, global)
	if err != nil {
		return err
	}

	options := []maestro.Option{
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(http.NewCaller()),
	}

	level, err := logrus.ParseLevel(global.LogLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)

	for _, flow := range global.Flows {
		options = append(options, maestro.WithDefinitions(hcl.DefinitionResolver(flow)))
	}

	for _, path := range global.Protobuffers {
		resolver, err := protoc.Collect(global.Protobuffers, path)
		if err != nil {
			return err
		}

		options = append(options, maestro.WithSchema(resolver))
	}

	_, err = maestro.ConstructSpecs(maestro.NewOptions(options...))
	if err != nil {
		return err
	}

	return nil
}
