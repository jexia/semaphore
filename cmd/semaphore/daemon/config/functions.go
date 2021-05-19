package config

import (
	"github.com/jexia/semaphore/v2/pkg/functions"
	"github.com/jexia/semaphore/v2/pkg/functions/lib/sprintf"
	"github.com/jexia/semaphore/v2/pkg/functions/lib/strconcat"
)

// DefaultFunctions represents the default functions collection
var DefaultFunctions = functions.Custom{
	"sprintf":   sprintf.Function,
	"strconcat": strconcat.Function,
}
