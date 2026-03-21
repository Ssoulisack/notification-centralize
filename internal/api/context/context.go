package apicontext

import (
	"context"
	"errors"

	"github.com/your-org/notification-center/internal/model"
)

type contextKey string

const projectKey contextKey = "project"

var ErrNoProject = errors.New("no project in context")

// WithProject returns a new context with the project attached.
func WithProject(ctx context.Context, project *model.Project) context.Context {
	return context.WithValue(ctx, projectKey, project)
}

// ProjectFromContext retrieves the project from context.
func ProjectFromContext(ctx context.Context) (*model.Project, error) {
	project, ok := ctx.Value(projectKey).(*model.Project)
	if !ok || project == nil {
		return nil, ErrNoProject
	}
	return project, nil
}

// ProjectIDFromContext retrieves just the project ID from context.
func ProjectIDFromContext(ctx context.Context) (string, error) {
	project, err := ProjectFromContext(ctx)
	if err != nil {
		return "", err
	}
	return project.ID, nil
}
