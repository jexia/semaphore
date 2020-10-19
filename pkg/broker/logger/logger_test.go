package logger

import (
	"os"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestPrintInfo(t *testing.T) {
	ctx := WithLogger(broker.NewBackground())
	Info(ctx, "mock message")
}

func TestPrintInfoNil(t *testing.T) {
	Info(nil, "mock message")
}

func TestWithFields(t *testing.T) {
	WithFields(broker.NewBackground())
}

func TestPrintInfoNilLogger(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("unexpected pass")
		}
	}()

	ctx := broker.NewBackground()
	Info(ctx, "mock message")
}

func TestPrintWarn(t *testing.T) {
	ctx := WithLogger(broker.NewBackground())
	Warn(ctx, "mock message")
}

func TestPrintWarnNil(t *testing.T) {
	Warn(nil, "mock message")
}

func TestPrintWarnNilLogger(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("unexpected pass")
		}
	}()

	ctx := broker.NewBackground()
	Warn(ctx, "mock message")
}

func TestPrintError(t *testing.T) {
	ctx := WithLogger(broker.NewBackground())
	Error(ctx, "mock message")
}

func TestPrintErrorNil(t *testing.T) {
	Error(nil, "mock message")
}

func TestPrintErrorNilLogger(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("unexpected pass")
		}
	}()

	ctx := broker.NewBackground()
	Error(ctx, "mock message")
}

func TestPrintDebug(t *testing.T) {
	ctx := WithLogger(broker.NewBackground())
	Debug(ctx, "mock message")
}

func TestPrintDebugNil(t *testing.T) {
	Debug(nil, "mock message")
}

func TestPrintDebugNilLogger(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("unexpected pass")
		}
	}()

	ctx := broker.NewBackground()
	Debug(ctx, "mock message")
}

func TestCopyParentAtomLevel(t *testing.T) {
	expected := zapcore.ErrorLevel
	parent := WithLogger(broker.WithModule(broker.NewBackground(), "main"))
	err := SetLevel(parent, "main", expected)
	if err != nil {
		t.Fatal(err)
	}

	child := WithLogger(parent)
	if child.Atom.Level() != expected {
		t.Fatalf("unexpected level %+v, expected %+v", child.Atom.Level(), expected)
	}
}

func TestSetLevelChildren(t *testing.T) {
	expected := zapcore.ErrorLevel
	parent := WithLogger(broker.WithModule(broker.NewBackground(), "main"))
	child := WithLogger(parent)

	err := SetLevel(parent, "*", expected)
	if err != nil {
		t.Fatal(err)
	}

	if child.Atom.Level() != expected {
		t.Fatalf("unexpected level %+v, expected %+v", child.Atom.Level(), expected)
	}
}

func TestSetLevelChildrenErr(t *testing.T) {
	parent := WithLogger(broker.WithModule(broker.NewBackground(), "x"))

	err := SetLevel(parent, "[x-]", zap.PanicLevel)
	if err == nil {
		t.Fatal("unexpected pass")
	}

	if parent.Atom.Level() != zap.ErrorLevel {
		t.Fatalf("unexpected level %+v, expected %+v", parent.Atom.Level(), zap.ErrorLevel)
	}
}

func TestInvalidPattern(t *testing.T) {
	expected := zapcore.ErrorLevel
	parent := WithLogger(broker.WithModule(broker.NewBackground(), "main"))
	_ = WithLogger(parent)

	err := SetLevel(parent, "[]", expected)
	if err == nil {
		t.Fatal("unexpected pass")
	}
}

func TestJSONEncoder(t *testing.T) {
	err := os.Setenv("LOG_ENCODER", "json")
	if err != nil {
		t.Fatal(err)
	}

	WithLogger(broker.NewBackground())
}
