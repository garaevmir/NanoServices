package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"github.com/nanoservices/users_service/mocks"
	"github.com/nanoservices/users_service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	repoMock := new(mocks.MockRepository)
	handler := NewHandlers(repoMock, "secret")

	t.Run("Successful registration", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"john_doe","password":"password123","email":"john@example.com"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		repoMock.On("GetRoleByName", mock.Anything, "user").
			Return(models.Role{}, pgx.ErrNoRows).Once()

		repoMock.On("CreateRole", mock.Anything, "user", "Default user role").
			Return("role-id-123", nil).Once()

		repoMock.On("CreateUser", mock.Anything, "john_doe", mock.AnythingOfType("string"), "role-id-123").
			Return("user-id-123", nil).Once()

		repoMock.On("CreateProfile", mock.Anything, "user-id-123", "", "", "john@example.com", "", "", "").
			Return("profile-id-123", nil).Once()

		_ = handler.Register(c)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "User registered successfully")
	})

	t.Run("Error during role creation", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"john_doe","password":"password123","email":"john@example.com"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		repoMock.On("GetRoleByName", mock.Anything, "user").
			Return(models.Role{}, pgx.ErrNoRows).Once()

		repoMock.On("CreateRole", mock.Anything, "user", "Default user role").
			Return("", errors.New("database error")).Once()

		_ = handler.Register(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to create role")
	})

	t.Run("Create user error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"john_doe","password":"password123","email":"john@example.com"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		repoMock.On("GetRoleByName", mock.Anything, "user").
			Return(models.Role{}, pgx.ErrNoRows).Once()

		repoMock.On("CreateRole", mock.Anything, "user", "Default user role").
			Return("role-id-123", nil).Once()

		repoMock.On("CreateUser", mock.Anything, "john_doe", mock.Anything, "role-id-123").
			Return("", pgx.ErrClosedPool).Once()

		_ = handler.Register(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to create user")
	})

	t.Run("Create profile error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username":"john_doe","password":"password123","email":"john@example.com"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		repoMock.On("GetRoleByName", mock.Anything, "user").
			Return(models.Role{}, pgx.ErrNoRows).Once()

		repoMock.On("CreateRole", mock.Anything, "user", "Default user role").
			Return("role-id-123", nil).Once()

		repoMock.On("CreateUser", mock.Anything, "john_doe", mock.AnythingOfType("string"), "role-id-123").
			Return("user-id-123", nil).Once()

		repoMock.On("CreateProfile", mock.Anything, "user-id-123", "", "", "john@example.com", "", "", "").
			Return("", pgx.ErrDeadConn).Once()

		_ = handler.Register(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to create profile")
	})
}

func TestLogin(t *testing.T) {
	repoMock := new(mocks.MockRepository)
	handler := NewHandlers(repoMock, "secret")

	t.Run("Successful login", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"john_doe","password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		passwordHash, _ := bcrypt.GenerateFromPassword([]byte("john_doe"+"password123"), bcrypt.DefaultCost)
		repoMock.On("GetUserByUsername", mock.Anything, "john_doe").
			Return(models.User{ID: "user-id-123", PasswordHash: string(passwordHash)}, nil).Once()

		_ = handler.Login(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Authentication successful")
	})

	t.Run("Invalid username", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"john_doe","password":"password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		repoMock.On("GetUserByUsername", mock.Anything, "john_doe").
			Return(models.User{}, pgx.ErrNoRows).Once()

		_ = handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid username")
	})

	t.Run("Invalid password", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"john_doe","password":"wrongpassword"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		repoMock.On("GetUserByUsername", mock.Anything, "john_doe").
			Return(models.User{ID: "user-id-123", PasswordHash: "hashedpassword"}, nil).Once()

		_ = handler.Login(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid password")
	})
}

