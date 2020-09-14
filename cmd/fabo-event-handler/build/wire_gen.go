// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package build

import (
	"context"
	"o.o/backend/cmd/fabo-event-handler/config"
	"o.o/backend/com/eventhandler/handler/api"
	"o.o/backend/com/eventhandler/webhook/sender"
	"o.o/backend/com/eventhandler/webhook/storage"
	"o.o/backend/com/fabo/main/fbmessaging"
	"o.o/backend/com/fabo/main/fbpage"
	"o.o/backend/com/fabo/main/fbuser"
	"o.o/backend/com/main"
	"o.o/backend/com/main/identity"
	"o.o/backend/com/shopping/customering/query"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/sqlstore"
)

// Injectors from wire.go:

func Build(ctx context.Context, cfg config.Config) (Output, func(), error) {
	redisRedis := cfg.Redis
	store := redis.Connect(redisRedis)
	service := health.New(store)
	miscService := &api.MiscService{}
	databases := cfg.Databases
	mainDB, err := com.BuildDatabaseMain(databases)
	if err != nil {
		return Output{}, nil, err
	}
	webhookDB, err := com.BuildDatabaseWebhook(databases)
	if err != nil {
		return Output{}, nil, err
	}
	changesStore := storage.NewChangesStore(webhookDB)
	partnerStore := sqlstore.BuildPartnerStore(mainDB)
	partnerStoreInterface := sqlstore.BindPartnerStore(partnerStore)
	webhookSender := sender.New(mainDB, store, changesStore, partnerStoreInterface)
	webhookService := &api.WebhookService{
		Sender: webhookSender,
	}
	servers := api.NewServers(miscService, webhookService)
	mainServer, err := BuildMainServer(cfg, service, servers)
	if err != nil {
		return Output{}, nil, err
	}
	v := BuildServers(mainServer)
	handler, err := BuildIntHandler(ctx, cfg)
	if err != nil {
		return Output{}, nil, err
	}
	customerQuery := query.NewCustomerQuery(mainDB)
	queryBus := query.CustomerQueryMessageBus(customerQuery)
	fbUserQuery := fbuser.NewFbUserQuery(mainDB, queryBus)
	fbuseringQueryBus := fbuser.FbUserQueryMessageBus(fbUserQuery)
	fbMessagingQuery := fbmessaging.NewFbMessagingQuery(mainDB)
	fbmessagingQueryBus := fbmessaging.FbMessagingQueryMessageBus(fbMessagingQuery)
	fbPageQuery := fbpage.NewFbPageQuery(mainDB)
	fbpagingQueryBus := fbpage.FbPageQueryMessageBus(fbPageQuery)
	queryService := identity.NewQueryService(mainDB)
	identityQueryBus := identity.QueryServiceMessageBus(queryService)
	handlerHandler, err := BuildWebhookHandler(ctx, cfg, mainDB, fbuseringQueryBus, fbmessagingQueryBus, fbpagingQueryBus, identityQueryBus)
	if err != nil {
		return Output{}, nil, err
	}
	v2 := BuildWaiters(handler, handlerHandler)
	pgeventService, err := BuildPgEventService(ctx, cfg)
	if err != nil {
		return Output{}, nil, err
	}
	output := Output{
		Servers:   v,
		Waiters:   v2,
		PgService: pgeventService,
		WhSender:  webhookSender,
		Health:    service,
	}
	return output, func() {
	}, nil
}
