package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nanoservices/users_service/models"
)

type RepositoryInt interface {
	CreateRole(ctx context.Context, name, description string) (string, error)
	GetRoleByName(ctx context.Context, name string) (models.Role, error)
	CreateUser(ctx context.Context, username, passwordHash, roleID string) (string, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	CreateProfile(ctx context.Context, userID, firstName, lastName, email, birthdate, phoneNumber, bio string) (string, error)
	GetProfileByUserID(ctx context.Context, userID string) (models.UserProfile, error)
	UpdateProfile(ctx context.Context, userID, firstName, lastName, email, phoneNumber, bio, birthdate string) error
}

type Repository struct {
	pool DB
}

func NewRepository(pool DB) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateRole(ctx context.Context, name, description string) (string, error) {
	query := `
        INSERT INTO roles (id, name, description, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
		`
	id := uuid.New().String()
	err := r.pool.QueryRow(ctx, query, id, name, description, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) GetRoleByName(ctx context.Context, name string) (models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1`
	row := r.pool.QueryRow(ctx, query, name)

	var role models.Role
	err := row.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return models.Role{}, err
	}
	return role, nil
}

func (r *Repository) CreateUser(ctx context.Context, username, passwordHash, roleID string) (string, error) {
	query := `
        INSERT INTO users (id, role_id, username, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`
	id := uuid.New().String()
	err := r.pool.QueryRow(ctx, query, id, roleID, username, passwordHash, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	query := `SELECT id, role_id, username, password_hash, created_at, updated_at FROM users WHERE username = $1`
	row := r.pool.QueryRow(ctx, query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.RoleID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *Repository) CreateProfile(ctx context.Context, userID, firstName, lastName, email, birthdate, phoneNumber, bio string) (string, error) {
	query := `
        INSERT INTO user_profiles (id, user_id, first_name, last_name, email, birthdate, bio, phone_number, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id`
	id := uuid.New().String()
	birthdateParsed, _ := time.Parse("2006-01-02", birthdate)
	err := r.pool.QueryRow(ctx, query, id, userID, firstName, lastName, email, birthdateParsed, phoneNumber, bio, time.Now()).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) GetProfileByUserID(ctx context.Context, userID string) (models.UserProfile, error) {
	query := `SELECT id, user_id, first_name, last_name, email, birthdate, phone_number, bio, created_at FROM user_profiles WHERE user_id = $1`
	row := r.pool.QueryRow(ctx, query, userID)

	var profile models.UserProfile
	err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
		&profile.Birthdate,
		&profile.PhoneNumber,
		&profile.Bio,
		&profile.CreatedAt,
	)
	if err != nil {
		return models.UserProfile{}, err
	}
	return profile, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, userID, firstName, lastName, email, phoneNumber, bio, birthdate string) error {
	query := `
        UPDATE user_profiles
        SET first_name = $1, last_name = $2, email = $3, birthdate = $4, phone_number = $5, bio = $6, updated_at = NOW()
        WHERE user_id = $7`
	_, err := r.pool.Exec(context.Background(), query, firstName, lastName, email, birthdate, phoneNumber, bio, userID)
	return err
}
