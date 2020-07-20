package core

import (
	"github.com/jexia/semaphore/pkg/checks"
	"github.com/jexia/semaphore/pkg/compare"
	"github.com/jexia/semaphore/pkg/core/api"
	"github.com/jexia/semaphore/pkg/core/instance"
	"github.com/jexia/semaphore/pkg/dependencies"
	"github.com/jexia/semaphore/pkg/functions"
	"github.com/jexia/semaphore/pkg/references"
)

// NewSpecs construct a specs manifest from the given options
func NewSpecs(ctx instance.Context, mem functions.Collection, options api.Options) (*api.Collection, error) {
	collection, err := ResolveProviders(ctx, options)
	if err != nil {
		return nil, err
	}

	ConstructErrorHandle(collection.Flows)

	err = checks.ManifestDuplicates(ctx, collection.Flows)
	if err != nil {
		return nil, err
	}

	err = references.DefineManifest(ctx, collection.Services, collection.Schema, collection.Flows)
	if err != nil {
		return nil, err
	}

	err = functions.PrepareManifestFunctions(ctx, mem, options.Functions, collection.Flows)
	if err != nil {
		return nil, err
	}

	dependencies.ResolveReferences(ctx, collection.Flows)

	err = compare.ManifestTypes(ctx, collection.Services, collection.Schema, collection.Flows)
	if err != nil {
		return nil, err
	}

	err = dependencies.ResolveManifest(ctx, collection.Flows)
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
