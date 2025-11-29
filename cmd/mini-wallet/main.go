package main

import (
	"errors"
	"fmt"
	"log/slog"
	"mini-wallet/internal/config"
	"mini-wallet/internal/domain/wallet"
	"mini-wallet/internal/infrastructure/http/handlers"
	mw "mini-wallet/internal/infrastructure/http/middleware"
	"mini-wallet/internal/infrastructure/storage/postgres"
	"mini-wallet/pkg/sl_logger/sl"
	"mini-wallet/pkg/sl_logger/slogpretty"
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
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	pgConfig := postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s port=%s password=%s dbname=%s sslmode=%s",
			cfg.Host,
			cfg.User,
			cfg.Port,
			cfg.Password,
			cfg.Dbname,
			cfg.Sslmode,
		),
		Seed:           cfg.Seed,
		MigrationsPath: cfg.MigrationsPath,
	}

	log.Info("CHECKING DB Conn,", slog.String("Trying to connect with DSN", pgConfig.DSN))
	storage, err := postgres.New(pgConfig, log)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
	}

	svc := wallet.NewService(log, storage)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(mw.New(log))

	r.Post("/api/v1/wallet", handlers.NewAddTransaction(log, svc))
	r.Get("/api/v1/wallets/{WALLET_UUID}", handlers.NewGetWallet(log, svc))
	
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	log.Info("starting HTTP server",
		slog.String("address", cfg.Address),
		slog.Duration("read_timeout", cfg.HTTPServer.Timeout),
		slog.Duration("write_timeout", cfg.HTTPServer.Timeout),
		slog.Duration("idle_timeout", cfg.IdleTimeout),
	)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("failed to start server", sl.Err(err))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
