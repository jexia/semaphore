package printer

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultHeader creates a default header.
func DefaultHeader(version string) Printer {
	if version == "" {
		version = "not set"
	}

	return Printer{
		"Code generated by Semaphore. DO NOT EDIT.",
		fmt.Sprintf("Semaphore version: %s", version),
		fmt.Sprintf("Timestamp: %s", time.Now().Format(http.TimeFormat)),
	}
}

// Options contains printer settings.
type Options struct {
	StreamStart string
	LineStart   string
	LineEnd     string
	StreamEnd   string
}

// Printer contains lines to be printed.
type Printer []string

// Print lines to the provided writer.
func (printer Printer) Print(dst io.Writer, options Options) error {
	if options.StreamStart != "" {
		if _, err := fmt.Fprint(dst, options.StreamStart); err != nil {
			return err
		}
	}

	for _, line := range printer {
		if _, err := fmt.Fprintf(dst, "%s%s%s", options.LineStart, line, options.LineEnd); err != nil {
			return err
		}
	}

	if options.StreamEnd != "" {
		_, err := fmt.Fprint(dst, options.StreamEnd)

		return err
	}

	return nil
}
