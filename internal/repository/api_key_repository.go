package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/notification-center/internal/models"
)

const (
	// APIKeyPrefix is the prefix for generated API keys.
	APIKeyPrefix = "nc_live_"

	// APIKeyByteLength is the number of random bytes in the key.
	APIKeyByteLength = 32
)

// APIKeyRepository handles API key database operations.
type APIKeyRepository struct {
	db *pgxpool.Pool
}

// NewAPIKeyRepository creates a new API key repository.
func NewAPIKeyRepository(db *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// Create creates a new API key and returns the full key (only shown once).
func (r *APIKeyRepository) Create(ctx context.Context, apiKey *models.APIKey) (string, error) {
	// Generate random key
	bytes := make([]byte, APIKeyByteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	fullKey := APIKeyPrefix + hex.EncodeToString(bytes)
	keyPrefix := fullKey[:len(APIKeyPrefix)+8]
	keyHash := hashKey(fullKey)

	if apiKey.ID == uuid.Nil {
		apiKey.ID = uuid.New()
	}
	apiKey.KeyPrefix = keyPrefix
	apiKey.KeyHash = keyHash
	apiKey.IsActive = true
	apiKey.CreatedAt = time.Now()

	query := `
		INSERT INTO api_keys (
			id, project_id, user_id, name, key_prefix, key_hash, scopes,
			expires_at, is_active, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(ctx, query,
		apiKey.ID, apiKey.ProjectID, apiKey.UserID, apiKey.Name,
		apiKey.KeyPrefix, apiKey.KeyHash, apiKey.Scopes,
		apiKey.ExpiresAt, apiKey.IsActive, apiKey.CreatedAt,
	)

	if err != nil {
		return "", err
	}

	return fullKey, nil
}

// GetByID retrieves an API key by ID.
func (r *APIKeyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.APIKey, error) {
	var key models.APIKey

	query := `
		SELECT id, project_id, user_id, name, key_prefix, scopes,
		       last_used_at, expires_at, is_active, created_at, revoked_at
		FROM api_keys
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&key.ID, &key.ProjectID, &key.UserID, &key.Name, &key.KeyPrefix,
		&key.Scopes, &key.LastUsedAt, &key.ExpiresAt, &key.IsActive,
		&key.CreatedAt, &key.RevokedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	return &key, err
}

// GetByHash retrieves an API key by its hash.
func (r *APIKeyRepository) GetByHash(ctx context.Context, hash string) (*models.APIKey, error) {
	var key models.APIKey

	query := `
		SELECT id, project_id, user_id, name, key_prefix, scopes,
		       last_used_at, expires_at, is_active, created_at, revoked_at
		FROM api_keys
		WHERE key_hash = $1 AND is_active = true
	`

	err := r.db.QueryRow(ctx, query, hash).Scan(
		&key.ID, &key.ProjectID, &key.UserID, &key.Name, &key.KeyPrefix,
		&key.Scopes, &key.LastUsedAt, &key.ExpiresAt, &key.IsActive,
		&key.CreatedAt, &key.RevokedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	return &key, err
}

// ListByProject retrieves all API keys for a project.
func (r *APIKeyRepository) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]models.APIKey, int, error) {
	if limit <= 0 {
		limit = 20
	}

	// Get total count
	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM api_keys WHERE project_id = $1", projectID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get keys
	query := `
		SELECT id, project_id, user_id, name, key_prefix, scopes,
		       last_used_at, expires_at, is_active, created_at, revoked_at
		FROM api_keys
		WHERE project_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var keys []models.APIKey
	for rows.Next() {
		var k models.APIKey
		err := rows.Scan(
			&k.ID, &k.ProjectID, &k.UserID, &k.Name, &k.KeyPrefix,
			&k.Scopes, &k.LastUsedAt, &k.ExpiresAt, &k.IsActive,
			&k.CreatedAt, &k.RevokedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		keys = append(keys, k)
	}

	return keys, total, nil
}

// ListByUser retrieves all API keys for a user.
func (r *APIKeyRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.APIKey, int, error) {
	if limit <= 0 {
		limit = 20
	}

	// Get total count
	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM api_keys WHERE user_id = $1", userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get keys
	query := `
		SELECT id, project_id, user_id, name, key_prefix, scopes,
		       last_used_at, expires_at, is_active, created_at, revoked_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var keys []models.APIKey
	for rows.Next() {
		var k models.APIKey
		err := rows.Scan(
			&k.ID, &k.ProjectID, &k.UserID, &k.Name, &k.KeyPrefix,
			&k.Scopes, &k.LastUsedAt, &k.ExpiresAt, &k.IsActive,
			&k.CreatedAt, &k.RevokedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		keys = append(keys, k)
	}

	return keys, total, nil
}

// Revoke revokes an API key.
func (r *APIKeyRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE api_keys
		SET is_active = false, revoked_at = $2
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// UpdateLastUsed updates the last_used_at timestamp.
func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE api_keys SET last_used_at = $2 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, time.Now())
	return err
}

// Delete permanently deletes an API key.
func (r *APIKeyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM api_keys WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// ValidateKey validates an API key and returns the associated key record.
func (r *APIKeyRepository) ValidateKey(ctx context.Context, fullKey string) (*models.APIKey, error) {
	hash := hashKey(fullKey)
	return r.GetByHash(ctx, hash)
}

// hashKey returns the SHA-256 hash of a key.
func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}
