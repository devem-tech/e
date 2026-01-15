package e_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/devem-tech/e"
)

func TestW(t *testing.T) {
	t.Parallel()

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		if e.W(nil) != nil {
			t.Error("W(nil) should return nil")
		}
	})

	t.Run("wrap new error", func(t *testing.T) {
		t.Parallel()

		err := errors.New("base error")
		wErr := e.W(err)

		if wErr.Error() != "base error" {
			t.Errorf("expected 'base error', got %v", wErr.Error())
		}

		stack := e.Stack(wErr)
		if len(stack) == 0 {
			t.Fatal("stack should not be empty")
		}

		// Проверяем, что первый кадр — это текущий тест (или функция, вызвавшая W)
		if !strings.Contains(stack[0].Func, "TestW") {
			t.Errorf("top frame should be TestW, got %s", stack[0].Func)
		}
	})

	t.Run("do not double wrap", func(t *testing.T) {
		t.Parallel()

		err := e.W(errors.New("base"))
		stack1 := e.Stack(err)

		err2 := e.W(err)
		stack2 := e.Stack(err2)

		if len(stack1) != len(stack2) {
			t.Error("should not wrap stack twice")
		}
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()

	base := errors.New("base")
	wrapped := e.Wrap(base, "context")

	if wrapped.Error() != "context: base" {
		t.Errorf("expected 'context: base', got %v", wrapped.Error())
	}

	if !errors.Is(wrapped, base) {
		t.Error("errors.Is should work with wrapped error")
	}

	if e.Stack(wrapped) == nil {
		t.Error("Wrap should attach a stack")
	}
}

func TestCallersFilter(t *testing.T) {
	t.Parallel()

	err := e.W(errors.New("test"))
	stack := e.Stack(err)

	for _, frame := range stack {
		// Проверка фильтрации runtime
		if strings.HasPrefix(frame.Func, "runtime.") {
			t.Errorf("stack contains runtime frame: %s", frame.Func)
		}
		// Проверка фильтрации собственного пакета
		if strings.HasPrefix(frame.Func, "github.com/devem-tech/e.") {
			t.Errorf("stack contains internal package frame: %s", frame.Func)
		}
	}
}

func BenchmarkW(b *testing.B) {
	err := errors.New("error")

	for b.Loop() {
		_ = e.W(err)
	}
}
