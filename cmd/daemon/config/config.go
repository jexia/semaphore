package config

import (
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// New constructs a new global config
func New() *Maestro {
	return &Maestro{
		HTTP:         HTTP{},
		GraphQL:      GraphQL{},
		GRPC:         GRPC{},
		Protobuffers: []string{},
		Flows:        []string{},
	}
}

// Read attempts to read the given configuration on the given path and decode it into the Maestro configuration
func Read(cmd *cobra.Command, target *Maestro) error {
	flag := cmd.Flag("config")
	if flag == nil {
		return nil
	}

	path := flag.Value.String()
	if path == "" {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	err = yaml.NewDecoder(file).Decode(&target)
	if err != nil {
		return err
	}

	return nil
}

// Maestro configurations
type Maestro struct {
	LogLevel     string   `yaml:"level"`
	HTTP         HTTP     `yaml:"http"`
	GraphQL      GraphQL  `yaml:"graphql"`
	GRPC         GRPC     `yaml:"grpc"`
	Protobuffers []string `yaml:"protobuffers"`
	Flows        []string `yaml:"flows"`
}

// HTTP configurations
type HTTP struct {
	Address string `yaml:"address"`
}

// GRPC configurations
type GRPC struct {
	Address string `yaml:"address"`
}

// GraphQL configurations
type GraphQL struct {
	Address string `yaml:"address"`
}
