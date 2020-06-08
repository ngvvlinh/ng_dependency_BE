package handler

import (
	"time"

	"o.o/api/main/catalog"
	"o.o/api/main/inventory"
	"o.o/api/main/location"
	"o.o/api/shopping/addressing"
	"o.o/api/shopping/customering"
	"o.o/api/top/external/types"
	"o.o/backend/com/eventhandler/pgevent"
	"o.o/backend/com/eventhandler/webhook/sender"
	com "o.o/backend/com/main"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/mq"
	"o.o/backend/pkg/common/sql/cmsql"
	historysqlstore "o.o/backend/pkg/etop-history/sqlstore"
	"o.o/common/l"
)

var ll = l.New()

const ConsumerGroup = "handler/webhook"

type Handler struct {
	db           *cmsql.Database
	historyStore historysqlstore.HistoryStoreFactory
	sender       *sender.WebhookSender

	catalogQuery   catalog.QueryBus
	customerQuery  customering.QueryBus
	inventoryQuery inventory.QueryBus
	addressQuery   addressing.QueryBus
	locationQuery  location.QueryBus
}

func New(
	db com.MainDB,
	sender *sender.WebhookSender,
	catalogQuery catalog.QueryBus,
	customerQuery customering.QueryBus,
	inventoryQuery inventory.QueryBus,
	addressQuery addressing.QueryBus,
	locationQuery location.QueryBus,
) *Handler {
	h := &Handler{
		db:             db,
		historyStore:   historysqlstore.NewHistoryStore(db),
		sender:         sender,
		catalogQuery:   catalogQuery,
		customerQuery:  customerQuery,
		inventoryQuery: inventoryQuery,
		addressQuery:   addressQuery,
		locationQuery:  locationQuery,
	}
	return h
}

func (h *Handler) TopicsAndHandlers() map[string]mq.EventHandler {
	return pgevent.WrapMapHandlers(map[string]pgevent.HandlerFunc{
		"fulfillment":                  h.HandleFulfillmentEvent,
		"order":                        h.HandleOrderEvent,
		"notification":                 nil,
		"money_transaction_shipping":   nil,
		"shop_product":                 h.HandleShopProductEvent,
		"shop_variant":                 h.HandleShopVariantEvent,
		"shop_customer":                h.HandleShopCustomerEvent,
		"shop_customer_group":          h.HandleShopCustomerGroupEvent,
		"shop_customer_group_customer": h.HandleShopCustomerGroupCustomerEvent,
		"inventory_variant":            h.HandleInventoryVariantEvent,
		"shop_trader_address":          h.HandleShopTraderAddressEvent,
		"shop_collection":              h.HandleShopProductCollectionEvent,
		"shop_product_collection":      h.HandleShopProductionCollectionRelationshipEvent,
	})
}

func pbChange(event *pgevent.PgEvent) *types.Change {
	return &types.Change{
		Time:       cmapi.PbTime(time.Unix(event.Timestamp, 0)),
		ChangeType: pbChangeType(event.Op),
		Entity:     event.Table,
	}
}

func pbChangeType(op pgevent.TGOP) string {
	switch op {
	case pgevent.OpInsert:
		return "create"
	case pgevent.OpUpdate:
		return "update"
	case pgevent.OpDelete:
		return "delete"
	default:
		return ""
	}
}
