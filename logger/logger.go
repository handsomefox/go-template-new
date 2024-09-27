package logger

import (
	"log/slog"
	"os"

	"project-template/env"

	"github.com/phsym/console-slog"
)

func New() *slog.Logger {
	switch env.Get() {
	case env.Local, env.Testing, env.Staging, env.Development:
		return slog.New(
			console.NewHandler(os.Stderr, &console.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelDebug,
			}),
		)
	case env.CI, env.Production:
		return slog.New(
			slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: true,
				Level:     slog.LevelInfo,
			}),
		)
	default:
		return slog.Default()
	}
}
