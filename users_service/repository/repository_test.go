package repository

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nanoservices/users_service/mocks"
	"github.com/nanoservices/users_service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRole(t *testing.T) {
	dbMock := new(mocks.DBMock)
	repo := NewRepository(dbMock)
	ctx := context.Background()
	rowMock := new(mocks.PgxRowMock)

	t.Run("Successful registration", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = "generated-id-123"
			}).Return(nil).Once()

		id, err := repo.CreateRole(ctx, "user", "some description")

		assert.NoError(t, err)
		assert.Equal(t, "generated-id-123", id)
	})

	t.Run("Error during role creation", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Return(assert.AnError).Once()

		id, err := repo.CreateRole(ctx, "user", "some description")

		assert.Error(t, err)
		assert.Empty(t, id)
	})
}

func TestGetRoleByName(t *testing.T) {
	dbMock := new(mocks.DBMock)
	repo := NewRepository(dbMock)
	ctx := context.Background()
	rowMock := new(mocks.PgxRowMock)

	t.Run("Successful role retrieval", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = "role-id-123"
				*args[1].(*string) = "admin"
				*args[2].(*string) = "Administrator role"
				*args[3].(*time.Time) = time.Now()
				*args[4].(*time.Time) = time.Now()
			}).Return(nil).Once()

		role, err := repo.GetRoleByName(ctx, "admin")

		assert.NoError(t, err)
		assert.Equal(t, "role-id-123", role.ID)
		assert.Equal(t, "admin", role.Name)
	})

	t.Run("Error during role retrieval", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError).Once()

		role, err := repo.GetRoleByName(ctx, "admin")

		assert.Error(t, err)
		assert.Equal(t, models.Role{}, role)
	})
}

func TestCreateUser(t *testing.T) {
	dbMock := new(mocks.DBMock)
	repo := NewRepository(dbMock)
	ctx := context.Background()
	rowMock := new(mocks.PgxRowMock)

	t.Run("Successful user creation", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = "user-id-123"
			}).Return(nil).Once()

		id, err := repo.CreateUser(ctx, "john_doe", "hashed-password", "role-id-123")

		assert.NoError(t, err)
		assert.Equal(t, "user-id-123", id)
	})

	t.Run("Error during user creation", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Return(assert.AnError).Once()

		id, err := repo.CreateUser(ctx, "john_doe", "hashed-password", "role-id-123")

		assert.Error(t, err)
		assert.Empty(t, id)
	})
}

func TestGetUserByUsername(t *testing.T) {
	dbMock := new(mocks.DBMock)
	repo := NewRepository(dbMock)
	ctx := context.Background()
	rowMock := new(mocks.PgxRowMock)

	t.Run("Successful user retrieval", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = "user-id-123"
				*args[1].(*string) = "role-id-123"
				*args[2].(*string) = "john_doe"
				*args[3].(*string) = "hashed-password"
				*args[4].(*time.Time) = time.Now()
				*args[5].(*time.Time) = time.Now()
			}).Return(nil).Once()

		user, err := repo.GetUserByUsername(ctx, "john_doe")

		assert.NoError(t, err)
		assert.Equal(t, "user-id-123", user.ID)
		assert.Equal(t, "john_doe", user.Username)
	})

	t.Run("Error during user retrieval", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError).Once()

		user, err := repo.GetUserByUsername(ctx, "john_doe")

		assert.Error(t, err)
		assert.Equal(t, models.User{}, user)
	})
}

func TestCreateProfile(t *testing.T) {
	dbMock := new(mocks.DBMock)
	repo := NewRepository(dbMock)
	ctx := context.Background()
	rowMock := new(mocks.PgxRowMock)

	t.Run("Successful profile creation", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = "profile-id-123"
			}).Return(nil).Once()

		id, err := repo.CreateProfile(ctx, "user-id-123", "John", "Doe", "john@example.com", "1990-01-01", "123456789", "Bio text")

		assert.NoError(t, err)
		assert.Equal(t, "profile-id-123", id)
	})

	t.Run("Error during profile creation", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything).
			Return(assert.AnError).Once()

		id, err := repo.CreateProfile(ctx, "user-id-123", "John", "Doe", "john@example.com", "1990-01-01", "123456789", "Bio text")

		assert.Error(t, err)
		assert.Empty(t, id)
	})
}

func TestGetProfileByUserID(t *testing.T) {
	dbMock := new(mocks.DBMock)
	repo := NewRepository(dbMock)
	ctx := context.Background()
	rowMock := new(mocks.PgxRowMock)

	t.Run("Successful profile retrieval", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				*args[0].(*string) = "profile-id-123"
				*args[1].(*string) = "user-id-123"
				*args[2].(*string) = "John"
				*args[3].(*string) = "Doe"
				*args[4].(*string) = "john@example.com"
				*args[5].(*time.Time) = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
				*args[6].(*string) = "123456789"
				*args[7].(*string) = "Bio text"
				*args[8].(*time.Time) = time.Now()
			}).Return(nil).Once()

		profile, err := repo.GetProfileByUserID(ctx, "user-id-123")

		assert.NoError(t, err)
		assert.Equal(t, "profile-id-123", profile.ID)
		assert.Equal(t, "John", profile.FirstName)
	})

	t.Run("Error during profile retrieval", func(t *testing.T) {
		dbMock.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(rowMock).Once()

		rowMock.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError).Once()

		profile, err := repo.GetProfileByUserID(ctx, "user-id-123")

		assert.Error(t, err)
		assert.Equal(t, models.UserProfile{}, profile)
	})
}

func TestUpdateProfile(t *testing.T) {
	dbMock := new(mocks.DBMock)
	repo := NewRepository(dbMock)
	ctx := context.Background()

	t.Run("Successful profile update", func(t *testing.T) {
		dbMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, nil).Once()

		err := repo.UpdateProfile(ctx, "user-id-123", "John", "Doe", "john@example.com", "123456789", "Updated bio", "1990-01-01")

		assert.NoError(t, err)
	})

	t.Run("Error during profile update", func(t *testing.T) {
		dbMock.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, assert.AnError).Once()

		err := repo.UpdateProfile(ctx, "user-id-123", "John", "Doe", "john@example.com", "123456789", "Updated bio", "1990-01-01")

		assert.Error(t, err)
	})
}
