package integration

import (
	"context"

	"o.o/api/top/int/integration"
	pbcm "o.o/api/top/types/common"
)

type MiscService struct{}

func (s *MiscService) Clone() integration.MiscService {
	res := *s
	return &res
}

func (s *MiscService) VersionInfo(ctx context.Context, q *pbcm.Empty) (*pbcm.VersionInfoResponse, error) {
	result := &pbcm.VersionInfoResponse{
		Service: "etop.Integration",
		Version: "0.1",
	}
	return result, nil
}
