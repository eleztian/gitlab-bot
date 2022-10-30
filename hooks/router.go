package hooks

import (
	"context"

	"github.com/xanzy/go-gitlab"
)

type Router interface {
	Match(eventType gitlab.EventType, event interface{}) bool
}

type Handler func(ctx context.Context, event interface{})

type eventTypeRouter struct {
	EventTypes []gitlab.EventType
}

func NewEventTypeRouter(eventTypes ...gitlab.EventType) Router {
	return &eventTypeRouter{
		EventTypes: eventTypes,
	}
}

func (r *eventTypeRouter) Match(eventType gitlab.EventType, event interface{}) bool {
	if len(r.EventTypes) == 0 {
		return true
	}
	for _, tp := range r.EventTypes {
		if tp == eventType {
			return true
		}
	}

	return false
}
