package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/definitions/hcl"
	"github.com/jexia/maestro/protocol/graphql"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Available execution flags
var (
	HTTPAddr    string
	GraphQLAddr string
	ProtoPaths  []string
	FlowPaths   []string
	LogLevel    string
)

// Cmd represents the maestro run command
var Cmd = &cobra.Command{
	Use:   "run",
	Short: "Run the flow definitions with the configured schema format(s)",
	RunE:  run,
}

func init() {
	Cmd.PersistentFlags().StringVar(&HTTPAddr, "http", "", "If set starts the HTTP listener on the given TCP address")
	Cmd.PersistentFlags().StringVar(&GraphQLAddr, "graphql", "", "If set starts the GraphQL listener on the given TCP address")
	Cmd.PersistentFlags().StringSliceVar(&ProtoPaths, "proto", []string{}, "If set are all proto definitions found inside the given path passed as schema definitions, all proto definitions are also passed as imports")
	Cmd.PersistentFlags().StringSliceVar(&FlowPaths, "flow", []string{}, "If set are all flow definitions inside the given path passed as flow definitions")
	Cmd.PersistentFlags().StringVar(&LogLevel, "level", "info", "Logging level")
}

func run(cmd *cobra.Command, args []string) error {
	level, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)

	options := []maestro.Option{
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(http.NewCaller()),
	}

	for _, flow := range FlowPaths {
		options = append(options, maestro.WithDefinitions(hcl.DefinitionResolver(flow)))
	}

	for _, path := range ProtoPaths {
		resolver, err := protoc.Collect(ProtoPaths, path)
		if err != nil {
			return err
		}

		options = append(options, maestro.WithSchema(resolver))
	}

	if HTTPAddr != "" {
		options = append(options, maestro.WithListener(http.NewListener(HTTPAddr, specs.Options{})))
	}

	if GraphQLAddr != "" {
		options = append(options, maestro.WithListener(graphql.NewListener(GraphQLAddr, specs.Options{})))
	}

	client, err := maestro.New(options...)
	if err != nil {
		return err
	}

	go sigterm(client)

	err = client.Serve()
	if err != nil {
		return err
	}

	return nil
}

func sigterm(client *maestro.Client) {
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	<-term
	client.Close()
}
