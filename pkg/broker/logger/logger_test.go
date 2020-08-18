package logger

import (
	"os"
	"testing"

	"github.com/jexia/semaphore/pkg/broker"
	"go.uber.org/zap/zapcore"
)

func TestPrintInfo(t *testing.T) {
	ctx := WithLogger(broker.NewBackground())
	Info(ctx, "mock message")
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
