package xshop

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/redis"
	cmService "etop.vn/backend/pkg/common/service"

	cmP "etop.vn/backend/pb/common"
	etopP "etop.vn/backend/pb/etop"
	shopW "etop.vn/backend/wrapper/external/shop"
)

var (
	idempgroup *idemp.RedisGroup
	authStore  auth.Generator
)

const PrefixIdempShopAPI = "IdempShopAPI"

func init() {
	bus.AddHandlers("apix",
		VersionInfo,
		CurrentAccount,
	)
}

func Init(sd cmService.Shutdowner, rd redis.Store, s auth.Generator) {
	authStore = s

	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempShopAPI, 0)
	sd.Register(idempgroup.Shutdown)
}

func VersionInfo(ctx context.Context, q *shopW.VersionInfoEndpoint) error {
	q.Result = &cmP.VersionInfoResponse{
		Service: "shop",
		Version: "1.0.0",
	}
	return nil
}

func CurrentAccount(ctx context.Context, q *shopW.CurrentAccountEndpoint) error {
	if q.Context.Shop == nil {
		return cm.Errorf(cm.Internal, nil, "")
	}
	q.Result = etopP.PbPublicAccountInfo(q.Context.Shop)
	return nil
}
