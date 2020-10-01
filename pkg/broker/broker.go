package broker

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewContext constructs a new empty context
func NewContext() *Context {
	return &Context{}
}

// NewBackground constructs a new context with a Atom log level set
// to fatal level avoiding to write logs.
func NewBackground() *Context {
	atom := zap.NewAtomicLevelAt(zapcore.ErrorLevel)

	return &Context{
		Atom: &atom,
	}
}

// A Context carries the module name and logger across module boundaries
type Context struct {
	Parent   *Context
	Name     string
	Module   string
	Zap      *zap.Logger
	Atom     *zap.AtomicLevel
	Children []*Context
}

// Child creates a new context sets the given context as
// it's parent and appends the newly created context as a child
func Child(parent *Context) *Context {
	if parent == nil {
		return &Context{}
	}

	child := &Context{
		Parent: parent,
		Name:   parent.Name,
		Module: parent.Module,
		Atom:   parent.Atom,
		Zap:    parent.Zap,
	}

	parent.Children = append(parent.Children, child)
	return child
}

// WithModule creates a new child context with the given module name
func WithModule(parent *Context, parts ...string) *Context {
	child := Child(parent)
	child.Name = strings.Join(parts, ".")

	if parent.Module != "" {
		child.Module = strings.Join(append([]string{parent.Module}, parts...), ".")
		return child
	}

	child.Module = child.Name
	return child
}
