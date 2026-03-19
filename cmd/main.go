package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GulzhanKarakul/payment-service/pkg/config"
	"github.com/GulzhanKarakul/payment-service/pkg/database"
)

func main() {
	// Srtuctured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("starting payment-service...")
	cfg := config.Load()

	// connect to database
	db, err := database.NewPostgres(database.DefaultConfig(cfg.DSN()))
	if err != nil{
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("connected to database")

	// Http server with timeout
	srv := &http.Server{
		Addr: ":" + cfg.Server.Port,
		Handler: nil,
		ReadTimeout: 15 * time.Second, // request timeout
		WriteTimeout: 15 * time.Second, // response timeout
		IdleTimeout: 60 * time.Second, // keep-alive timeout
	}
	logger.Info("server started", "port", cfg.Server.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("forced shutdown", "error", err)
	}

	logger.Info("server stopped")
}