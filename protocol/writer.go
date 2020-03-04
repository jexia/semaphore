package protocol

import (
	"io"
)

// NewResponseWriter constructs a new response writer for the given io writer
func NewResponseWriter(writer io.Writer) *Writer {
	return &Writer{
		writer: writer,
	}
}

// Writer represents a response writer
type Writer struct {
	writer io.Writer
	header Header
	status int
}

// Header returns the response header
func (rw *Writer) Header() Header {
	return rw.header
}

// Write writes the given byte buffer to the underlaying io Writer
func (rw *Writer) Write(bb []byte) (int, error) {
	return rw.writer.Write(bb)
}

// WriteHeader writes the status code header
func (rw *Writer) WriteHeader(status int) {
	rw.status = status
}
