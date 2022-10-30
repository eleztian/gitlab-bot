package common

import (
	"context"

	"github.com/xanzy/go-gitlab"
)

type Event struct {
	EventType gitlab.EventType
	Event     interface{}
}

type Handler interface {
	Do(ctx context.Context, client *gitlab.Client, event Event) error
}
