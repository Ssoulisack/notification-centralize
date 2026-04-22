package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/notification-center/internal/domain/models"
)

const (
	APIKeyPrefix     = "nc_live_"
	APIKeyByteLength = 32
)

type APIKeyRepository interface {
	Create(ctx context.Context, apiKey *models.APIKey) (string, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.APIKey, error)
	GetByHash(ctx context.Context, hash string) (*models.APIKey, error)
	ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]models.APIKey, int, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.APIKey, int, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	UpdateLastUsed(ctx context.Context, id uuid.UUID) error
}

type apiKeyRepository struct {
	db *gorm.DB
}

func NewAPIKeyRepository(db *gorm.DB) APIKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) Create(ctx context.Context, apiKey *models.APIKey) (string, error) {
	b := make([]byte, APIKeyByteLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	fullKey := APIKeyPrefix + hex.EncodeToString(b)
	apiKey.KeyPrefix = fullKey[:len(APIKeyPrefix)+8]
	apiKey.KeyHash = hashKey(fullKey)
	apiKey.IsActive = true
	apiKey.CreatedAt = time.Now()
	if apiKey.ID == uuid.Nil {
		apiKey.ID = uuid.New()
	}

	if err := r.db.WithContext(ctx).Create(apiKey).Error; err != nil {
		return "", err
	}
	return fullKey, nil
}

func (r *apiKeyRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.APIKey, error) {
	var key models.APIKey
	err := r.db.WithContext(ctx).First(&key, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &key, err
}

func (r *apiKeyRepository) GetByHash(ctx context.Context, hash string) (*models.APIKey, error) {
	var key models.APIKey
	err := r.db.WithContext(ctx).Where("key_hash = ? AND is_active = true", hash).First(&key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &key, err
}

func (r *apiKeyRepository) ListByProject(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]models.APIKey, int, error) {
	if limit <= 0 {
		limit = 20
	}
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.APIKey{}).Where("project_id = ?", projectID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var keys []models.APIKey
	err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&keys).Error
	return keys, int(total), err
}

func (r *apiKeyRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.APIKey, int, error) {
	if limit <= 0 {
		limit = 20
	}
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.APIKey{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var keys []models.APIKey
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&keys).Error
	return keys, int(total), err
}

func (r *apiKeyRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Model(&models.APIKey{}).Where("id = ?", id).Updates(map[string]any{
		"is_active":  false,
		"revoked_at": time.Now(),
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *apiKeyRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.APIKey{}).Where("id = ?", id).Update("last_used_at", time.Now()).Error
}

func hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}
