package run

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jexia/maestro"
	"github.com/jexia/maestro/codec/json"
	"github.com/jexia/maestro/codec/proto"
	"github.com/jexia/maestro/protocol/http"
	"github.com/jexia/maestro/schema/protoc"
	"github.com/jexia/maestro/specs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Available execution flags
var (
	HTTPAddr     string
	ProtoPath    string
	ProtoImports []string
	LogLevel     string
)

// Cmd represents the maestro run command
var Cmd = &cobra.Command{
	Use:   "run [path]",
	Short: "Run the given flow definitions with the configured schema format",
	Args:  cobra.MinimumNArgs(1),
	RunE:  run,
}

func init() {
	Cmd.PersistentFlags().StringVar(&HTTPAddr, "http", "", "If set starts the HTTP listener on the given TCP address")
	Cmd.PersistentFlags().StringVar(&ProtoPath, "proto", "", "If set are all proto definitions found inside the given path passed as schema definitions")
	Cmd.PersistentFlags().StringSliceVar(&ProtoImports, "proto-import", []string{}, "Proto import definitions")
	Cmd.PersistentFlags().StringVar(&LogLevel, "level", "info", "Logging level")
}

func run(cmd *cobra.Command, args []string) error {
	level, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)

	collection, err := protoc.Collect(append(ProtoImports, ProtoPath), ProtoPath)
	if err != nil {
		return err
	}

	flows := args[0]

	options := []maestro.Option{
		maestro.WithPath(flows, true),
		maestro.WithCodec(json.NewConstructor()),
		maestro.WithCodec(proto.NewConstructor()),
		maestro.WithCaller(http.NewCaller()),
		maestro.WithSchemaCollection(collection),
	}

	if HTTPAddr != "" {
		listener, err := http.NewListener(HTTPAddr, specs.Options{})
		if err != nil {
			return err
		}

		options = append(options, maestro.WithListener(listener))
	}

	client, err := maestro.New(options...)
	if err != nil {
		return err
	}

	go sigterm(client)

	errs := client.Serve()
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func sigterm(client *maestro.Client) {
	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGINT, syscall.SIGTERM)
	<-term
	client.Close()
}
