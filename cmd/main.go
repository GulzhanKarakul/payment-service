package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GulzhanKarakul/payment-service/pkg/config"
	"github.com/GulzhanKarakul/payment-service/pkg/database"
	"github.com/GulzhanKarakul/payment-service/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)
	log.Info("starting payment-service...")

	// connect to database
	db, err := database.NewPostgres(database.DefaultConfig(cfg.DSN()))
	if err != nil{
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	log.Info("connected to database")

	// Http server with timeout
	srv := &http.Server{
		Addr: ":" + cfg.Server.Port,
		Handler: nil,
		ReadTimeout: 15 * time.Second, // request timeout
		WriteTimeout: 15 * time.Second, // response timeout
		IdleTimeout: 60 * time.Second, // keep-alive timeout
	}
	log.Info("server started", "port", cfg.Server.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", "error", err)
	}

	log.Info("server stopped")
}