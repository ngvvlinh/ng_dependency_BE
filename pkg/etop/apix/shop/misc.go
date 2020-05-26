package xshop

import (
	"context"

	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/apix/shipping"
)

type MiscService struct {
	Shipping *shipping.Shipping
}

func (s *MiscService) Clone() *MiscService { res := *s; return &res }

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

func (s *MiscService) GetLocationList(ctx context.Context, r *GetLocationListEndpoint) error {
	resp, err := s.Shipping.GetLocationList(ctx)
	r.Result = resp
	return err
}
