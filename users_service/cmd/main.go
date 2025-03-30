package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/nanoservices/users_service/handlers"
	authMiddleware "github.com/nanoservices/users_service/middleware"
	"github.com/nanoservices/users_service/repository"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=10",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	repo := repository.NewRepository(pool)
	handlers := handlers.NewHandlers(repo, os.Getenv("JWT_SECRET"))

	e.POST("/api/register", handlers.Register)
	e.POST("/api/login", handlers.Login)
	api := e.Group("")
	api.Use(authMiddleware.JWTAuth(os.Getenv("JWT_SECRET")))
	api.GET("/api/profile", handlers.Profile)
	api.POST("/api/profile", handlers.UpdateProfile)

	s := &http.Server{
		Addr: ":8081",
	}

	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("HTTP server shutdown error:", err)
	}

	select {
	case <-ctx.Done():
		log.Printf("timeout of 5 seconds.\n")
	}
}
