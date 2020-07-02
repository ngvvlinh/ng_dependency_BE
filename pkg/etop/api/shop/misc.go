package shop

import (
	"context"

	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/authorize/session"
)

type MiscService struct {
	session.Session
}

func (s *MiscService) Clone() api.MiscService { res := *s; return &res }

func (s *MiscService) VersionInfo(ctx context.Context, q *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	result := &pbcm.VersionInfoResponse{
		Service: "etop.Shop",
		Version: "0.1",
	}
	return result, nil
}
