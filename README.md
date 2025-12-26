# E

`e` is a minimalist, high-performance Go library for attaching stack traces to errors.

## 🚀 Key Features

- 🧬 **Zero-allocation capture**: stores only raw `uintptr` (program counters).
- 🦉 **Lazy formatting**: decodes the stack trace only when it's actually accessed.
- 🧹 **Built-in filtering**: automatically skips `runtime` and internal library frames.

## 📦 Installation

```bash
go get -u github.com/devem-tech/e
```

## 🔧 Usage

### Basic Wrapping

Attach a stack trace to an error:

```go
return e.W(err)
```

Add contextual message and a stack trace:

```go
return e.Wrap(err, "failed to get user")
```

The stack trace is added only **if the error does not already contain one**.

#### Example

```go
package main

import (
    "fmt"

    "github.com/devem-tech/e"
)

func findUser(id int) error {
    return fmt.Errorf("database connection refused")
}

func logic() error {
    err := findUser(1)
    if err != nil {
        return e.Wrap(err, "failed to get user")
    }

    return nil
}

func main() {
    err := logic()
    if err != nil {
        fmt.Println("Error:", err)
        
        // Retrieve frames manually if needed
        for _, f := range e.Stack(err) {
            fmt.Printf("%s\n\t%s\n", f.Func, f.File)
        }
    }
}
```

### Logging with `slog` (recommended)

The error type implements `slog.LogValuer`, so stack traces are automatically embedded into structured logs:

```go
slog.Error("operation failed", "error", err)
```

#### Example JSON output

```json
{
  "level": "ERROR",
  "msg": "operation failed",
  "error": {
    "msg": "failed to get user: database connection refused",
    "stack": [
      {
        "func": "main.logic",
        "file": "/app/main.go:15"
      },
      {
        "func": "main.main",
        "file": "/app/main.go:22"
      }
    ]
  }
}
```

### Manual access to stack frames

If you need to process stack frames manually:

```go
for _, f := range e.Stack(err) {
    fmt.Println(f.Func, f.File)
}
```