func TestProfile(t *testing.T) {
	repoMock := new(mocks.MockRepository)
	handler := NewHandlers(repoMock, "secret")

	t.Run("Successful profile retrieval", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/profile", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user-id-123")

		repoMock.On("GetProfileByUserID", mock.Anything, "user-id-123").
			Return(models.UserProfile{
				ID:          "profile-id-123",
				UserID:      "user-id-123",
				FirstName:   "John",
				LastName:    "Doe",
				Email:       "john@example.com",
				Birthdate:   time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				PhoneNumber: "123456789",
				Bio:         "Bio text",
			}, nil).Once()

		_ = handler.Profile(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "John")
	})

	t.Run("Profile not found", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/profile", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user-id-123")

		repoMock.On("GetProfileByUserID", mock.Anything, "user-id-123").
			Return(models.UserProfile{}, pgx.ErrNoRows).Once()

		_ = handler.Profile(c)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Profile not found")
	})

	t.Run("Get profile error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/profile", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user-id-123")

		repoMock.On("GetProfileByUserID", mock.Anything, "user-id-123").
			Return(models.UserProfile{}, pgx.ErrDeadConn).Once()

		_ = handler.Profile(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to fetch profile")
	})
}

func TestUpdateProfile(t *testing.T) {
	repoMock := new(mocks.MockRepository)
	handler := NewHandlers(repoMock, "secret")

	t.Run("Successful profile update", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(`{"first_name":"John","last_name":"Doe","email":"john@example.com","birthdate":"1990-01-01","phone_number":"123456789","bio":"Updated bio"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user-id-123")

		repoMock.On("GetProfileByUserID", mock.Anything, "user-id-123").
			Return(models.UserProfile{
				ID:          "profile-id-123",
				UserID:      "user-id-123",
				FirstName:   "",
				LastName:    "",
				Email:       "",
				Birthdate:   time.Time{},
				PhoneNumber: "",
				Bio:         "",
			}, nil).Once()

		repoMock.On("UpdateProfile", mock.Anything, "user-id-123", "John", "Doe", "john@example.com", "123456789", "Updated bio", "1990-01-01").
			Return(nil).Once()

		_ = handler.UpdateProfile(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Profile updated successfully")
	})

	t.Run("Profile not found", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(`{"first_name":"John","last_name":"Doe","email":"john@example.com","birthdate":"1990-01-01","phone_number":"123456789","bio":"Updated bio"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user-id-123")

		repoMock.On("GetProfileByUserID", mock.Anything, "user-id-123").
			Return(models.UserProfile{}, pgx.ErrNoRows).Once()

		_ = handler.UpdateProfile(c)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "Profile not found")
	})

	t.Run("Get profile error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(`{"first_name":"John","last_name":"Doe","email":"john@example.com","birthdate":"1990-01-01","phone_number":"123456789","bio":"Updated bio"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user-id-123")

		repoMock.On("GetProfileByUserID", mock.Anything, "user-id-123").
			Return(models.UserProfile{}, pgx.ErrDeadConn).Once()

		_ = handler.UpdateProfile(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to fetch profile")
	})

	t.Run("Successful profile update", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(`{"first_name":"John","last_name":"Doe","email":"john@example.com","birthdate":"1990-01-01","phone_number":"123456789","bio":"Updated bio"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user_id", "user-id-123")

		repoMock.On("GetProfileByUserID", mock.Anything, "user-id-123").
			Return(models.UserProfile{
				ID:          "profile-id-123",
				UserID:      "user-id-123",
				FirstName:   "",
				LastName:    "",
				Email:       "",
				Birthdate:   time.Time{},
				PhoneNumber: "",
				Bio:         "",
			}, nil).Once()

		repoMock.On("UpdateProfile", mock.Anything, "user-id-123", "John", "Doe", "john@example.com", "123456789", "Updated bio", "1990-01-01").
			Return(pgx.ErrConnBusy).Once()

		_ = handler.UpdateProfile(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Failed to update profile")
	})
}
