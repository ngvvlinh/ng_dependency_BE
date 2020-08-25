package xshop

import (
	"context"

	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/apix/shipping"
	"o.o/backend/pkg/etop/authorize/session"
)

type MiscService struct {
	session.Session

	Shipping *shipping.Shipping
}

func (s *MiscService) Clone() api.MiscService { res := *s; return &res }

func (s *MiscService) VersionInfo(ctx context.Context, q *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	result := &pbcm.VersionInfoResponse{
		Service: "shop",
		Version: "1.0.0",
	}
	return result, nil
}

func (s *MiscService) CurrentAccount(ctx context.Context, q *pbcm.Empty) (*etop.PublicAccountInfo, error) {
	if s.SS.Shop() == nil {
		return nil, cm.Errorf(cm.Internal, nil, "")
	}
	result := convertpb.PbPublicAccountInfo(s.SS.Shop())
	return result, nil
}

func (s *MiscService) GetLocationList(ctx context.Context, r *pbcm.Empty) (*externaltypes.LocationResponse, error) {
	resp, err := s.Shipping.GetLocationList(ctx)
	return resp, err
}
