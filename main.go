package maestro

import (
	"os"
	"path/filepath"

	"github.com/jexia/maestro/schema"
	"github.com/jexia/maestro/specs"
	"github.com/jexia/maestro/specs/intermediate"
	"github.com/jexia/maestro/specs/strict"
	"github.com/jexia/maestro/specs/trace"
	"github.com/jexia/maestro/utils"
)

// Option represents a constructor func which sets a given option
type Option func(*Options)

// Options represents all the available options
type Options struct {
	Path      string
	Recursive bool
	Schema    schema.Collection
	Functions specs.CustomDefinedFunctions
}

// NewOptions constructs a options object from the given option constructors
func NewOptions(options ...Option) Options {
	result := Options{}
	for _, option := range options {
		option(&result)
	}
	return result
}

// WithPath defines the definitions path to be used
func WithPath(path string, recursive bool) Option {
	return func(options *Options) {
		options.Path = path
		options.Recursive = recursive
	}
}

// WithSchemaCollection defines the schema collection to be used
func WithSchemaCollection(collection schema.Collection) Option {
	return func(options *Options) {
		options.Schema = collection
	}
}

// WithFunctions defines the custom defined functions to be used
func WithFunctions(functions specs.CustomDefinedFunctions) Option {
	return func(options *Options) {
		options.Functions = functions
	}
}

// New constructs a new Maestro instance
func New(opts ...Option) (*specs.Manifest, error) {
	options := NewOptions(opts...)

	if options.Path == "" {
		return nil, trace.New(trace.WithMessage("undefined path in options"))
	}

	if options.Schema == nil {
		return nil, trace.New(trace.WithMessage("undefined schema in options"))
	}

	files, err := utils.ReadDir(options.Path, options.Recursive, intermediate.Ext)
	if err != nil {
		return nil, err
	}

	manifest := &specs.Manifest{}

	for _, file := range files {
		reader, err := os.Open(filepath.Join(file.Path, file.Name()))
		if err != nil {
			return nil, err
		}

		definition, err := intermediate.UnmarshalHCL(file.Name(), reader)
		if err != nil {
			return nil, err
		}

		result, err := intermediate.ParseManifest(definition, options.Functions)
		if err != nil {
			return nil, err
		}

		manifest.MergeLeft(result)

		err = specs.CheckManifestDuplicates(file.Name(), manifest)
		if err != nil {
			return nil, err
		}
	}

	err = specs.ResolveManifestDependencies(manifest)
	if err != nil {
		panic(err)
	}

	err = strict.Define(options.Schema, manifest)
	if err != nil {
		panic(err)
	}

	return manifest, nil
}
