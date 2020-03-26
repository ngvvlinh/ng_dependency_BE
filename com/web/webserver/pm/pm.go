package pm

import (
	"context"

	"o.o/api/main/catalog"
	"o.o/api/webserver"
	"o.o/backend/pkg/common/bus"
	"o.o/capi"
	"o.o/capi/dot"
)

type ProcessManager struct {
	eventBus           capi.EventBus
	webserverCommanBus webserver.CommandBus
	webserverQueryBus  webserver.QueryBus
}

func New(
	eventBusArgs capi.EventBus,
	webserverCommanBusArg webserver.CommandBus,
	webserverQueryBusArg webserver.QueryBus,
) *ProcessManager {
	return &ProcessManager{
		eventBus:           eventBusArgs,
		webserverCommanBus: webserverCommanBusArg,
		webserverQueryBus:  webserverQueryBusArg,
	}
}

func (m *ProcessManager) RegisterEventHandlers(eventBus bus.EventRegistry) {
	eventBus.AddEventListener(m.ShopProductDeletedEvent)
}

func (m *ProcessManager) ShopProductDeletedEvent(ctx context.Context, event *catalog.ShopProductDeletedEvent) error {
	query := &webserver.ListWsWebsitesQuery{
		ShopID: event.ShopID,
	}
	err := m.webserverQueryBus.Dispatch(ctx, query)
	if err != nil {
		return err
	}
	var mapProductIDDelted = make(map[dot.ID]bool)
	for _, v := range event.ProductIDs {
		mapProductIDDelted[v] = true
	}
	for _, wsWebsite := range query.Result.WsWebsites {
		var outStandingProducts = wsWebsite.OutstandingProduct
		var newProducts = wsWebsite.NewProduct
		// remove deleted product in out standing product
		var outStandingProductsTemp []dot.ID
		for _, productID := range outStandingProducts.ProductIDs {
			if !mapProductIDDelted[productID] {
				outStandingProductsTemp = append(outStandingProductsTemp, productID)
			}
		}
		outStandingProducts.ProductIDs = outStandingProductsTemp

		var newProductsTemp []dot.ID
		for _, productID := range newProducts.ProductIDs {
			if !mapProductIDDelted[productID] {
				newProductsTemp = append(newProductsTemp, productID)
			}
		}
		newProducts.ProductIDs = newProductsTemp
		cmd := &webserver.UpdateWsWebsiteCommand{
			ShopID:             event.ShopID,
			ID:                 wsWebsite.ID,
			OutstandingProduct: outStandingProducts,
			NewProduct:         newProducts,
		}
		err = m.webserverCommanBus.Dispatch(ctx, cmd)
		if err != nil {
			return err
		}
	}
	return nil
}
