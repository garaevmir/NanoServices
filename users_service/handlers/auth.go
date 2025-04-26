package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/nanoservices/users_service/models"
	"github.com/nanoservices/users_service/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo   repository.RepositoryInt
	secret string
}

func NewHandlers(repo repository.RepositoryInt, secret string) *UserHandler {
	return &UserHandler{repo: repo, secret: secret}
}

func (h *UserHandler) Register(c echo.Context) error {
	var input models.Register
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	roleID, err := h.repo.GetRoleByName(c.Request().Context(), "user")
	if err != nil {
		roleID.ID, err = h.repo.CreateRole(c.Request().Context(), "user", "Default user role")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create role"})
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Username+input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	userID, err := h.repo.CreateUser(c.Request().Context(), input.Username, string(hashedPassword), roleID.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	_, err = h.repo.CreateProfile(c.Request().Context(), userID, "", "", input.Email, "", "", "")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create profile"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully", "id": userID})
}

func (h *UserHandler) Login(c echo.Context) error {
	var input models.Login
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	user, err := h.repo.GetUserByUsername(c.Request().Context(), input.Username)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid username"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Username+input.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID})
	tokenString, err := token.SignedString([]byte(h.secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "generating token error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Authentication successful", "Token": tokenString})
}

func (h *UserHandler) Profile(c echo.Context) error {
	userId := c.Get("user_id").(string)

	profile, err := h.repo.GetProfileByUserID(c.Request().Context(), userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch profile"})
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	var input models.UpdateProfile
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	userId := c.Get("user_id").(string)
	currentProfile, err := h.repo.GetProfileByUserID(c.Request().Context(), userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch profile"})
	}

	if input.FirstName != "" {
		currentProfile.FirstName = input.FirstName
	}
	if input.LastName != "" {
		currentProfile.LastName = input.LastName
	}
	if input.Email != "" {
		currentProfile.Email = input.Email
	}
	if input.Birthdate != "" {
		birthdateParsed, err := time.Parse("2006-01-02", input.Birthdate)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid birthdate format"})
		}
		currentProfile.Birthdate = birthdateParsed
	}
	if input.PhoneNumber != "" {
		currentProfile.PhoneNumber = input.PhoneNumber
	}
	if input.Bio != "" {
		currentProfile.Bio = input.Bio
	}

	err = h.repo.UpdateProfile(
		c.Request().Context(),
		userId,
		currentProfile.FirstName,
		currentProfile.LastName,
		currentProfile.Email,
		currentProfile.PhoneNumber,
		currentProfile.Bio,
		currentProfile.Birthdate.Format("2006-01-02"),
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update profile"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}
