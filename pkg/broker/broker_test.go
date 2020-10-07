package broker

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext()
	if ctx == nil {
		t.Fatal("unexpected empty context")
	}
}

func TestNewBackground(t *testing.T) {
	ctx := NewBackground()
	if ctx == nil {
		t.Fatal("unexpected empty context")
	}

	if ctx.Atom == nil {
		t.Fatal("context log level Atom not set")
	}

	if ctx.Atom.Level() != zapcore.ErrorLevel {
		t.Fatalf("unexpected atom log level %s, expected %s", ctx.Atom.Level(), zapcore.ErrorLevel)
	}
}

func TestChild(t *testing.T) {
	parent := NewContext()
	if parent == nil {
		t.Fatal("unexpected empty parent context")
	}

	child := Child(parent)
	if child == nil {
		t.Fatal("unexpected empty child context")
	}

	if len(parent.Children) != 1 {
		t.Fatal("child not set")
	}

	if child.Parent == nil {
		t.Fatal("parent not set")
	}
}

func TestChildNilParent(t *testing.T) {
	child := Child(nil)
	if child == nil {
		t.Fatal("unexpected empty child context")
	}

	if child.Parent != nil {
		t.Fatalf("unexpected parent %+v", child.Parent)
	}
}

func TestWithModule(t *testing.T) {
	expected := "mock.module"

	parent := NewContext()
	if parent == nil {
		t.Fatal("unexpected empty parent")
	}

	module := WithModule(parent, "mock", "module")
	if module == nil {
		t.Fatal("unexpected empty module")
	}

	if module.Name != expected {
		t.Fatal("unexpected module name")
	}

	if module.Module != expected {
		t.Fatal("unexpected module path")
	}
}

func TestWithParentModulePath(t *testing.T) {
	expectedN := "mock.module"
	expectedP := "top.mock.module"

	parent := WithModule(NewContext(), "top")
	if parent == nil {
		t.Fatal("unexpected empty parent")
	}

	module := WithModule(parent, "mock", "module")
	if module == nil {
		t.Fatal("unexpected empty module")
	}

	if module.Name != expectedN {
		t.Fatal("unexpected module name")
	}

	if module.Module != expectedP {
		t.Fatal("unexpected module path")
	}
}
