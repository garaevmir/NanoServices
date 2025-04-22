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

	authMiddleware "github.com/nanoservices/gateway/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userServiceURL := os.Getenv("USER_SERVICE_URL")

	e.POST("/api/register", func(c echo.Context) error {
		return proxyRequest(c, userServiceURL+"/api/register")
	})

	e.POST("/api/login", func(c echo.Context) error {
		return proxyRequest(c, userServiceURL+"/api/login")
	})

	e.GET("/api/profile", func(c echo.Context) error {
		return proxyRequest(c, userServiceURL+"/api/profile")
	})

	e.POST("/api/profile", func(c echo.Context) error {
		return proxyRequest(c, userServiceURL+"/api/profile")
	})
	initGRPC()

	apiGroup := e.Group("")
	apiGroup.Use(authMiddleware.JWTAuth(os.Getenv("JWT_SECRET")))

	apiGroup.POST("/api/posts", CreatePost)

	apiGroup.GET("/api/posts/:id", GetPost)

	apiGroup.PUT("/api/posts/:id", UpdatePost)

	apiGroup.DELETE("/api/posts/:id", DeletePost)

	apiGroup.GET("/api/posts_list", ListPosts)

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
