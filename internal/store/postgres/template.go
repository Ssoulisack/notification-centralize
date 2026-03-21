package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/notification-center/internal/model"
)

type TemplateStore struct {
	pool *pgxpool.Pool
}

func NewTemplateStore(pool *pgxpool.Pool) *TemplateStore {
	return &TemplateStore{pool: pool}
}

func (s *TemplateStore) GetByID(ctx context.Context, projectID, id string) (*model.Template, error) {
	query := `
		SELECT id, project_id, name, channel, subject_template, body_template, created_at, updated_at
		FROM notification_templates WHERE project_id = $1 AND id = $2`

	var t model.Template
	err := s.pool.QueryRow(ctx, query, projectID, id).Scan(
		&t.ID, &t.ProjectID, &t.Name, &t.Channel, &t.SubjectTemplate,
		&t.BodyTemplate, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get template %s: %w", id, err)
	}
	return &t, nil
}

func (s *TemplateStore) GetByName(ctx context.Context, projectID, name string, channel model.Channel) (*model.Template, error) {
	query := `
		SELECT id, project_id, name, channel, subject_template, body_template, created_at, updated_at
		FROM notification_templates WHERE project_id = $1 AND name = $2 AND channel = $3`

	var t model.Template
	err := s.pool.QueryRow(ctx, query, projectID, name, channel).Scan(
		&t.ID, &t.ProjectID, &t.Name, &t.Channel, &t.SubjectTemplate,
		&t.BodyTemplate, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get template %s/%s: %w", name, channel, err)
	}
	return &t, nil
}

func (s *TemplateStore) Create(ctx context.Context, t *model.Template) error {
	query := `
		INSERT INTO notification_templates (id, project_id, name, channel, subject_template, body_template, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`

	_, err := s.pool.Exec(ctx, query, t.ID, t.ProjectID, t.Name, t.Channel, t.SubjectTemplate, t.BodyTemplate)
	return err
}

func (s *TemplateStore) Update(ctx context.Context, t *model.Template) error {
	query := `
		UPDATE notification_templates
		SET subject_template = $3, body_template = $4, updated_at = NOW()
		WHERE project_id = $1 AND id = $2`

	_, err := s.pool.Exec(ctx, query, t.ProjectID, t.ID, t.SubjectTemplate, t.BodyTemplate)
	return err
}

func (s *TemplateStore) List(ctx context.Context, projectID string) ([]*model.Template, error) {
	query := `SELECT id, project_id, name, channel, subject_template, body_template, created_at, updated_at
		FROM notification_templates WHERE project_id = $1 ORDER BY name`

	rows, err := s.pool.Query(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*model.Template
	for rows.Next() {
		var t model.Template
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Name, &t.Channel, &t.SubjectTemplate,
			&t.BodyTemplate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, &t)
	}
	return templates, nil
}
