package functions

import (
	"github.com/jexia/semaphore/cmd/semaphore/functions/strings/sprintf"
	"github.com/jexia/semaphore/cmd/semaphore/functions/strings/strconcat"
	"github.com/jexia/semaphore/pkg/functions"
)

// Default represents the default functions collection
var Default = functions.Custom{
	"sprintf":   sprintf.Function,
	"strconcat": strconcat.Function,
}
