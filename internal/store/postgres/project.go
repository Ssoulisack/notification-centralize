package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/notification-center/internal/model"
)

type ProjectStore struct {
	pool *pgxpool.Pool
}

func NewProjectStore(pool *pgxpool.Pool) *ProjectStore {
	return &ProjectStore{pool: pool}
}

func (s *ProjectStore) GetByAPIKey(ctx context.Context, apiKey string) (*model.Project, error) {
	query := `
		SELECT id, name, api_key, created_at, updated_at
		FROM projects WHERE api_key = $1`

	var p model.Project
	err := s.pool.QueryRow(ctx, query, apiKey).Scan(
		&p.ID, &p.Name, &p.APIKey, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get project by api key: %w", err)
	}
	return &p, nil
}

func (s *ProjectStore) GetByID(ctx context.Context, id string) (*model.Project, error) {
	query := `
		SELECT id, name, api_key, created_at, updated_at
		FROM projects WHERE id = $1`

	var p model.Project
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.APIKey, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get project %s: %w", id, err)
	}
	return &p, nil
}

func (s *ProjectStore) Create(ctx context.Context, p *model.Project) error {
	query := `
		INSERT INTO projects (id, name, api_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := s.pool.Exec(ctx, query, p.ID, p.Name, p.APIKey, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

func (s *ProjectStore) List(ctx context.Context) ([]*model.Project, error) {
	query := `
		SELECT id, name, api_key, created_at, updated_at
		FROM projects ORDER BY name`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []*model.Project
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.APIKey, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, &p)
	}
	return projects, nil
}
