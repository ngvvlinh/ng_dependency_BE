// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package build

import (
	"context"
	"o.o/backend/cmd/etop-event-handler/config"
	"o.o/backend/com/eventhandler/handler/api"
	"o.o/backend/com/eventhandler/webhook/sender"
	"o.o/backend/com/eventhandler/webhook/storage"
	"o.o/backend/com/main"
	"o.o/backend/com/main/catalog/query"
	query4 "o.o/backend/com/main/inventory/query"
	"o.o/backend/com/main/location"
	"o.o/backend/com/main/shipnow"
	query3 "o.o/backend/com/main/stocktaking/query"
	query2 "o.o/backend/com/shopping/customering/query"
	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/bus"
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
	webhookSender := sender.New(mainDB, store, changesStore)
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
	queryService := query.New(mainDB)
	queryBus := query.QueryServiceMessageBus(queryService)
	customerQuery := query2.NewCustomerQuery(mainDB)
	customeringQueryBus := query2.CustomerQueryMessageBus(customerQuery)
	stocktakeQuery := query3.NewQueryStocktake(mainDB)
	stocktakingQueryBus := query3.StocktakeQueryMessageBus(stocktakeQuery)
	busBus := bus.New()
	inventoryQueryService := query4.NewQueryInventory(stocktakingQueryBus, busBus, mainDB)
	inventoryQueryBus := query4.InventoryQueryServiceMessageBus(inventoryQueryService)
	addressQuery := query2.NewAddressQuery(mainDB)
	addressingQueryBus := query2.AddressQueryMessageBus(addressQuery)
	locationQuery := location.New(mainDB)
	locationQueryBus := location.QueryMessageBus(locationQuery)
	shipnowQueryService := shipnow.NewQueryService(mainDB)
	shipnowQueryBus := shipnow.QueryServiceMessageBus(shipnowQueryService)
	handlerHandler, err := BuildWebhookHandler(ctx, cfg, mainDB, webhookSender, queryBus, customeringQueryBus, inventoryQueryBus, addressingQueryBus, locationQueryBus, shipnowQueryBus)
	if err != nil {
		return Output{}, nil, err
	}
	v2 := BuildWaiters(handler, handlerHandler)
	notifierDB := _wireNotifierDBValue
	sqlstoreStore := sqlstore.New(mainDB, notifierDB, locationQueryBus, busBus)
	pgeventService, err := BuildPgEventService(ctx, cfg)
	if err != nil {
		return Output{}, nil, err
	}
	output := Output{
		Servers:   v,
		Waiters:   v2,
		Store:     sqlstoreStore,
		PgService: pgeventService,
		WhSender:  webhookSender,
		Health:    service,
	}
	return output, func() {
	}, nil
}

var (
	_wireNotifierDBValue = com.NotifierDB(nil)
)
