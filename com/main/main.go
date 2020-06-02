package com

import (
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
)

type MainDB *cmsql.Database
type LogDB *cmsql.Database
type NotifierDB *cmsql.Database
type AffiliateDB *cmsql.Database
type WebhookDB *cmsql.Database

func BuildDatabaseMain(c cc.Databases) (MainDB, error) {
	cfg, err := c.Get("postgres")
	if err != nil {
		return nil, err
	}
	return cmsql.Connect(cfg)
}

func BuildDatabaseLogs(c cc.Databases) (LogDB, error) {
	cfg, err := c.Get("postgres_logs")
	if err != nil {
		return nil, err
	}
	return cmsql.Connect(cfg)
}

func BuildDatabaseWebhook(c cc.Databases) (WebhookDB, error) {
	cfg, err := c.Get("postgres_webhook")
	if err != nil {
		return nil, err
	}
	return cmsql.Connect(cfg)
}

func BuildDatabaseNotifier(c cc.Databases) (NotifierDB, error) {
	cfg, err := c.Get("postgres_notifier")
	if err != nil {
		return nil, err
	}
	return cmsql.Connect(cfg)
}

func BuildDatabaseAffiliateDB(c cc.Databases) (AffiliateDB, error) {
	cfg, err := c.Get("postgres_affiliate")
	if err != nil {
		return nil, err
	}
	return cmsql.Connect(cfg)
}
