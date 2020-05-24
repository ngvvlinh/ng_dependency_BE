package api

import (
	"context"

	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
)

type MiscService struct{}

func (s *MiscService) Clone() *MiscService { res := *s; return &res }

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop",
		Version: "0.1",
		Meta: map[string]string{
			"name": wl.X(ctx).Name,
			"key":  wl.X(ctx).Key,
			"host": wl.X(ctx).Host,
		},
	}
	return nil
}
