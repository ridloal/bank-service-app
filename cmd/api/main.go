package main

import (
	"bank-service-app/internal/delivery/http"
	"bank-service-app/internal/delivery/http/middleware"
	"bank-service-app/internal/repository"
	"bank-service-app/internal/usecase"
	"bank-service-app/pkg/config"
	"bank-service-app/pkg/logger"
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()
	defer log.Sync()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.ErrorWithContext("Failed to load config", err)
		os.Exit(1)
	}

	// Database connection
	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		log.ErrorWithContext("Failed to connect to database", err)
		os.Exit(1)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.ErrorWithContext("Failed to ping database", err)
		os.Exit(1)
	}

	log.InfoWithContext("Successfully connected to database")

	// Repository
	nasabahRepo := repository.NewPostgresNasabahRepository(db, log)
	transaksiRepo := repository.NewPostgresTransaksiRepository(db, log)

	// Usecase
	nasabahUsecase := usecase.NewNasabahUsecase(nasabahRepo, transaksiRepo, log)
	transaksiUsecase := usecase.NewTransaksiUsecase(nasabahRepo, transaksiRepo, log)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(echomiddleware.RequestID())
	e.Use(middleware.RequestLogger(log))

	// Register handlers
	http.NewNasabahHandler(e, nasabahUsecase, transaksiUsecase)

	// Server configuration
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	// Start server with graceful shutdown
	go func() {
		if err := e.Start(serverAddr); err != nil {
			log.InfoWithContext("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.ErrorWithContext("Failed to gracefully shutdown the server", err)
	}
}
