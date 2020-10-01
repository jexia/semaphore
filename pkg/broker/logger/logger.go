package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jexia/semaphore/pkg/broker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// WithFields adds structured context to the defined logger. Fields added
// to the child don't affect the parent, and vice versa.
func WithFields(parent *broker.Context, fields ...zap.Field) *broker.Context {
	ctx := WithLogger(parent)
	ctx.Zap = ctx.Zap.With(fields...)

	return ctx
}

// WithLogger creates a child context
func WithLogger(parent *broker.Context) *broker.Context {
	ctx := broker.Child(parent)
	atom := zap.NewAtomicLevel()

	if parent.Atom != nil {
		atom.SetLevel(parent.Atom.Level())
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.EpochNanosTimeEncoder

	var encoder zapcore.Encoder

	switch os.Getenv("LOG_ENCODER") {
	case "json":
		encoder = zapcore.NewJSONEncoder(config)
	default:
		encoder = zapcore.NewConsoleEncoder(config)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.Lock(os.Stdout),
		atom,
	)

	ctx.Zap = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel)).Named(ctx.Module)
	ctx.Atom = &atom
	return ctx
}

// SetLevel sets all modules matching the given pattern with the given log level
func SetLevel(ctx *broker.Context, pattern string, level zapcore.Level) error {
	matched, err := filepath.Match(pattern, ctx.Module)
	if err != nil {
		return fmt.Errorf("failed to match pattern: %w", err)
	}

	for _, child := range ctx.Children {
		// errors could only occure inside the pattern which are validate above
		// this error could safely be ignored
		_ = SetLevel(child, pattern, level)
	}

	if (matched || pattern == ctx.Module) && ctx.Atom != nil {
		ctx.Atom.SetLevel(level)
	}

	return nil
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(ctx *broker.Context, msg string, fields ...zap.Field) {
	if ctx == nil {
		return
	}

	if ctx.Zap == nil {
		panic("context logger not set")
	}

	ctx.Zap.Error(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(ctx *broker.Context, msg string, fields ...zap.Field) {
	if ctx == nil {
		return
	}

	if ctx.Zap == nil {
		panic("context logger not set")
	}

	ctx.Zap.Warn(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(ctx *broker.Context, msg string, fields ...zap.Field) {
	if ctx == nil {
		return
	}

	if ctx.Zap == nil {
		panic("context logger not set")
	}

	ctx.Zap.Info(msg, fields...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(ctx *broker.Context, msg string, fields ...zap.Field) {
	if ctx == nil {
		return
	}

	if ctx.Zap == nil {
		panic("context logger not set")
	}

	ctx.Zap.Debug(msg, fields...)
}
