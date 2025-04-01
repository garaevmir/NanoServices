package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	pb "github.com/nanoservices/gateway/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var postClient pb.PostServiceClient

func initGRPC() {
	conn, err := grpc.NewClient(
		"events_service:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	postClient = pb.NewPostServiceClient(conn)
}

func CreatePost(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		IsPrivate   bool     `json:"is_private"`
		Tags        []string `json:"tags"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	grpcReq := &pb.CreatePostRequest{
		Title:       req.Title,
		Description: req.Description,
		UserId:      userID,
		IsPrivate:   req.IsPrivate,
		Tags:        req.Tags,
	}

	res, err := postClient.CreatePost(c.Request().Context(), grpcReq)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(http.StatusCreated, res)
}

func GetPost(c echo.Context) error {
	userID := c.Get("user_id").(string)
	postID := c.Param("id")

	res, err := postClient.GetPost(c.Request().Context(), &pb.GetPostRequest{
		PostId: postID,
		UserId: userID,
	})

	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func UpdatePost(c echo.Context) error {
	userID := c.Get("user_id").(string)
	postID := c.Param("id")

	var req struct {
		Title       *string  `json:"title"`
		Description *string  `json:"description"`
		IsPrivate   *bool    `json:"is_private"`
		Tags        []string `json:"tags"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	grpcReq := &pb.UpdatePostRequest{
		PostId: postID,
		UserId: userID,
		Tags:   req.Tags,
	}

	if req.Title != nil {
		titleValue := *req.Title
		grpcReq.Title = &titleValue
	}
	if req.Description != nil {
		descValue := *req.Description
		grpcReq.Description = &descValue
	}
	if req.IsPrivate != nil {
		privateValue := *req.IsPrivate
		grpcReq.IsPrivate = &privateValue
	}

	res, err := postClient.UpdatePost(c.Request().Context(), grpcReq)
	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func DeletePost(c echo.Context) error {
	userID := c.Get("user_id").(string)
	postID := c.Param("id")

	_, err := postClient.DeletePost(c.Request().Context(), &pb.DeletePostRequest{
		PostId: postID,
		UserId: userID,
	})

	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "post deleted"})
}

func ListPosts(c echo.Context) error {
	userID := c.Get("user_id").(string)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	res, err := postClient.ListPosts(c.Request().Context(), &pb.ListPostsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		UserId:   userID,
	})

	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func handleGRPCError(c echo.Context, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": "internal server error"})
	}

	switch st.Code() {
	case codes.NotFound:
		return c.JSON(http.StatusNotFound,
			map[string]string{"error": "post not found"})
	case codes.PermissionDenied:
		return c.JSON(http.StatusForbidden,
			map[string]string{"error": "access denied"})
	case codes.InvalidArgument:
		return c.JSON(http.StatusBadRequest,
			map[string]string{"error": "invalid arguments"})
	default:
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": "internal server error"})
	}
}
