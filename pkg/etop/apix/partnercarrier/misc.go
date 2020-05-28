package partnercarrier

import (
	"context"

	"o.o/api/top/external/partnercarrier"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/apix/shipping"
	"o.o/backend/pkg/etop/authorize/session"
)

type MiscService struct {
	session.Session
	Shipping *shipping.Shipping
}

func (s *MiscService) Clone() partnercarrier.MiscService {
	res := *s
	return &res
}

func (s *MiscService) GetLocationList(ctx context.Context, q *pbcm.Empty) (*externaltypes.LocationResponse, error) {
	return s.Shipping.GetLocationList(ctx)
}

func (s *MiscService) CurrentAccount(ctx context.Context, q *pbcm.Empty) (*externaltypes.Partner, error) {
	partner := s.SS.Partner()
	if partner == nil {
		return nil, cm.Errorf(cm.Internal, nil, "")
	}
	return convertpb.PbPartner(partner), nil
}
