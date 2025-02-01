package main

import (
	"bank-service-app/internal/delivery/http"
	"bank-service-app/internal/repository"
	"bank-service-app/internal/usecase"
	"bank-service-app/pkg/config"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Database connection
	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Repository
	nasabahRepo := repository.NewPostgresNasabahRepository(db)
	transaksiRepo := repository.NewPostgresTransaksiRepository(db)

	// Usecase
	nasabahUsecase := usecase.NewNasabahUsecase(nasabahRepo, transaksiRepo)
	transaksiUsecase := usecase.NewTransaksiUsecase(nasabahRepo, transaksiRepo)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Register handlers
	http.NewNasabahHandler(e, nasabahUsecase, transaksiUsecase)

	// Server configuration
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	// Start server with graceful shutdown
	go func() {
		if err := e.Start(serverAddr); err != nil {
			log.Printf("Shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Failed to gracefully shutdown the server:", err)
	}
}
