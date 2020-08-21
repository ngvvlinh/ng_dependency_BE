package database_all

import (
	com "o.o/backend/com/main"
	"o.o/backend/com/web"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
)

type Config struct {
	Postgres          cc.Postgres `yaml:"postgres"`
	PostgresWebServer cc.Postgres `yaml:"postgres_web_server"`
	PostgresLogs      cc.Postgres `yaml:"postgres_logs"`
	PostgresNotifier  cc.Postgres `yaml:"postgres_notifier"`
	PostgresAffiliate cc.Postgres `yaml:"postgres_affiliate"`
}

func DefaultConfig() Config {
	cfg := Config{
		Postgres:          cc.DefaultPostgres(),
		PostgresAffiliate: cc.DefaultPostgres(),
		PostgresLogs:      cc.DefaultPostgres(),
		PostgresNotifier:  cc.DefaultPostgres(),
		PostgresWebServer: cc.DefaultPostgres(),
	}
	cfg.Postgres.Database = "etop_dev"
	cfg.PostgresAffiliate.Database = "etop_dev"
	cfg.PostgresNotifier.Database = "etop_dev"
	cfg.PostgresWebServer.Database = "etop_dev"
	cfg.PostgresLogs.Database = "etop_dev"
	return cfg
}

func (cfg *Config) MustLoadEnv() {
	cc.PostgresMustLoadEnv(&cfg.Postgres)
	cc.PostgresMustLoadEnv(&cfg.PostgresLogs, "ET_POSTGRES_LOGS")
	cc.PostgresMustLoadEnv(&cfg.PostgresNotifier, "ET_POSTGRES_NOTIFIER")
	cc.PostgresMustLoadEnv(&cfg.PostgresAffiliate, "ET_POSTGRES_AFFILIATE")
	cc.PostgresMustLoadEnv(&cfg.PostgresWebServer, "ET_POSTGRES_WEB_SERVER")
}

type Databases struct {
	Main      com.MainDB
	Log       com.LogDB
	Notifier  com.NotifierDB
	Affiliate com.AffiliateDB
	WebServer web.WebServerDB
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
	dbs.Affiliate, err = cmsql.Connect(cfg.PostgresAffiliate)
	if err != nil {
		return dbs, err
	}
	dbs.WebServer, err = cmsql.Connect(cfg.PostgresWebServer)
	if err != nil {
		return dbs, err
	}
	return dbs, nil
}
