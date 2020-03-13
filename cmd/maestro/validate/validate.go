package validate

import (
	"github.com/spf13/cobra"
)

// Available execution flags
var (
	RecursiveLookup bool
	ProtoPath       string
	ProtoImports    []string
	LogLevel        string
)

// Cmd represents the maestro validate command
var Cmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate the flow definitions with the configured schema format(s)",
	Args:  cobra.MinimumNArgs(1),
	RunE:  run,
}

func init() {
	Cmd.PersistentFlags().StringVar(&ProtoPath, "proto", "", "If set are all proto definitions found inside the given path passed as schema definitions")
	Cmd.PersistentFlags().StringVar(&LogLevel, "level", "info", "Logging level")
}

func run(cmd *cobra.Command, args []string) error {
	// flows := args[0]

	// options := []maestro.Option{
	// 	maestro.WithPath(flows, RecursiveLookup),
	// 	maestro.WithCodec(json.NewConstructor()),
	// 	maestro.WithCodec(proto.NewConstructor()),
	// 	maestro.WithCaller(http.NewCaller()),
	// 	maestro.WithSchema(collection),
	// }
	return nil
}
