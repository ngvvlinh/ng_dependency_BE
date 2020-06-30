package api

import (
	"context"

	"o.o/api/top/int/etop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
)

type MiscService struct{}

func (s *MiscService) Clone() etop.MiscService { res := *s; return &res }

func (s *MiscService) VersionInfo(ctx context.Context, q *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	result := &pbcm.VersionInfoResponse{
		Service: "etop",
		Version: "0.1",
		Meta: map[string]string{
			"name": wl.X(ctx).Name,
			"key":  wl.X(ctx).Key,
			"host": wl.X(ctx).Host,
		},
	}
	return result, nil
}
