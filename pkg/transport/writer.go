package transport

import (
	"io"

	"github.com/jexia/maestro/pkg/metadata"
)

// NewResponseWriter constructs a new response writer for the given io writer
func NewResponseWriter(writer io.Writer) *Writer {
	return &Writer{
		header: metadata.MD{},
		writer: writer,
	}
}

// Writer represents a response writer
type Writer struct {
	writer io.Writer
	header metadata.MD
}

// Header returns the response header
func (rw *Writer) Header() metadata.MD {
	return rw.header
}

// Write writes the given byte buffer to the underlaying io Writer
func (rw *Writer) Write(bb []byte) (int, error) {
	return rw.writer.Write(bb)
}
