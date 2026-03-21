package provider

import (
	"context"

	"github.com/your-org/notification-center/internal/model"
)

// Provider is the interface every notification channel must implement.
type Provider interface {
	// Send delivers a notification through this channel.
	Send(ctx context.Context, n *model.Notification) error

	// Channel returns which channel this provider handles.
	Channel() model.Channel

	// Name returns a human-readable provider name (e.g. "twilio", "fcm").
	Name() string
}

// Registry maps channels to their provider implementation.
type Registry struct {
	providers map[model.Channel]Provider
}

// NewRegistry creates an empty provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[model.Channel]Provider),
	}
}

// Register adds a provider for a channel.
func (r *Registry) Register(p Provider) {
	r.providers[p.Channel()] = p
}

// Get returns the provider for a channel, if registered.
func (r *Registry) Get(ch model.Channel) (Provider, bool) {
	p, ok := r.providers[ch]
	return p, ok
}

// Channels returns all registered channels.
func (r *Registry) Channels() []model.Channel {
	channels := make([]model.Channel, 0, len(r.providers))
	for ch := range r.providers {
		channels = append(channels, ch)
	}
	return channels
}
