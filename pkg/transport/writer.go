package transport

import (
	"io"

	"github.com/jexia/maestro/pkg/metadata"
)

const (
	// StatusOK represents a 200 header status code
	StatusOK = 200
)

// NewResponseWriter constructs a new response writer for the given io writer
func NewResponseWriter(writer io.WriteCloser) *Writer {
	return &Writer{
		header: metadata.MD{},
		writer: writer,
	}
}

// Writer represents a response writer
type Writer struct {
	writer io.WriteCloser
	header metadata.MD
	status int
}

// Header returns the response header
func (rw *Writer) Header() metadata.MD {
	return rw.header
}

// HeaderStatus sets the header status
func (rw *Writer) HeaderStatus(status int) {
	rw.status = status
}

// Status returns the header status
func (rw *Writer) Status() int {
	return rw.status
}

// Write writes the given byte buffer to the underlaying io Writer
func (rw *Writer) Write(bb []byte) (int, error) {
	return rw.writer.Write(bb)
}

// Close closes the underlaying byte writer
func (rw *Writer) Close() error {
	return rw.writer.Close()
}
