package sadmin

import (
	"context"

	"o.o/api/top/int/sadmin"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/authorize/session"
)

type MiscService struct {
	session.Session
}

func (s *MiscService) Clone() sadmin.MiscService {
	res := *s
	return &res
}

func (s *MiscService) VersionInfo(ctx context.Context, q *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	res := &pbcm.VersionInfoResponse{
		Service: "etop.SuperAdmin",
		Version: "0.1",
	}
	return res, nil
}
