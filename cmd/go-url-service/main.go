package main

import (
	"go-url-service/internal/config"
	"go-url-service/internal/http-server/handlers/redirect"
	"go-url-service/internal/http-server/handlers/url/save"
	"go-url-service/internal/http-server/middleware/logger"
	"go-url-service/internal/lib/logger/sl"
	"go-url-service/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.NewConfig()

	log := setupLogger(cfg.Env)
	log.Info("Starting service", "env", cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to create storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(logger.New(log))
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("basic-auth", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Get("/{alias}", redirect.New(log, storage))
	})

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Service started", slog.String("address", cfg.Address))

	if err := server.ListenAndServe(); err != nil {
		log.Error("Failed to start server", sl.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
