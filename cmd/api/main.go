package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abass-codes/peakcloud/internal/cache"
	"github.com/abass-codes/peakcloud/internal/config"
	"github.com/abass-codes/peakcloud/internal/database"
	"github.com/abass-codes/peakcloud/internal/health"
	httpserver "github.com/abass-codes/peakcloud/internal/http"
	"github.com/abass-codes/peakcloud/internal/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = log.Sync()
	}()

	ctx := context.Background()

	db, err := database.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("connect to postgres", zap.Error(err))
	}
	defer db.Close()

	redisClient, err := cache.NewRedis(
		ctx,
		cfg.RedisAddr,
		cfg.RedisPassword,
		cfg.RedisDB,
	)
	if err != nil {
		log.Fatal("connect to redis", zap.Error(err))
	}
	defer func() {
		_ = redisClient.Close()
	}()

	healthService := health.NewService(db, redisClient)
	healthHandler := health.NewHandler(healthService)

	router := httpserver.NewRouter(healthHandler)

	server := &http.Server{
		Addr:              cfg.APIHost + ":" + cfg.APIPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(
			"PeakCloud API started",
			zap.String("address", server.Addr),
			zap.String("environment", cfg.AppEnv),
		)

		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)

	signal.Notify(
		shutdown,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", zap.Error(err))
		}

	case sig := <-shutdown:
		log.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", zap.Error(err))
	}

	log.Info("PeakCloud API stopped")
}
