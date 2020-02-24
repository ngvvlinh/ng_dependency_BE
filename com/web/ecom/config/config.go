package config

import cc "etop.vn/backend/pkg/common/config"

type Config struct {
	HTTP     cc.HTTP `yaml:"http"`
	MainSite string  `yaml:"main_site"`
}
