package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/notification-center/internal/middleware"
	"github.com/your-org/notification-center/internal/models"
)

// UserSyncService handles syncing users from Keycloak to the local database.
type UserSyncService struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

// NewUserSyncService creates a new user sync service.
func NewUserSyncService(db *pgxpool.Pool, logger *slog.Logger) *UserSyncService {
	return &UserSyncService{
		db:     db,
		logger: logger,
	}
}

// SyncUser creates or updates a user from Keycloak claims.
func (s *UserSyncService) SyncUser(ctx context.Context, claims *middleware.KeycloakClaims) (*models.User, error) {
	// Check if user exists
	var user models.User

	query := `
		SELECT id, keycloak_id, email, username, first_name, last_name,
		       avatar_url, email_verified, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE keycloak_id = $1
	`

	err := s.db.QueryRow(ctx, query, claims.Subject).Scan(
		&user.ID, &user.KeycloakID, &user.Email, &user.Username,
		&user.FirstName, &user.LastName, &user.AvatarURL,
		&user.EmailVerified, &user.IsActive, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		// User doesn't exist, create new one
		return s.createUser(ctx, claims)
	}

	// Update existing user
	return s.updateUser(ctx, &user, claims)
}

// createUser creates a new user from Keycloak claims.
func (s *UserSyncService) createUser(ctx context.Context, claims *middleware.KeycloakClaims) (*models.User, error) {
	user := &models.User{
		ID:            uuid.New(),
		KeycloakID:    claims.Subject,
		Email:         claims.Email,
		Username:      claims.PreferredUsername,
		FirstName:     claims.GivenName,
		LastName:      claims.FamilyName,
		EmailVerified: claims.EmailVerified,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	now := time.Now()
	user.LastLoginAt = &now

	query := `
		INSERT INTO users (
			id, keycloak_id, email, username, first_name, last_name,
			email_verified, is_active, last_login_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(ctx, query,
		user.ID, user.KeycloakID, user.Email, user.Username,
		user.FirstName, user.LastName, user.EmailVerified,
		user.IsActive, user.LastLoginAt, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		s.logger.Error("failed to create user", "error", err)
		return nil, err
	}

	s.logger.Info("created new user", "user_id", user.ID, "email", user.Email)

	return user, nil
}

// updateUser updates an existing user from Keycloak claims.
func (s *UserSyncService) updateUser(ctx context.Context, user *models.User, claims *middleware.KeycloakClaims) (*models.User, error) {
	user.Email = claims.Email
	user.Username = claims.PreferredUsername
	user.FirstName = claims.GivenName
	user.LastName = claims.FamilyName
	user.EmailVerified = claims.EmailVerified

	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now

	query := `
		UPDATE users
		SET email = $2, username = $3, first_name = $4, last_name = $5,
		    email_verified = $6, last_login_at = $7, updated_at = $8
		WHERE id = $1
		RETURNING updated_at
	`

	err := s.db.QueryRow(ctx, query,
		user.ID, user.Email, user.Username, user.FirstName, user.LastName,
		user.EmailVerified, user.LastLoginAt, user.UpdatedAt,
	).Scan(&user.UpdatedAt)

	if err != nil {
		s.logger.Error("failed to update user", "error", err)
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by their database ID.
func (s *UserSyncService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, keycloak_id, email, username, first_name, last_name,
		       avatar_url, email_verified, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.KeycloakID, &user.Email, &user.Username,
		&user.FirstName, &user.LastName, &user.AvatarURL,
		&user.EmailVerified, &user.IsActive, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByKeycloakID retrieves a user by their Keycloak ID.
func (s *UserSyncService) GetUserByKeycloakID(ctx context.Context, keycloakID string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, keycloak_id, email, username, first_name, last_name,
		       avatar_url, email_verified, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE keycloak_id = $1
	`

	err := s.db.QueryRow(ctx, query, keycloakID).Scan(
		&user.ID, &user.KeycloakID, &user.Email, &user.Username,
		&user.FirstName, &user.LastName, &user.AvatarURL,
		&user.EmailVerified, &user.IsActive, &user.LastLoginAt,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
