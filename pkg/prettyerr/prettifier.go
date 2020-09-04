package prettyerr

// Prettifier builds Error
type Prettifier interface {
	// Prettify the error and return Error
	Prettify() Error
}

// PrettyfierFunc is a functional way to write Pretifier.
type PrettyfierFunc func() Error

func (f PrettyfierFunc) Prettify() Error {
	return f()
}

// Strategy selects and returns a prettifier for the given error
type Strategy interface {
	Match(error) Prettifier
}

// PrettifierStrategy is a strategy based on calling Prettify method on Prettifier.
// If the error does not implement Prettifier, it fallbacks to the generic one.
type PrettifierStrategy struct{}

func (PrettifierStrategy) Match(err error) Prettifier {
	prettifier, ok := err.(Prettifier)
	if ok {
		return prettifier
	}

	return PrettyfierFunc(func() Error {
		return Error{
			Original: err,
			Message:  err.Error(),
			Details:  nil,
			Code:     GenericErrorCode,
		}
	})
}
