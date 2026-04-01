package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/notification-center/internal/models"
)

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("record not found")

// ProjectRepository handles project database operations.
type ProjectRepository struct {
	db *pgxpool.Pool
}

// NewProjectRepository creates a new project repository.
func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create creates a new project.
func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	query := `
		INSERT INTO projects (id, name, description, slug, settings, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	if project.ID == uuid.Nil {
		project.ID = uuid.New()
	}
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	project.IsActive = true

	_, err := r.db.Exec(ctx, query,
		project.ID, project.Name, project.Description, project.Slug,
		project.Settings, project.IsActive, project.CreatedAt, project.UpdatedAt,
	)

	return err
}

// GetByID retrieves a project by ID.
func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project

	query := `
		SELECT id, name, description, slug, settings, is_active, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Slug,
		&project.Settings, &project.IsActive, &project.CreatedAt, &project.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	return &project, err
}

// GetBySlug retrieves a project by slug.
func (r *ProjectRepository) GetBySlug(ctx context.Context, slug string) (*models.Project, error) {
	var project models.Project

	query := `
		SELECT id, name, description, slug, settings, is_active, created_at, updated_at
		FROM projects
		WHERE slug = $1
	`

	err := r.db.QueryRow(ctx, query, slug).Scan(
		&project.ID, &project.Name, &project.Description, &project.Slug,
		&project.Settings, &project.IsActive, &project.CreatedAt, &project.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	return &project, err
}

// Update updates a project.
func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	query := `
		UPDATE projects
		SET name = $2, description = $3, settings = $4, is_active = $5, updated_at = $6
		WHERE id = $1
	`

	project.UpdatedAt = time.Now()

	result, err := r.db.Exec(ctx, query,
		project.ID, project.Name, project.Description, project.Settings,
		project.IsActive, project.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete soft-deletes a project by setting is_active to false.
func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE projects SET is_active = false, updated_at = $2 WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, time.Now())
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// List retrieves all active projects.
func (r *ProjectRepository) List(ctx context.Context, limit, offset int) ([]models.Project, int, error) {
	if limit <= 0 {
		limit = 20
	}

	// Get total count
	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM projects WHERE is_active = true").Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get projects
	query := `
		SELECT id, name, description, slug, settings, is_active, created_at, updated_at
		FROM projects
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Slug,
			&p.Settings, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}

	return projects, total, nil
}

// ListByUser retrieves projects that a user is a member of.
func (r *ProjectRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Project, int, error) {
	if limit <= 0 {
		limit = 20
	}

	// Get total count
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM projects p
		JOIN project_members pm ON pm.project_id = p.id
		WHERE pm.user_id = $1 AND p.is_active = true
	`
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get projects
	query := `
		SELECT p.id, p.name, p.description, p.slug, p.settings, p.is_active, p.created_at, p.updated_at
		FROM projects p
		JOIN project_members pm ON pm.project_id = p.id
		WHERE pm.user_id = $1 AND p.is_active = true
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Slug,
			&p.Settings, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}

	return projects, total, nil
}

// AddMember adds a user to a project with a role.
func (r *ProjectRepository) AddMember(ctx context.Context, member *models.ProjectMember) error {
	query := `
		INSERT INTO project_members (id, project_id, user_id, role_id, joined_at, invited_by)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	member.JoinedAt = time.Now()

	_, err := r.db.Exec(ctx, query,
		member.ID, member.ProjectID, member.UserID, member.RoleID,
		member.JoinedAt, member.InvitedBy,
	)

	return err
}

// GetMember retrieves a project member.
func (r *ProjectRepository) GetMember(ctx context.Context, projectID, userID uuid.UUID) (*models.ProjectMember, error) {
	var member models.ProjectMember

	query := `
		SELECT pm.id, pm.project_id, pm.user_id, pm.role_id, pm.joined_at, pm.invited_by,
		       r.name as role_name, r.permissions
		FROM project_members pm
		JOIN roles r ON r.id = pm.role_id
		WHERE pm.project_id = $1 AND pm.user_id = $2
	`

	var roleName models.MemberRole
	var permissions []byte

	err := r.db.QueryRow(ctx, query, projectID, userID).Scan(
		&member.ID, &member.ProjectID, &member.UserID, &member.RoleID,
		&member.JoinedAt, &member.InvitedBy, &roleName, &permissions,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	member.Role = &models.Role{
		ID:   member.RoleID,
		Name: roleName,
	}

	return &member, nil
}

// UpdateMemberRole updates a member's role.
func (r *ProjectRepository) UpdateMemberRole(ctx context.Context, projectID, userID, roleID uuid.UUID) error {
	query := `UPDATE project_members SET role_id = $3 WHERE project_id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, query, projectID, userID, roleID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// RemoveMember removes a user from a project.
func (r *ProjectRepository) RemoveMember(ctx context.Context, projectID, userID uuid.UUID) error {
	query := `DELETE FROM project_members WHERE project_id = $1 AND user_id = $2`

	result, err := r.db.Exec(ctx, query, projectID, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// ListMembers retrieves all members of a project.
func (r *ProjectRepository) ListMembers(ctx context.Context, projectID uuid.UUID, limit, offset int) ([]models.ProjectMember, int, error) {
	if limit <= 0 {
		limit = 50
	}

	// Get total count
	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM project_members WHERE project_id = $1", projectID).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get members
	query := `
		SELECT pm.id, pm.project_id, pm.user_id, pm.role_id, pm.joined_at, pm.invited_by,
		       u.id, u.keycloak_id, u.email, u.username, u.first_name, u.last_name,
		       r.id, r.name
		FROM project_members pm
		JOIN users u ON u.id = pm.user_id
		JOIN roles r ON r.id = pm.role_id
		WHERE pm.project_id = $1
		ORDER BY pm.joined_at ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, projectID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var members []models.ProjectMember
	for rows.Next() {
		var m models.ProjectMember
		var u models.User
		var role models.Role

		err := rows.Scan(
			&m.ID, &m.ProjectID, &m.UserID, &m.RoleID, &m.JoinedAt, &m.InvitedBy,
			&u.ID, &u.KeycloakID, &u.Email, &u.Username, &u.FirstName, &u.LastName,
			&role.ID, &role.Name,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan member: %w", err)
		}

		m.User = &u
		m.Role = &role
		members = append(members, m)
	}

	return members, total, nil
}

// GetRoleByName retrieves a role by its name.
func (r *ProjectRepository) GetRoleByName(ctx context.Context, name models.MemberRole) (*models.Role, error) {
	var role models.Role

	query := `SELECT id, name, permissions, created_at FROM roles WHERE name = $1`

	var permissions []byte
	err := r.db.QueryRow(ctx, query, name).Scan(
		&role.ID, &role.Name, &permissions, &role.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	return &role, err
}
