package vnp

import (
	"context"

	"o.o/api/top/external/mc/vnp"
	pbcm "o.o/api/top/types/common"
	xshop "o.o/backend/pkg/etop/apix/shop"
	"o.o/backend/pkg/etop/authorize/session"
)

var _ vnp.ShipnowService = &VNPostService{}

type VNPostService struct {
	session.Session
	xshop.ShipnowService
}

func (s *VNPostService) Clone() vnp.ShipnowService {
	res := *s
	res.Session.Link(&res.ShipnowService)
	return &res
}

func (s *VNPostService) Ping(ctx context.Context, empty *pbcm.Empty) (*pbcm.Empty, error) {
	return &pbcm.Empty{}, nil
}
