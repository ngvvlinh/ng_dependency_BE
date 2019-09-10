package config

import (
	"os"
	"path/filepath"

	cc "etop.vn/backend/pkg/common/config"
)

type Vtiger struct {
	ServiceURL  string `yaml:"service_url"`
	Username    string `yaml:"username"`
	APIKey      string `yaml:"api_key"`
	MappingFile string `yaml:"mapping_file"`
}

func DefaultVtiger() Vtiger {
	exampleMappingFile := filepath.Join(
		os.Getenv("ETOPDIR"),
		"backend/cmd/etop-server/config/field_mapping_example.json",
	)
	return Vtiger{
		ServiceURL:  "http://vtiger/webservice.php",
		Username:    "admin",
		APIKey:      "",
		MappingFile: exampleMappingFile,
	}
}

func (c *Vtiger) MustLoadEnv(prefix ...string) {
	p := "ET_VTIGER"
	if len(prefix) != 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_SERVICE_URL":  &c.ServiceURL,
		p + "_USERNAME":     &c.Username,
		p + "_API_KEY":      &c.APIKey,
		p + "_MAPPING_FILE": &c.MappingFile,
	}.MustLoad()
}

type Vht struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func DefaultVht() Vht {
	return Vht{
		Username: "",
		Password: "",
	}
}

func (c *Vht) MustLoadEnv(prefix ...string) {
	p := "ET_VHT"
	if len(prefix) != 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_USERNAME": &c.Username,
		p + "_PASSWORD": &c.Password,
	}.MustLoad()
}

type Config struct {
	Postgres cc.Postgres `yaml:"postgres"`
	HTTP     cc.HTTP     `yaml:"http"`
	Env      string      `yaml:"env"`
	Vtiger   Vtiger      `yaml:"vtiger"`
	Vht      Vht         `yaml:"vht"`
}

func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8080,
		},
		Env:    "dev",
		Vtiger: DefaultVtiger(),
		Vht:    DefaultVht(),
	}
	cfg.Postgres.Database = "etop_crm"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	return cfg, err
}
