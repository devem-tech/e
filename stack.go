package e

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const (
	// depth defines the maximum number of stack frames to capture.
	depth = 32
)

// selfPkg stores the package path to filter out internal library frames from the stack trace.
var selfPkg = reflect.TypeOf(withStack{}).PkgPath() + "."

// Frame represents a single stack frame with function name, file path, and line number.
type Frame struct {
	Func string `json:"func"`
	File string `json:"file"`
}

// withStack is an error wrapper that attaches a stack trace.
type withStack struct {
	err error
	pcs []uintptr
}

// Error implements the error interface.
func (w *withStack) Error() string {
	return w.err.Error()
}

// Unwrap implements the error unwrap interface for compatibility with errors.Is/As.
func (w *withStack) Unwrap() error {
	return w.err
}

// LogValue implements slog.LogValuer for structured logging support.
// It allows loggers like slog.JSONHandler to automatically format the stack trace.
func (w *withStack) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("msg", w.err.Error()),
		slog.Any("stack", w.Stack()),
	)
}

// Stack decodes the captured program counters into human-readable Frames.
// It lazily evaluates the stack trace only when requested and filters out internal frames.
func (w *withStack) Stack() []Frame {
	if len(w.pcs) == 0 {
		return nil
	}

	frames := make([]Frame, 0, len(w.pcs))
	it := runtime.CallersFrames(w.pcs)

	for {
		f, more := it.Next()

		// Filter out runtime and this library's internal frames.
		if strings.HasPrefix(f.Function, "runtime.") || strings.HasPrefix(f.Function, selfPkg) {
			if !more {
				break
			}

			continue
		}

		frames = append(frames, Frame{
			Func: f.Function,
			File: f.File + ":" + strconv.Itoa(f.Line),
		})

		if !more {
			break
		}
	}

	return frames
}

// W (Wrap) attaches a stack trace to the error if it doesn't already have one.
func W(err error) error {
	if err == nil {
		return nil
	}

	// Fast path: check if the error already contains a stack trace without using reflection.
	if _, ok := err.(*withStack); ok { //nolint:errorlint
		return err
	}

	// Slow path: check if the error implements the stack trace interface using errors.As.
	var ws interface{ Stack() []Frame }
	if errors.As(err, &ws) {
		return err
	}

	return &withStack{
		err: err,
		pcs: callers(),
	}
}

// Wrap adds a context message and attaches a stack trace (if not already present).
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return W(fmt.Errorf("%s: %w", msg, err))
}

// Stack returns the stack trace frames from the error if available.
func Stack(err error) []Frame {
	var ws interface{ Stack() []Frame }
	if errors.As(err, &ws) {
		return ws.Stack()
	}

	return nil
}

// callers captures the program counters for the current goroutine.
func callers() []uintptr {
	var pcs [depth]uintptr

	n := runtime.Callers(3, pcs[:])
	res := make([]uintptr, n)
	copy(res, pcs[:n])

	return res
}
