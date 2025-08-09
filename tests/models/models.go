package models

import "time"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CreatePostRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	IsPrivate   bool     `json:"is_private,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type PostResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UserID      string    `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsPrivate   bool      `json:"is_private"`
	Tags        []string  `json:"tags"`
}

type PostsListResponse struct {
	Posts []PostResponse `json:"posts"`
	Total int            `json:"total"`
}

type StatsResponse struct {
	PostID      string    `json:"post_id"`
	Views       int       `json:"views"`
	Likes       int       `json:"likes"`
	Comments    int       `json:"comments"`
	LastUpdated time.Time `json:"last_updated"`
}

type TrendItem struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type ProfileResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Birthdate   string    `json:"birthdate"`
	PhoneNumber string    `json:"phone_number"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProfileUpdateRequest struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Birthdate   string `json:"birthdate,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Bio         string `json:"bio,omitempty"`
}

type CommentRequest struct {
	Content string `json:"content"`
}

type CommentResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type TopPostItem struct {
	PostID string `json:"post_id"`
	Title  string `json:"title"`
	Count  int    `json:"count"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
