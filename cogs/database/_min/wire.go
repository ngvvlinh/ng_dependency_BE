package database_min

import (
	"github.com/google/wire"

	com "o.o/backend/com/main"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
)

var WireSet = wire.NewSet(
	wire.FieldsOf(new(Databases), "main", "log", "notifier"),
	BuildDatabases,
)

type Config struct {
	Postgres         cc.Postgres `yaml:"postgres"`
	PostgresLogs     cc.Postgres `yaml:"postgres_logs"`
	PostgresNotifier cc.Postgres `yaml:"postgres_notifier"`
}

func DefaultConfig() Config {
	cfg := Config{
		Postgres:         cc.DefaultPostgres(),
		PostgresNotifier: cc.DefaultPostgres(),
		PostgresLogs:     cc.DefaultPostgres(),
	}
	cfg.Postgres.Database = "etop_dev"
	return cfg
}

func (cfg *Config) MustLoadEnv() {
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cc.PostgresMustLoadEnv(&cfg.PostgresLogs, "ET_POSTGRES_LOGS")
	cc.PostgresMustLoadEnv(&cfg.PostgresNotifier, "ET_POSTGRES_NOTIFIER")
}

type Databases struct {
	Main     com.MainDB
	Log      com.LogDB
	Notifier com.NotifierDB
	Webhook  com.WebhookDB
}

func BuildDatabases(cfg Config) (dbs Databases, err error) {
	dbs.Main, err = cmsql.Connect(cfg.Postgres)
	if err != nil {
		return dbs, err
	}
	dbs.Log, err = cmsql.Connect(cfg.PostgresLogs)
	if err != nil {
		return dbs, err
	}
	dbs.Notifier, err = cmsql.Connect(cfg.PostgresNotifier)
	if err != nil {
		return dbs, err
	}
	return dbs, nil
}
