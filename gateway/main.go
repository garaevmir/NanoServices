package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userServiceURL := os.Getenv("USER_SERVICE_URL")

	e.POST("/register", func(c echo.Context) error {
		return proxyRequest(c, userServiceURL+"/register")
	})

	e.POST("/login", func(c echo.Context) error {
		return proxyRequest(c, userServiceURL+"/login")
	})

	e.GET("/profile", func(c echo.Context) error {
		return proxyRequest(c, userServiceURL+"/profile")
	})

	e.POST("/profile", func(c echo.Context) error {
		return proxyRequest(c, userServiceURL+"/profile")
	})

	s := &http.Server{
		Addr: ":8080",
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

func proxyRequest(c echo.Context, targetURL string) error {
	client := &http.Client{}
	reqBody, _ := io.ReadAll(c.Request().Body)
	fmt.Println(string(reqBody))
	req, _ := http.NewRequest(c.Request().Method, targetURL, bytes.NewBuffer(reqBody))
	req.Header = c.Request().Header

	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error proxying request")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.String(resp.StatusCode, string(body))
}
