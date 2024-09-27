package util

import (
	"context"
	"log/slog"
	"time"
)

func Must[T any](fn func() (T, error)) T {
	v, err := fn()
	if err != nil {
		panic(err)
	}
	return v
}

func CaptureExecutionTime(ctx context.Context, scopeName string) func() {
	start := time.Now()
	return func() {
		slog.LogAttrs(ctx, slog.LevelDebug, "Execution time", slog.String("scope_name", scopeName), slog.Duration("took", time.Since(start)))
	}
}
