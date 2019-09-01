package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/services/crm-service/mapping"
	vhtservice "etop.vn/backend/pkg/services/crm-service/vht/service"
	vtigerservice "etop.vn/backend/pkg/services/crm-service/vtiger/service"
)

type Config struct {
	Postgres    cc.Postgres          `yaml:"postgres"`
	Redis       cc.Redis             `yaml:"redis"`
	HTTP        cc.HTTP              `yaml:"http"`
	Env         string               `yaml:"env"`
	Vtiger      vtigerservice.Config `yaml:"vtiger"`
	Vht         vhtservice.Config    `yaml:vht`
	MappingFile string               `yaml:"mapping_file"`
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
		Vht: vhtservice.Config{
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

// ReadMappingFile read mapping json file for mapping fields between vtiger and etop
func ReadMappingFile(filename string) (configMap mapping.ConfigMap, _ error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &configMap)
	return
}