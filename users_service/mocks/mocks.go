package mocks

import (
	"context"

	"github.com/nanoservices/users_service/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateRole(ctx context.Context, name, description string) (string, error) {
	args := m.Called(ctx, name, description)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockRepository) GetRoleByName(ctx context.Context, name string) (models.Role, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(models.Role), args.Error(1)
}

func (m *MockRepository) CreateUser(ctx context.Context, username, passwordHash, roleID string) (string, error) {
	args := m.Called(ctx, username, passwordHash, roleID)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockRepository) CreateProfile(ctx context.Context, userID, firstName, lastName, email, birthdate, phoneNumber, bio string) (string, error) {
	args := m.Called(ctx, userID, firstName, lastName, email, birthdate, phoneNumber, bio)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockRepository) GetProfileByUserID(ctx context.Context, userID string) (models.UserProfile, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.UserProfile), args.Error(1)
}

func (m *MockRepository) UpdateProfile(ctx context.Context, userID, firstName, lastName, email, phoneNumber, bio, birthdate string) error {
	args := m.Called(ctx, userID, firstName, lastName, email, phoneNumber, bio, birthdate)
	return args.Error(0)
}
