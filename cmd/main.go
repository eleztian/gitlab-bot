package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"

	"github.com/eleztian/gitlab-bot/config"
	"github.com/eleztian/gitlab-bot/hooks"
	"github.com/eleztian/gitlab-bot/plugin"
)

var (
	gitlabToken    = os.Getenv("GITLAB_TOKEN")
	gitlabURL      = os.Getenv("GITLAB_URL")
	hookServerAddr = os.Getenv("BOOT_HOOK_ADDR")
)

var client *gitlab.Client

func init() {
	var err error
	client, err = gitlab.NewClient(gitlabToken,
		gitlab.WithBaseURL(gitlabURL))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}
}

var configFilePath = flag.String("c", "conf/config.yml", "config file path")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadFile(*configFilePath)
	if err != nil {
		logrus.WithError(err).Fatalln("failed to load config file")
		return
	}

	err = plugin.LoadAndWatchPlugins(cfg.PluginDir)
	if err != nil {
		logrus.WithError(err).Fatalln("plugin load failed")
		return
	}

	var newHandler = func(cfg *config.Config) http.Handler {
		hookServer := hooks.NewHookServer(client)
		for name, h := range cfg.Routers {
			logrus.WithField("Prefix", name).
				WithField("Types", h.Types).
				WithField("Handler", h.Handler).Info("add route")
			hookServer.AddRouter(name, hooks.NewEventTypeRouter(h.Types...), hooks.PluginHandler(h.Handler))
		}
		return hookServer.Handler(ctx)
	}

	switchHandler := &SwitchHandler{}
	switchHandler.SetHandler(newHandler(cfg))
	err = config.WatchFile(*configFilePath, func(cfg *config.Config) {
		switchHandler.SetHandler(newHandler(cfg))
	})
	if err != nil {
		logrus.WithError(err).Fatalln("failed to watch config file")
		return
	}
	logrus.WithField("HookServerAddr", hookServerAddr).Info("start hook server....")
	err = http.ListenAndServe(hookServerAddr, switchHandler)
	if err != nil {
		logrus.WithError(err).WithField("Addr", hookServerAddr).Fatalln("http listen and server failed")
	}
}

type SwitchHandler struct {
	atomic.Value
}

func (sh *SwitchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h := sh.Load().(http.Handler)
	h.ServeHTTP(writer, request)
}

func (sh *SwitchHandler) SetHandler(h http.Handler) {
	sh.Store(h)
}
