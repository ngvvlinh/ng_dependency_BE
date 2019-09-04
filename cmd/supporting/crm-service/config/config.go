package config

import (
	"os"
	"path/filepath"

	cc "etop.vn/backend/pkg/common/config"
)

type Vtiger struct {
	ServiceURL string `yaml:"service_url"`
	Username   string `yaml:"username"`
	APIKey     string `yaml:"api_key"`
}

type Vht struct {
	ServiceURL string `yaml:"service_url"`
	UserName   string `yaml:"user_name"`
	PassWord   string `yaml:"pass_word"`
}

type Config struct {
	Postgres    cc.Postgres `yaml:"postgres"`
	Redis       cc.Redis    `yaml:"redis"`
	HTTP        cc.HTTP     `yaml:"http"`
	Env         string      `yaml:"env"`
	MappingFile string      `yaml:"mapping_file"`
	Vtiger      Vtiger      `yaml:"vtiger"`
	Vht         Vht         `yaml:"vht"`
}

var exampleMappingFile = filepath.Join(
	os.Getenv("ETOPDIR"),
	"backend/cmd/etop-server/config/field_mapping_example.json",
)

func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		Redis:    cc.DefaultRedis(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8080,
		},
		Env: "dev",
		Vtiger: Vtiger{
			ServiceURL: "http://vtiger/webservice.php",
			Username:   "admin",
			APIKey:     "q5dZOnJYGlmPY2nc",
		},
		Vht: Vht{
			UserName: "5635810cde4c14ebf6a41341f4e68395",
			PassWord: "36828473ce0d87db8cc29798f6b8aa1e",
		},
		MappingFile: exampleMappingFile,
	}
	cfg.Postgres.Database = "etop_crm"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	return cfg, err
}
