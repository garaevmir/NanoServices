package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authMiddleware "github.com/nanoservices/gateway/middleware"
	"github.com/segmentio/kafka-go"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var kafkaWriter *kafka.Writer

func initKafka() {
	kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "user_registrations",
	})
}

func main() {
	initKafka()
	defer kafkaWriter.Close()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userServiceURL := os.Getenv("USER_SERVICE_URL")

	e.POST("/api/register", func(c echo.Context) error {
		body, statusCode, err := proxyRequest(c, userServiceURL+"/api/register")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal server error")
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(body, &resp); err != nil {
			log.Printf("Failed to parse response: %v", err)
			return c.String(http.StatusInternalServerError, "Invalid response format")
		}

		userID, ok := resp["id"].(string)
		if !ok {
			log.Printf("Missing 'id' in response: %v", resp)
			return c.String(http.StatusInternalServerError, "User ID not found")
		}

		msg := fmt.Sprintf(`{"user_id": "%s", "timestamp": "%s"}`,
			userID, time.Now().Format(time.RFC3339))
		if err := kafkaWriter.WriteMessages(context.Background(),
			kafka.Message{Value: []byte(msg)}); err != nil {
			log.Printf("Failed to send Kafka message: %v", err)
		}

		return c.String(statusCode, string(body))
	})

	e.POST("/api/login", func(c echo.Context) error {
		body, statusCode, _ := proxyRequest(c, userServiceURL+"/api/login")
		return c.String(statusCode, string(body))
	})

	e.GET("/api/profile", func(c echo.Context) error {
		body, statusCode, _ := proxyRequest(c, userServiceURL+"/api/profile")
		return c.String(statusCode, string(body))
	})

	e.POST("/api/profile", func(c echo.Context) error {
		body, statusCode, _ := proxyRequest(c, userServiceURL+"/api/profile")
		return c.String(statusCode, string(body))
	})
	initGRPC()
	initStatsGRPC()

	apiGroup := e.Group("")
	apiGroup.Use(authMiddleware.JWTAuth(os.Getenv("JWT_SECRET")))

	apiGroup.POST("/api/posts", CreatePost)

	apiGroup.GET("/api/posts/:id", GetPost)

	apiGroup.PUT("/api/posts/:id", UpdatePost)

	apiGroup.DELETE("/api/posts/:id", DeletePost)

	apiGroup.GET("/api/posts_list", ListPosts)
	apiGroup.POST("/api/posts/view/:id", ViewPost)
	apiGroup.POST("/api/posts/like/:id", LikePost)
	apiGroup.POST("/api/posts/comment/:id", CommentPost)
	apiGroup.GET("/api/posts/comments/:id", GetComments)

	apiGroup.GET("/stats/posts/:id", GetPostStats)
	apiGroup.GET("/stats/posts/:id/views/trend", GetViewsTrend)
	apiGroup.GET("/stats/posts/:id/likes/trend", GetLikesTrend)
	apiGroup.GET("/stats/posts/:id/comments/trend", GetCommentsTrend)
	apiGroup.GET("/stats/top/posts", GetTopPosts)
	apiGroup.GET("/stats/top/users", GetTopUsers)

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

func proxyRequest(c echo.Context, targetURL string) ([]byte, int, error) {
	client := &http.Client{}
	reqBody, _ := io.ReadAll(c.Request().Body)
	req, _ := http.NewRequest(
		c.Request().Method,
		targetURL,
		bytes.NewBuffer(reqBody),
	)
	req.Header = c.Request().Header.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("proxy error: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body, resp.StatusCode, nil
}
