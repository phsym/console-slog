# console-slog

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/phsym/console-slog) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/phsym/console-slog/master/LICENSE) [![Build](https://github.com/phsym/console-slog/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/phsym/zeroslog/actions/workflows/go.yml)

A handler for slog that prints colorized logs, similar to zerolog's console writer output

## Example
```go
package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/phsym/console-slog"
)

func main() {
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelDebug}),
	)
	slog.SetDefault(logger)
	slog.Info("Hello world!", "foo", "bar")
	slog.Debug("Debug message")
	slog.Warn("Warning message")
	slog.Error("Error message", "err", errors.New("the error"))

	logger = logger.With("foo", "bar").
		WithGroup("the-group").
		With("bar", "baz")

	logger.Info("group info", "attr", "value")
}
```

![output](./doc/img/output.png)

When setting `console.HandlerOptions.AddSource` to `true`:
```go
console.NewHandler(os.Stderr, &console.HandlerOptions{Level: slog.LevelDebug, AddSource: true})
```
![output-with-source](./doc/img/output-with-source.png)