package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/notification-center/internal/model"
)

type DeviceStore struct {
	pool *pgxpool.Pool
}

func NewDeviceStore(pool *pgxpool.Pool) *DeviceStore {
	return &DeviceStore{pool: pool}
}

func (s *DeviceStore) Register(ctx context.Context, d *model.DeviceToken) error {
	query := `
		INSERT INTO device_tokens (id, project_id, user_id, token, platform, app_version, created_at, last_used_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (project_id, token) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			platform = EXCLUDED.platform,
			app_version = EXCLUDED.app_version,
			last_used_at = EXCLUDED.last_used_at`

	_, err := s.pool.Exec(ctx, query,
		d.ID, d.ProjectID, d.UserID, d.Token, d.Platform, d.AppVersion, d.CreatedAt, d.LastUsedAt,
	)
	if err != nil {
		return fmt.Errorf("register device: %w", err)
	}
	return nil
}

func (s *DeviceStore) GetByUser(ctx context.Context, projectID, userID string) ([]*model.DeviceToken, error) {
	query := `
		SELECT id, project_id, user_id, token, platform, app_version, created_at, last_used_at
		FROM device_tokens
		WHERE project_id = $1 AND user_id = $2
		ORDER BY last_used_at DESC`

	rows, err := s.pool.Query(ctx, query, projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("get devices: %w", err)
	}
	defer rows.Close()

	var devices []*model.DeviceToken
	for rows.Next() {
		var d model.DeviceToken
		if err := rows.Scan(&d.ID, &d.ProjectID, &d.UserID, &d.Token, &d.Platform,
			&d.AppVersion, &d.CreatedAt, &d.LastUsedAt); err != nil {
			return nil, fmt.Errorf("scan device: %w", err)
		}
		devices = append(devices, &d)
	}
	return devices, nil
}

func (s *DeviceStore) Remove(ctx context.Context, projectID, tokenID string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM device_tokens WHERE project_id = $1 AND id = $2`, projectID, tokenID)
	return err
}

func (s *DeviceStore) RemoveByToken(ctx context.Context, projectID, token string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM device_tokens WHERE project_id = $1 AND token = $2`, projectID, token)
	return err
}
