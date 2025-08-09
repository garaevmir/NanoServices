package api

import (
	"fmt"
	"net/http"

	"github.com/nanoservices/tests/models"
	"github.com/nanoservices/tests/utils"
)

func RegisterUser(user models.User) error {
	return utils.SendRequest("POST", "/register", user, http.StatusCreated, nil, "")
}

func LoginUser(user models.User) (string, error) {
	var tokenResp models.TokenResponse
	err := utils.SendRequest("POST", "/login", user, http.StatusOK, &tokenResp, "")
	return tokenResp.Token, err
}

func CreatePost(token string, post models.CreatePostRequest) (string, error) {
	var postResp models.PostResponse
	err := utils.SendRequest("POST", "/posts", post, http.StatusCreated, &postResp, token)
	return postResp.ID, err
}

func ListPosts(token string) ([]models.PostResponse, int, error) {
	var postsResp models.PostsListResponse
	err := utils.SendRequest("GET", "/posts_list", nil, http.StatusOK, &postsResp, token)
	return postsResp.Posts, postsResp.Total, err
}

func GetPostStats(token, postID string) (*models.StatsResponse, error) {
	var stats models.StatsResponse
	err := utils.SendRequest("GET", "/stats/posts/"+postID, nil, http.StatusOK, &stats, token)
	return &stats, err
}

func ViewPost(token, postID string) error {
	return utils.SendRequest("POST", "/posts/view/"+postID, nil, http.StatusOK, nil, token)
}

func GetViewsTrend(token, postID, period string) ([]models.TrendItem, error) {
	url := fmt.Sprintf("/stats/posts/%s/views/trend?period=%s", postID, period)

	var trendItems []models.TrendItem
	err := utils.SendRequest("GET", url, nil, http.StatusOK, &trendItems, token)
	return trendItems, err
}

func UpdateProfile(token string, data models.ProfileUpdateRequest) error {
	return utils.SendRequest("POST", "/profile", data, http.StatusOK, nil, token)
}

func GetProfile(token string) (*models.ProfileResponse, error) {
	var profile models.ProfileResponse
	err := utils.SendRequest("GET", "/profile", nil, http.StatusOK, &profile, token)
	return &profile, err
}

func AddComment(token, postID, content string) error {
	url := fmt.Sprintf("/posts/comment/%s", postID)
	body := models.CommentRequest{Content: content}
	return utils.SendRequest("POST", url, body, http.StatusCreated, nil, token)
}

func GetTopPosts(token, metric, period string) ([]models.TopPostItem, error) {
	url := fmt.Sprintf("/stats/top/posts?metric=%s&period=%s", metric, period)

	var topPosts []models.TopPostItem
	err := utils.SendRequest("GET", url, nil, http.StatusOK, &topPosts, token)
	return topPosts, err
}
