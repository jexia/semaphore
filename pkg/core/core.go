package core

import (
	"github.com/jexia/semaphore/pkg/checks"
	"github.com/jexia/semaphore/pkg/compare"
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/dependencies"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
	"github.com/jexia/semaphore/pkg/specs"
)

// NewSpecs construct a specs manifest from the given options
func NewSpecs(ctx instance.Context, mem functions.Collection, options api.Options) (*specs.Collection, error) {
	collection, err := ResolveProviders(ctx, options)
	if err != nil {
		return nil, err
	}

	ConstructErrorHandle(collection.FlowsManifest)

	err = checks.ManifestDuplicates(ctx, collection.FlowsManifest)
	if err != nil {
		return nil, err
	}

	err = references.DefineManifest(ctx, collection.ServicesManifest, collection.SchemaManifest, collection.FlowsManifest)
	if err != nil {
		return nil, err
	}

	err = functions.PrepareManifestFunctions(ctx, mem, options.Functions, collection.FlowsManifest)
	if err != nil {
		return nil, err
	}

	dependencies.ResolveReferences(ctx, collection.FlowsManifest)

	err = compare.ManifestTypes(ctx, collection.ServicesManifest, collection.SchemaManifest, collection.FlowsManifest)
	if err != nil {
		return nil, err
	}

	err = dependencies.ResolveManifest(ctx, collection.FlowsManifest)
	if err != nil {
		return nil, err
	}

	if options.AfterConstructor != nil {
		err = options.AfterConstructor(ctx, collection)
		if err != nil {
			return nil, err
		}
	}

	return collection, nil
}
