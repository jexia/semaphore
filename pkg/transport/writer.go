package transport

import (
	"io"

	"github.com/jexia/semaphore/v2/pkg/codec/metadata"
)

const (
	// StatusOK represents a 200 header status code
	StatusOK = 200
	// StatusInternalErr represents a 500 header status code
	StatusInternalErr = 500
)

var codes = map[int]string{
	StatusOK:          "OK",
	StatusInternalErr: "Internal Server Error",
}

// StatusMessage attempts to lookup the message for the given status code.
// If no message has been found fot the given status code is a empty string returned.
func StatusMessage(status int) string {
	return codes[status]
}

// NewResponseWriter constructs a new response writer for the given io writer
func NewResponseWriter(writer io.WriteCloser) *Writer {
	return &Writer{
		header: metadata.MD{},
		writer: writer,
	}
}

// Writer represents a response writer
type Writer struct {
	writer  io.WriteCloser
	header  metadata.MD
	status  int
	message string
}

// Header returns the response header
func (rw *Writer) Header() metadata.MD {
	return rw.header
}

// HeaderStatus sets the header status
func (rw *Writer) HeaderStatus(status int) {
	rw.status = status
}

// HeaderMessage sets the header message
func (rw *Writer) HeaderMessage(message string) {
	rw.message = message
}

// Status returns the header status
func (rw *Writer) Status() int {
	return rw.status
}

// Message returns the header status message
func (rw *Writer) Message() string {
	return rw.message
}

// Write writes the given byte buffer to the underlaying io Writer
func (rw *Writer) Write(bb []byte) (int, error) {
	return rw.writer.Write(bb)
}

// Close closes the underlaying byte writer
func (rw *Writer) Close() error {
	return rw.writer.Close()
}
