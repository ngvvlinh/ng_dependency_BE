package main

import (
	"context"
	"flag"

	"o.o/api/top/types/etc/shipping_provider"
	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/com/main/shipping/modelx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/pkg/integration/shipping/ghtk"
	"o.o/common/l"
)

var (
	cfg config.Config
	ll  = l.New()
)

func main() {
	cc.InitFlags()
	flag.Parse()
	var err error
	if cfg, err = config.Load(false); err != nil {
		ll.Fatal("Error while loading config", l.Error(err))
	}
	cmenv.SetEnvironment(cfg.Env)

	if db, err := cmsql.Connect(cfg.Postgres); err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	} else {
		sqlstore.Init(db)
	}
	if cfg.GHTK.AccountDefault.Token != "" {
		if err := ghtk.Init(cfg.GHTK); err != nil {
			ll.Fatal("Unable to connect to GHTK", l.Error(err))
		}
		ll.S.Info("GHTK: connect success")
	} else {
		ll.Fatal("GHTK: No token")
	}

	ctx := context.Background()
	cmd := &modelx.GetUnCompleteFulfillmentsQuery{
		ShippingProviders: []shipping_provider.ShippingProvider{
			shipping_provider.GHTK,
		},
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		ll.Fatal("Can not get uncomplete ghtk fulfillments", l.Error(err))
	}
	if len(cmd.Result) == 0 {
		ll.Fatal("There aren't uncommplete ghtk fulfillments")
	}
	updateFfms, err := ghtk.SyncOrders(cmd.Result)
	if err != nil {
		ll.Fatal("Lỗi :: ", l.Error(err))
	}
	ll.Info("len :: ", l.Int("len", len(updateFfms)))
	if len(updateFfms) > 0 {
		cmdUpdate := &modelx.UpdateFulfillmentsWithoutTransactionCommand{
			Fulfillments: updateFfms,
		}
		if err := bus.Dispatch(ctx, cmdUpdate); err != nil {
			ll.Error("Không thể cập nhật ffm", l.Error(err))
		}
	}
}
