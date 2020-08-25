package partner

import (
	"context"

	api "o.o/api/top/external/partner"
	externaltypes "o.o/api/top/external/types"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/etop/apix/convertpb"
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
		Service: "partner",
		Version: "1.0.0",
	}
	return result, nil
}

func (s *MiscService) GetLocationList(ctx context.Context, r *pbcm.Empty) (*externaltypes.LocationResponse, error) {
	resp, err := s.Shipping.GetLocationList(ctx)
	return resp, err
}

func (s *MiscService) CurrentAccount(ctx context.Context, q *pbcm.Empty) (*externaltypes.Partner, error) {
	if s.SS.Partner() == nil {
		return nil, cm.Errorf(cm.Internal, nil, "")
	}
	result := convertpb.PbPartner(s.SS.Partner())
	if wl.X(ctx).IsWhiteLabel() {
		result.Meta = map[string]string{
			"wl_name": wl.X(ctx).Name,
			"wl_key":  wl.X(ctx).Key,
			"wl_host": wl.X(ctx).Host,
		}
	}
	return result, nil
}
