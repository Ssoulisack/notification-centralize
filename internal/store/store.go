package store

import (
	"context"

	"github.com/your-org/notification-center/internal/model"
)

// ProjectStore manages projects (tenants).
type ProjectStore interface {
	GetByAPIKey(ctx context.Context, apiKey string) (*model.Project, error)
	GetByID(ctx context.Context, id string) (*model.Project, error)
	Create(ctx context.Context, p *model.Project) error
	List(ctx context.Context) ([]*model.Project, error)
}

// NotificationStore persists notification logs.
type NotificationStore interface {
	Create(ctx context.Context, n *model.Notification) error
	GetByID(ctx context.Context, projectID, id string) (*model.Notification, error)
	UpdateStatus(ctx context.Context, projectID, id string, status model.Status, errMsg string) error
	ListByUser(ctx context.Context, projectID, userID string, limit, offset int) ([]*model.Notification, error)
}

// DeviceStore manages push notification device tokens.
type DeviceStore interface {
	Register(ctx context.Context, d *model.DeviceToken) error
	GetByUser(ctx context.Context, projectID, userID string) ([]*model.DeviceToken, error)
	Remove(ctx context.Context, projectID, tokenID string) error
	RemoveByToken(ctx context.Context, projectID, token string) error
}

// PreferenceStore manages user notification preferences.
type PreferenceStore interface {
	Get(ctx context.Context, projectID, userID string) (*model.UserPreference, error)
	Upsert(ctx context.Context, pref *model.UserPreference) error
}

// TemplateStore manages notification templates.
type TemplateStore interface {
	GetByID(ctx context.Context, projectID, id string) (*model.Template, error)
	GetByName(ctx context.Context, projectID, name string, channel model.Channel) (*model.Template, error)
	Create(ctx context.Context, t *model.Template) error
	Update(ctx context.Context, t *model.Template) error
	List(ctx context.Context, projectID string) ([]*model.Template, error)
}
