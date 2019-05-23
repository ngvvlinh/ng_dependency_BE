package xshop

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/auth"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/redis"
	cmservice "etop.vn/backend/pkg/common/service"

	pbcm "etop.vn/backend/pb/common"
	pbetop "etop.vn/backend/pb/etop"
	wrapshop "etop.vn/backend/wrapper/external/shop"
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

func Init(sd cmservice.Shutdowner, rd redis.Store, s auth.Generator) {
	authStore = s

	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempShopAPI, 0)
	sd.Register(idempgroup.Shutdown)
}

func VersionInfo(ctx context.Context, q *wrapshop.VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "shop",
		Version: "1.0.0",
	}
	return nil
}

func CurrentAccount(ctx context.Context, q *wrapshop.CurrentAccountEndpoint) error {
	if q.Context.Shop == nil {
		return cm.Errorf(cm.Internal, nil, "")
	}
	q.Result = pbetop.PbPublicAccountInfo(q.Context.Shop)
	return nil
}
