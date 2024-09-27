package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"project-template/config"
	"project-template/database"
	"project-template/env"
	"project-template/http"
	"project-template/logger"
	"project-template/services/user"

	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "Fatal server error", slog.Any("error", err))
	}
}

func run(ctx context.Context) error {
	log := logger.New()
	slog.SetDefault(log)

	conf, err := config.New()
	if err != nil {
		return err
	}

	slog.LogAttrs(ctx, slog.LevelDebug, "Debug messages are enabled", slog.Any("environment", env.Get()))
	slog.LogAttrs(ctx, slog.LevelInfo, "Loaded configuration", slog.Any("configuration", conf))

	db, err := database.New(&conf.Database)
	if err != nil {
		return err
	}

	server, err := http.NewServer(&conf.Server)
	if err != nil {
		return err
	}

	server.Use(
		middleware.Heartbeat("/health"),
		slogchi.New(log),
		middleware.Recoverer,
		middleware.CleanPath,
		middleware.RequestID,
		middleware.RealIP,
		middleware.AllowContentType("application/json", "application/json; charset=utf-8"),
	)

	userService := user.New(db)
	server.MountSubrouter("/api/users", userService.Bind())

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go server.Run(ctx)
	defer server.Close()

	select {
	case sig := <-exit:
		slog.LogAttrs(ctx, slog.LevelInfo, "Graceful shutdown after receiving signal", slog.Any("signal", sig))
		return nil
	case err := <-server.Notify():
		slog.LogAttrs(ctx, slog.LevelError, "Graceful shutdown after server error", slog.Any("error", err))
		return err
	}
}
