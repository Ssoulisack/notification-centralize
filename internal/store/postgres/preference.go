package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/notification-center/internal/model"
)

type PreferenceStore struct {
	pool *pgxpool.Pool
}

func NewPreferenceStore(pool *pgxpool.Pool) *PreferenceStore {
	return &PreferenceStore{pool: pool}
}

func (s *PreferenceStore) Get(ctx context.Context, projectID, userID string) (*model.UserPreference, error) {
	query := `
		SELECT project_id, user_id, enabled_channels, quiet_start, quiet_end,
			opted_out_events, updated_at
		FROM user_preferences
		WHERE project_id = $1 AND user_id = $2`

	var p model.UserPreference
	err := s.pool.QueryRow(ctx, query, projectID, userID).Scan(
		&p.ProjectID, &p.UserID, &p.EnabledChannels, &p.QuietHoursStart,
		&p.QuietHoursEnd, &p.OptedOutEvents, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get preference for user %s: %w", userID, err)
	}
	return &p, nil
}

func (s *PreferenceStore) Upsert(ctx context.Context, pref *model.UserPreference) error {
	query := `
		INSERT INTO user_preferences (project_id, user_id, enabled_channels, quiet_start,
			quiet_end, opted_out_events, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (project_id, user_id) DO UPDATE SET
			enabled_channels = EXCLUDED.enabled_channels,
			quiet_start = EXCLUDED.quiet_start,
			quiet_end = EXCLUDED.quiet_end,
			opted_out_events = EXCLUDED.opted_out_events,
			updated_at = NOW()`

	_, err := s.pool.Exec(ctx, query,
		pref.ProjectID, pref.UserID, pref.EnabledChannels, pref.QuietHoursStart,
		pref.QuietHoursEnd, pref.OptedOutEvents,
	)
	if err != nil {
		return fmt.Errorf("upsert preference: %w", err)
	}
	return nil
}
