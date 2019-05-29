package main

import (
	"context"
	"flag"

	"etop.vn/backend/pkg/services/shipping/modelx"

	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/sqlstore"
	"etop.vn/backend/pkg/integration/ghtk"

	"etop.vn/backend/cmd/etop-server/config"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/l"
	_ "etop.vn/backend/pkg/etop/sqlstore"
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
	cm.SetEnvironment(cfg.Env)

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
		ShippingProviders: []model.ShippingProvider{
			model.TypeGHTK,
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