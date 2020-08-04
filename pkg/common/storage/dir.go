package storage

import (
	cm "o.o/backend/pkg/common"
)

type Purpose string

type DirConfigs map[Purpose]DirConfig

func (dirs DirConfigs) Get(s Purpose) (DirConfig, bool) {
	if s == "" {
		s = "default"
	}
	dirCfg, ok := dirs[s]
	return dirCfg, ok
}

type DirConfig struct {
	Path      string `yaml:"path"`
	URLPath   string `yaml:"url_path"`
	URLPrefix string `yaml:"url_prefix"`
}

func (dir DirConfig) Validate() error {
	if dir.Path == "" {
		return cm.Errorf(cm.Internal, nil, "empty dir")
	}
	if (dir.URLPrefix == "") != (dir.URLPath == "") {
		return cm.Errorf(cm.Internal, nil, "must provide both url_prefix and url_path (or leave both empty)")
	}
	return nil
}
