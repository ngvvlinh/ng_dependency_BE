package shop

import (
	"context"

	pbcm "o.o/api/top/types/common"
)

type MiscService struct{}

func (s *MiscService) Clone() *MiscService { res := *s; return &res }

func (s *MiscService) VersionInfo(ctx context.Context, q *VersionInfoEndpoint) error {
	q.Result = &pbcm.VersionInfoResponse{
		Service: "etop.Shop",
		Version: "0.1",
	}
	return nil
}
