package main

import (
	"bank-service-app/internal/delivery/http"
	"bank-service-app/internal/delivery/http/middleware"
	"bank-service-app/internal/domain"
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
	log := initLogger()
	defer log.Sync()

	cfg := initConfig(log)
	db := initDatabase(cfg, log)
	defer db.Close()

	repos := initRepositories(db, log)
	usecases := initUsecases(repos, log)

	e := initEcho()
	setupMiddleware(e, log)
	setupHandlers(e, usecases)

	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	startServer(e, serverAddr, cfg, log)
}

func initLogger() *logger.Logger {
	return logger.NewLogger()
}

func initConfig(log *logger.Logger) *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.ErrorWithContext("Failed to load config", err)
		os.Exit(1)
	}
	return cfg
}

func initDatabase(cfg *config.Config, log *logger.Logger) *sql.DB {
	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		log.ErrorWithContext("Failed to connect to database", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		log.ErrorWithContext("Failed to ping database", err)
		os.Exit(1)
	}

	log.InfoWithContext("Successfully connected to database")
	return db
}

type repositories struct {
	nasabah   domain.NasabahRepository
	transaksi domain.TransaksiRepository
}

func initRepositories(db *sql.DB, log *logger.Logger) *repositories {
	return &repositories{
		nasabah:   repository.NewPostgresNasabahRepository(db, log),
		transaksi: repository.NewPostgresTransaksiRepository(db, log),
	}
}

type usecases struct {
	nasabah   domain.NasabahUsecase
	transaksi domain.TransaksiUsecase
}

func initUsecases(repos *repositories, log *logger.Logger) *usecases {
	return &usecases{
		nasabah:   usecase.NewNasabahUsecase(repos.nasabah, repos.transaksi, log),
		transaksi: usecase.NewTransaksiUsecase(repos.nasabah, repos.transaksi, log),
	}
}

func initEcho() *echo.Echo {
	return echo.New()
}

func setupMiddleware(e *echo.Echo, log *logger.Logger) {
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(echomiddleware.RequestID())
	e.Use(middleware.RequestLogger(log))
}

func setupHandlers(e *echo.Echo, usecases *usecases) {
	http.NewNasabahHandler(e, usecases.nasabah, usecases.transaksi)
}

func startServer(e *echo.Echo, serverAddr string, cfg *config.Config, log *logger.Logger) {
	// Start server
	go func() {
		if err := e.Start(serverAddr); err != nil {
			log.InfoWithContext("Shutting down the server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.ErrorWithContext("Failed to gracefully shutdown the server", err)
	}
}
