package config

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Addr      string `json:"addr" yaml:"addr"`
	PluginDir string `json:"plugin_dir" yaml:"plugin_dir"`
	Routers   map[string]struct {
		Types   []gitlab.EventType `json:"types" yaml:"types"`
		Handler string             `json:"handler" yaml:"handler"`
	}
}

func LoadFile(filePath string) (*Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	res := &Config{
		Addr:      "0.0.0.0:7890",
		PluginDir: "plugins",
		Routers: map[string]struct {
			Types   []gitlab.EventType `json:"types" yaml:"types"`
			Handler string             `json:"handler" yaml:"handler"`
		}{},
	}
	err = yaml.Unmarshal(content, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func WatchFile(filepath string, handler func(cfg *Config)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = watcher.Add(filepath)
	if err != nil {
		return err
	}

	go func() {
		for e := range watcher.Events {
			if !e.Has(fsnotify.Write) {
				continue
			}
			cfg, err := LoadFile(filepath)
			if err != nil {
				logrus.WithError(err).Error("failed to reload config file")
			}
			handler(cfg)
		}
	}()

	return nil
}

func LoadAndWatchFile(filepath string, handler func(cfg *Config)) (*Config, error) {
	cfg, err := LoadFile(filepath)
	if err != nil {
		return nil, err
	}
	return cfg, WatchFile(filepath, handler)
}
