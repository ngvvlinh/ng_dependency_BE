package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/services/crm-service/mapping"
	vtigerservice "etop.vn/backend/pkg/services/crm-service/vtiger/service"
)

type Config struct {
	Postgres cc.Postgres          `yaml:"postgres"`
	Redis    cc.Redis             `yaml:"redis"`
	HTTP     cc.HTTP              `yaml:"http"`
	Env      string               `yaml:"env"`
	Vtiger   vtigerservice.Config `yaml:"vtiger"`

	MappingFile string `yaml:"mapping_file"`
}

var exampleMappingFile = filepath.Join(
	os.Getenv("ETOPDIR"),
	"backend/cmd/supporting/crm-service/config/field_mapping_example.json",
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
		Vtiger: vtigerservice.Config{
			ServiceURL: "http://vtiger/webservice.php",
			Username:   "admin",
			APIKey:     "q5dZOnJYGlmPY2nc",
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

// ReadMappingFile read mapping json file for mapping fields between vtiger and etop
func ReadMappingFile(filename string) (configMap mapping.ConfigMap, _ error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &configMap)
	return
}
