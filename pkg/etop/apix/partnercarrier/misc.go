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
	session.Sessioner
	ss *session.Session
}

func NewMiscService(ss *session.Session) *MiscService {
	return &MiscService{
		ss: ss,
	}
}

func (s *MiscService) Clone() partnercarrier.MiscService {
	res := *s
	res.Sessioner, res.ss = s.ss.Split()
	return &res
}

func (s *MiscService) GetLocationList(ctx context.Context, q *pbcm.Empty) (*externaltypes.LocationResponse, error) {
	return shipping.GetLocationList(ctx)
}

func (s *MiscService) CurrentAccount(ctx context.Context, q *pbcm.Empty) (*externaltypes.Partner, error) {
	partner := s.ss.Partner()
	if partner == nil {
		return nil, cm.Errorf(cm.Internal, nil, "")
	}
	return convertpb.PbPartner(partner), nil
}
