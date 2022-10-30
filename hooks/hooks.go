package hooks

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"

	"github.com/eleztian/gitlab-bot/common"
	"github.com/eleztian/gitlab-bot/plugin"
)

type HookServer struct {
	ctx     context.Context
	client  *gitlab.Client
	routers map[string]map[Router]common.Handler
}

func (h *HookServer) Handler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		h.MatchAndRoute(ctx, writer, request)
	})
}

func NewHookServer(client *gitlab.Client) *HookServer {
	return &HookServer{
		client:  client,
		routers: map[string]map[Router]common.Handler{},
	}
}

func (h *HookServer) AddRouter(name string, router Router, handler common.Handler) {
	if _, ok := h.routers[name]; !ok {
		h.routers[name] = map[Router]common.Handler{}
	}
	h.routers[name][router] = handler
}

func (h *HookServer) MatchAndRoute(ctx context.Context, writer http.ResponseWriter, request *http.Request) {
	content, _ := io.ReadAll(request.Body)
	eventType := gitlab.HookEventType(request)
	event, err := gitlab.ParseHook(eventType, content)
	if err != nil {
		logrus.WithError(err).Error("failed to parse request event")
		return
	}

	for route, handlers := range h.routers {
		if request.URL.Path == route {
			for r, handler := range handlers {
				if r.Match(eventType, event) {
					err = handler.Do(ctx, h.client, common.Event{
						EventType: eventType,
						Event:     event,
					})
					if err != nil {
						logrus.WithError(err).Error("handler event failed")
					}
				}
			}
		}
	}
}

func PluginHandler(name string) common.Handler {
	v, ok := plugin.Plugins.Load(name)
	if !ok {
		return common.Handler(NotFoundHandler{})
	}
	return v.(common.Handler)
}

type NotFoundHandler struct{}

func (n NotFoundHandler) Do(ctx context.Context, client *gitlab.Client, event common.Event) error {
	return errors.New("ot found handler")
}
