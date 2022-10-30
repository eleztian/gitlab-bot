package plugin

import (
	"errors"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"

	"github.com/eleztian/gitlab-bot/common"
)

var Plugins = sync.Map{}

func LoadPlugin(path string) (name string, handler common.Handler, err error) {
	_, fileName := filepath.Split(path)
	name = strings.TrimRight(fileName, ".so")
	p, err := plugin.Open(path)
	if err != nil {
		return "", nil, err
	}
	h, err := p.Lookup("Handler")
	if err != nil {
		return "", nil, err
	}

	var ok bool
	handler, ok = h.(common.Handler)
	if !ok {
		return "", nil, errors.New("invalid Handler")
	}

	return name, handler, nil
}

func watchPluginDir(dir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = watcher.Add(dir)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op.Has(fsnotify.Rename | fsnotify.Create | fsnotify.Write) {
					if strings.HasSuffix(event.Name, ".so") {
						pluginName, pluginHandler, err := LoadPlugin(event.Name)
						if err != nil {
							logrus.WithError(err).WithField("Plugin", event.Name).Error("failed to load plugin")
							continue
						}
						Plugins.Store(pluginName, pluginHandler)
					}
				}

			case err := <-watcher.Errors:
				logrus.WithError(err).Error("watch error")
			}
		}
	}()
	return nil
}

func LoadAndWatchPlugins(dir string) error {
	pluginFs, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range pluginFs {
		if !strings.HasSuffix(f.Name(), ".so") {
			continue
		}
		pluginFile := filepath.Join(dir, f.Name())

		pluginName, pluginHandler, err := LoadPlugin(pluginFile)
		if err != nil {
			logrus.WithError(err).WithField("Plugin", pluginFile).Error("failed to load plugin")
			continue
		}
		Plugins.Store(pluginName, pluginHandler)
	}

	return watchPluginDir(dir)
}
