package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"server/internal/config"
	"server/internal/http-server/handlers/url/save"
	"server/internal/lib/logger/hanlder/slogpretty"
	"server/internal/lib/logger/sl"
	"server/internal/storage/sqlite"
)

func main() {
	cnf := config.MustLoad()

	log := setupLogger(cnf.Env)

	log.Info("Starting project...", slog.String("env", cnf.Env))

	log.Info("Init storage: " + cnf.Storage.StoragePath + cnf.Storage.StorageName)
	storage, err := sqlite.New(cnf.Storage.StoragePath + cnf.Storage.StorageName)
	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/url", save.New(log, storage))

	log.Info(fmt.Sprintf("Server is starting on http://%s", cnf.HTTPServer.GetAddress()))

	srv := &http.Server{
		Addr:         cnf.HTTPServer.GetAddress(),
		Handler:      router,
		ReadTimeout:  cnf.HTTPServer.Timeout,
		WriteTimeout: cnf.HTTPServer.Timeout,
		IdleTimeout:  cnf.HTTPServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", err)
	}

}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettyLogger()
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
