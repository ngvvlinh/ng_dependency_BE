// +build wireinject

package sqlstore

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	New,
	NewAccountAuthStore,
	NewAccountStore,
	NewAccountUserStore,
	NewAddressStore,
	NewCategoryStore,
	NewExportAttemptStore,
	NewHistoryStore,
	NewMoneyTxStore,
	NewOrderStore,
	NewPartnerStore,
	NewShopStore,
	NewShopVariantStore,
	NewUserStore,
	wire.Struct(new(Login), "*"),

	// TODO(vu): remove
	wire.Bind(new(AccountStoreInterface), new(*AccountStore)),
	wire.Bind(new(AccountUserStoreInterface), new(*AccountUserStore)),
	wire.Bind(new(AddressStoreInterface), new(*AddressStore)),
	wire.Bind(new(CategoryStoreInterface), new(*CategoryStore)),
	wire.Bind(new(ExportAttemptStoreInterface), new(*ExportAttemptStore)),
	wire.Bind(new(HistoryStoreInterface), new(*HistoryStore)),
	wire.Bind(new(LoginInterface), new(*Login)),
	wire.Bind(new(MoneyTxStoreInterface), new(*MoneyTxStore)),
	wire.Bind(new(OrderStoreInterface), new(*OrderStore)),
	wire.Bind(new(PartnerStoreInterface), new(*PartnerStore)),
	wire.Bind(new(ShopStoreInterface), new(*ShopStore)),
	wire.Bind(new(ShopVariantStoreInterface), new(*ShopVariantStore)),
	wire.Bind(new(UserStoreInterface), new(*UserStore)),

	BuildExportAttempStore,
	BuildPartnerStore,
	BuildUserStore,
)
