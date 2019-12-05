package xshop

import (
	"context"

	pbcm "etop.vn/api/top/types/common"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/redis"
	cmservice "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/etop/api/convertpb"
)

var (
	idempgroup *idemp.RedisGroup
)

const PrefixIdempShopAPI = "IdempShopAPI"

func init() {
	bus.AddHandlers("apix",
		miscService.VersionInfo,
		miscService.CurrentAccount,
	)
}

func Init(sd cmservice.Shutdowner, rd redis.Store) {
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdempShopAPI, 0)
	sd.Register(idempgroup.Shutdown)
}

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "shop",
		Version: "1.0.0",
	}
	return nil
}

func (s *MiscService) CurrentAccount(ctx context.Context, q *CurrentAccountEndpoint) error {
	if q.Context.Shop == nil {
		return cm.Errorf(cm.Internal, nil, "")
	}
	q.Result = convertpb.PbPublicAccountInfo(q.Context.Shop)
	return nil
}
