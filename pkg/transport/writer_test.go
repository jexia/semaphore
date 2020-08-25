package transport

import (
	"testing"

	"github.com/jexia/semaphore/pkg/codec/metadata"
)

type MockWriterCloser struct{}

func (mock *MockWriterCloser) Write(p []byte) (n int, err error) { return len(p), nil }
func (mock *MockWriterCloser) Close() error                      { return nil }

func TestNewResponseWriter(t *testing.T) {
	writer := NewResponseWriter(&MockWriterCloser{})
	if writer == nil {
		t.Fatal("unexpected empty writer")
	}
}

func TestNewResponseWriterNil(t *testing.T) {
	writer := NewResponseWriter(nil)
	if writer == nil {
		t.Fatal("unexpected empty writer")
	}
}

func TestWriterMethods(t *testing.T) {
	wc := &MockWriterCloser{}
	header := metadata.MD{"key": "value"}
	status := 500
	message := "unexpected mock"

	writer := &Writer{
		writer:  wc,
		header:  header,
		status:  status,
		message: message,
	}

	if len(writer.Header()) != len(header) {
		t.Errorf("unexpected header length %d, expected %d", len(writer.Header()), len(header))
	}

	if writer.Status() != status {
		t.Errorf("unexpected status %d, expected %d", writer.Status(), status)
	}

	if writer.Message() != message {
		t.Errorf("unexpected message %s, expected %s", writer.Message(), message)
	}

	writer.HeaderStatus(404)
	if writer.Status() != 404 {
		t.Errorf("unexpected header status %d, expected %d", writer.Status(), 404)
	}

	writer.HeaderMessage("mock")
	if writer.Message() != "mock" {
		t.Errorf("unexpected header message %s, expected %s", writer.Message(), "mock")
	}

	bb := []byte("message")
	l, err := writer.Write(bb)
	if l != len(bb) {
		t.Errorf("unexpected length %d, expected %d", l, len(bb))
	}

	if err != nil {
		t.Error(err)
	}

	err = writer.Close()
	if err != nil {
		t.Error(err)
	}
}
