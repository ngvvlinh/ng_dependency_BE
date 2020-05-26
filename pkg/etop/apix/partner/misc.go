package partner

import (
	"context"

	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/apix/shipping"
)

type MiscService struct {
	Shipping *shipping.Shipping
}

func (s *MiscService) Clone() *MiscService { res := *s; return &res }

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "partner",
		Version: "1.0.0",
	}
	return nil
}

func (s *MiscService) GetLocationList(ctx context.Context, r *GetLocationListEndpoint) error {
	resp, err := s.Shipping.GetLocationList(ctx)
	r.Result = resp
	return err
}

func (s *MiscService) CurrentAccount(ctx context.Context, q *CurrentAccountEndpoint) error {
	if q.Context.Partner == nil {
		return cm.Errorf(cm.Internal, nil, "")
	}
	q.Result = convertpb.PbPartner(q.Context.Partner)
	if wl.X(ctx).IsWhiteLabel() {
		q.Result.Meta = map[string]string{
			"wl_name": wl.X(ctx).Name,
			"wl_key":  wl.X(ctx).Key,
			"wl_host": wl.X(ctx).Host,
		}
	}
	return nil
}
