// +build wireinject

package sqlstore

import (
	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	NewAccountAuthStore,
	NewPartnerStore,
	NewExportAttemptStore,
	NewUserStore,

	BuildExportAttemptStore, BindExportAttemptStore,
	BuildPartnerStore, BindPartnerStore,
	BuildUserStore, BindUserStore,

	BindAccountStore, wire.Struct(new(AccountStore), "*"),
	BindAccountUserStore, wire.Struct(new(AccountUserStore), "*"),
	BindAddressStore, wire.Struct(new(AddressStore), "*"),
	BindCategoryStore, wire.Struct(new(CategoryStore), "*"),
	BindHistoryStore, wire.Struct(new(HistoryStore), "*"),
	BindOrderStore, wire.Struct(new(OrderStore), "*"),
	BindShopStore, wire.Struct(new(ShopStore), "*"),
	BindShopVariantStore, wire.Struct(new(ShopVariantStore), "*"),
	BindLogin, wire.Struct(new(Login), "*"),
)
