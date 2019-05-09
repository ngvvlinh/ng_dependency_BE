package config

import cc "etop.vn/backend/pkg/common/config"

type Config struct {
	Postgres         cc.Postgres `yaml:"postgres"`
	PostgresNotifier cc.Postgres `yaml:"postgres_notifier"`
	HTTP             cc.HTTP     `yaml:"http"`
	Kafka            cc.Kafka    `yaml:"kafka"`
	Env              string      `yaml:"env"`
	Secret           string      `yaml:"secret"`
}

func Default() Config {
	cfg := Config{
		Postgres:         cc.DefaultPostgres(),
		PostgresNotifier: cc.DefaultPostgres(),
		HTTP: cc.HTTP{
			Host: "",
			Port: 8082,
		},
		Kafka:  cc.DefaultKafka(),
		Env:    "dev",
		Secret: "secret",
	}
	cfg.Postgres.Database = "etop_dev"
	return cfg
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cc.PostgresMustLoadEnv(&cfg.PostgresNotifier, "ET_POSTGRES_NOTIFIER")
	cc.EnvMap{
		"ET_SECRET": &cfg.Secret,
	}.MustLoad()
	return cfg, err
}
