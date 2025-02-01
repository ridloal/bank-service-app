package main

import (
	"bank-service-app/internal/delivery/http"
	"bank-service-app/internal/repository"
	"bank-service-app/internal/usecase"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

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

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
