package config

import cc "etop.vn/backend/pkg/common/config"

type Config struct {
	Postgres        cc.Postgres `yaml:"postgres"`
	Redis           cc.Redis    `yaml:"redis"`
	HTTP            cc.HTTP     `yaml:"http"`
	Env             string      `yaml:"env"`
	VtigerService   string      `yaml:"vtiger_service"`
	VtigerUsername  string      `yaml:"vtiger_username"`
	VtigerAccesskey string      `yaml:"vtiger_accesskey"`
}

func Default() Config {
	cfg := Config{
		Postgres: cc.DefaultPostgres(),
		Redis:    cc.DefaultRedis(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8080,
		},
		Env:             "dev",
		VtigerService:   "http://vtiger/webservice.php",
		VtigerUsername:  "admin",
		VtigerAccesskey: "q5dZOnJYGlmPY2nc",
	}
	cfg.Postgres.Database = "etop_crm"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	return cfg, err
}
