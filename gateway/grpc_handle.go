package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

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
			map[string]string{"error": err.Error()})
	}
}

func ViewPost(c echo.Context) error {
	userID := c.Get("user_id").(string)
	postID := c.Param("id")
	log.Print(postID)

	res, err := postClient.ViewPost(c.Request().Context(), &pb.ViewPostRequest{
		PostId: postID,
		UserId: userID,
	})
	if err != nil {
		return handleGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

func LikePost(c echo.Context) error {
	userID := c.Get("user_id").(string)
	postID := c.Param("id")

	res, err := postClient.LikePost(c.Request().Context(), &pb.LikePostRequest{
		PostId: postID,
		UserId: userID,
	})
	if err != nil {
		return handleGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

func CommentPost(c echo.Context) error {
	userID := c.Get("user_id").(string)
	postID := c.Param("id")

	var req struct {
		Content string `json:"content"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	res, err := postClient.CommentPost(c.Request().Context(), &pb.CommentPostRequest{
		PostId:  postID,
		UserId:  userID,
		Content: req.Content,
	})
	if err != nil {
		return handleGRPCError(c, err)
	}
	return c.JSON(http.StatusCreated, res)
}

func GetComments(c echo.Context) error {
	userID := c.Get("user_id").(string)
	postID := c.Param("id")

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	res, err := postClient.GetComments(c.Request().Context(), &pb.GetCommentsRequest{
		PostId:   postID,
		Page:     int32(page),
		PageSize: int32(pageSize),
		UserId:   userID,
	})
	if err != nil {
		return handleGRPCError(c, err)
	}
	return c.JSON(http.StatusOK, res)
}

var statsClient pb.StatsServiceClient

func initStatsGRPC() {
	conn, err := grpc.NewClient(
		"stats_service:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create stats client: %v", err)
	}
	statsClient = pb.NewStatsServiceClient(conn)
}

func GetPostStats(c echo.Context) error {
	postID := c.Param("id")

	res, err := statsClient.GetPostStats(c.Request().Context(), &pb.PostStatsRequest{
		PostId: postID,
	})

	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"views":    res.Views,
		"likes":    res.Likes,
		"comments": res.Comments,
	})
}

func GetViewsTrend(c echo.Context) error {
	return handleTrendRequest(c, "view")
}

func GetLikesTrend(c echo.Context) error {
	return handleTrendRequest(c, "like")
}

func GetCommentsTrend(c echo.Context) error {
	return handleTrendRequest(c, "comment")
}

func handleTrendRequest(c echo.Context, metric string) error {
	postID := c.Param("id")
	period := c.QueryParam("period")

	var res *pb.PostTrendResponse
	var err error

	switch metric {
	case "view":
		res, err = statsClient.GetViewsTrend(c.Request().Context(), &pb.PostTrendRequest{
			PostId: postID,
			Period: period,
		})
	case "like":
		res, err = statsClient.GetLikesTrend(c.Request().Context(), &pb.PostTrendRequest{
			PostId: postID,
			Period: period,
		})
	case "comment":
		res, err = statsClient.GetCommentsTrend(c.Request().Context(), &pb.PostTrendRequest{
			PostId: postID,
			Period: period,
		})
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid metric"})
	}

	if err != nil {
		return handleGRPCError(c, err)
	}

	response := make([]map[string]interface{}, len(res.Data))
	for i, item := range res.Data {
		response[i] = map[string]interface{}{
			"date":  item.Date,
			"count": item.Count,
		}
	}
	return c.JSON(http.StatusOK, response)
}

func GetTopPosts(c echo.Context) error {
	metric := c.QueryParam("metric")

	var metricEnum pb.Metric
	switch strings.ToLower(metric) {
	case "views":
		metricEnum = pb.Metric_VIEWS
	case "likes":
		metricEnum = pb.Metric_LIKES
	case "comments":
		metricEnum = pb.Metric_COMMENTS
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid metric. Use: views/likes/comments",
		})
	}

	res, err := statsClient.GetTopPosts(c.Request().Context(), &pb.TopRequest{
		Metric: metricEnum,
	})

	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, res.Posts)
}

func GetTopUsers(c echo.Context) error {
	metric := c.QueryParam("metric")
	if metric == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "metric parameter is required"})
	}

	var metricEnum pb.Metric
	switch strings.ToLower(metric) {
	case "views":
		metricEnum = pb.Metric_VIEWS
	case "likes":
		metricEnum = pb.Metric_LIKES
	case "comments":
		metricEnum = pb.Metric_COMMENTS
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid metric. Use: views/likes/comments",
		})
	}

	res, err := statsClient.GetTopUsers(c.Request().Context(), &pb.TopRequest{
		Metric: metricEnum,
	})

	if err != nil {
		return handleGRPCError(c, err)
	}

	return c.JSON(http.StatusOK, res.Users)
}
